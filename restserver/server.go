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
}

var driversRank []*types.DriverRank

func NewRestServer(addr string, pointTabCircuit []*types.Circuit, pointTabTeam []*types.Team) *RestServer {
	return &RestServer{addr: addr, pointTabCircuit: pointTabCircuit, pointTabTeam: pointTabTeam}
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

	driversInfosPersonnalities := make([]types.PersonalityInfo, 0)

	for _, team := range rsa.pointTabTeam {
		team := *team
		for _, driver := range team.Drivers {
			fmt.Println(driver)
			driverInfo := types.PersonalityInfo{
				IdDriver:    driver.Id,
				Lastname:    driver.Lastname,
				Personality: driver.Personality.TraitsValue,
			}
			driversInfosPersonnalities = append(driversInfosPersonnalities, driverInfo)
		}

	}
	fmt.Println(driversInfosPersonnalities[0])

	serial, _ := json.Marshal(driversInfosPersonnalities)
	w.WriteHeader(http.StatusOK)
	w.Write(serial)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (rsa *RestServer) Start() {
	// création du multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/api/startSimulation", rsa.startSimulation)
	mux.HandleFunc("/api/driversChampionshipRank", rsa.getChampionshipRank)
	mux.HandleFunc("/personalities", rsa.getPersonnalities)

	corsHandler := corsMiddleware(mux)

	// création du serveur http
	s := &http.Server{
		Addr:           rsa.addr,
		Handler:        corsHandler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20}

	// lancement du serveur
	log.Println("Listening on", rsa.addr)
	go log.Fatal(s.ListenAndServe())

}
