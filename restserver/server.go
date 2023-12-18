package restserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"gitlab.utc.fr/vaursdam/formule-1-ia04/simulator"
	"gitlab.utc.fr/vaursdam/formule-1-ia04/types"
)

type RestServer struct {
	sync.Mutex
	addr            string
	pointTabCircuit []*types.Circuit
	pointTabTeam    []*types.Team
	drivers         []types.Driver
}

var driversRank []*types.DriverRank

func NewRestServer(addr string, pointTabCircuit []*types.Circuit, pointTabTeam []*types.Team, drivers []types.Driver) *RestServer {
	return &RestServer{addr: addr, pointTabCircuit: pointTabCircuit, pointTabTeam: pointTabTeam, drivers: drivers}
}

// Test de la méthode
func (rsa *RestServer) checkMethod(method string, w http.ResponseWriter, r *http.Request) bool {
	if r.Method != method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "method %q not allowed", r.Method)
		return false
	}
	return true
}

func (rsa *RestServer) startSimulation(w http.ResponseWriter, r *http.Request) {

	// vérification de la méthode de la requête
	if !rsa.checkMethod("POST", w, r) {
		return
	}

	championship := types.NewChampionship("2023", "Championship 1", rsa.pointTabCircuit, rsa.pointTabTeam)
	s := simulator.NewSimulator([]types.Championship{*championship})
	driversRank = s.LaunchSimulation()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Simulation terminée"))
}

// Obtenir le classements des pilotes à la fin d'un championnat
func (rsa *RestServer) getChampionshipRank(w http.ResponseWriter, r *http.Request) {
	// vérification de la méthode de la requête
	if !rsa.checkMethod("GET", w, r) {
		return
	}

	serial, _ := json.Marshal(driversRank)
	w.Write(serial)
}

// Obtenir les personnalités d'une simulation
func (rsa *RestServer) getPersonnalities(w http.ResponseWriter, r *http.Request) {
	// vérification de la méthode de la requête
	if !rsa.checkMethod("GET", w, r) {
		return
	}

	driversInfosPersonnalities := make([]types.PersonnalityInfo, 0)

	for _, team := range rsa.pointTabTeam {
		team := *team
		for _, driver := range team.Drivers {
			driverInfo := types.PersonnalityInfo{
				IdDriver:     driver.Id,
				Lastname:     driver.Lastname,
				Personnality: driver.Personnality.TraitsValue,
			}
			driversInfosPersonnalities = append(driversInfosPersonnalities, driverInfo)
		}

	}

	serial, _ := json.Marshal(driversInfosPersonnalities)
	w.WriteHeader(http.StatusOK)
	w.Write(serial)
}

func (rsa *RestServer) Start() {
	// création du multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/api/startSimulation", rsa.startSimulation)
	mux.HandleFunc("/api/driversChampionshipRank", rsa.getChampionshipRank)
	mux.HandleFunc("/personnalities", rsa.getPersonnalities)

	// création du serveur http
	s := &http.Server{
		Addr:           rsa.addr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20}

	// lancement du serveur
	log.Println("Listening on", rsa.addr)
	go log.Fatal(s.ListenAndServe())

}
