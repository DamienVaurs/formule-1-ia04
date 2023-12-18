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

// Lancement de 50 simulations
func (rsa *RestServer) start50Simulation(w http.ResponseWriter, r *http.Request) {

	// vérification de la méthode de la requête
	if !rsa.checkMethod("POST", w, r) {
		return
	}

	//On a les équipes et les circuits, on lance la simulation
	championship := make([]types.Championship, 0)
	for i := 0; i < 50; i++ {
		championship = append(championship, *types.NewChampionship("2023", "Championship 1", rsa.pointTabCircuit, rsa.pointTabTeam))
	}
	s := simulator.NewSimulator(championship)
	s.LaunchSimulation()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Simulation de 50 championnats terminée"))
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

// Obtenir les pilotes avant une simulation
func (rsa *RestServer) getDrivers(w http.ResponseWriter, r *http.Request) {
	// vérification de la méthode de la requête
	if !rsa.checkMethod("GET", w, r) {
		return
	}

	serial, _ := json.Marshal(rsa.drivers)
	w.Write(serial)
}

func (rsa *RestServer) Start() {
	// création du multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/api/startSimulation", rsa.startSimulation)
	mux.HandleFunc("/api/start50Simulation", rsa.start50Simulation)
	mux.HandleFunc("/api/driversChampionshipRank", rsa.getChampionshipRank)
	mux.HandleFunc("/api/drivers", rsa.getDrivers)

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