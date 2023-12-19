package restserver

import (
	"bytes"
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

// Décodage de la requête /personalities/update
func (*RestServer) decodeUpdatePersonalityRequest(r *http.Request) (req []types.UpdatePersonalityInfo, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	return
}

func (rsa *RestServer) startSimulation(w http.ResponseWriter, r *http.Request) {

	// vérification de la méthode de la requête
	if r.Method != "POST" {
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
	if r.Method != "GET" {
		return
	}

	serial, _ := json.Marshal(driversRank)
	w.Write(serial)
}

// Obtenir les personnalités d'une simulation
func (rsa *RestServer) getAndUpdatePersonalities(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" { // Obtenir les personnalités
		driversInfosPersonnalities := make([]types.PersonalityInfo, 0)

		for _, team := range rsa.pointTabTeam {
			team := *team
			for _, driver := range team.Drivers {
				driverInfo := types.PersonalityInfo{
					IdDriver:    driver.Id,
					Lastname:    driver.Lastname,
					Personality: driver.Personality.TraitsValue,
				}
				driversInfosPersonnalities = append(driversInfosPersonnalities, driverInfo)
			}

		}

		serial, _ := json.Marshal(driversInfosPersonnalities)
		w.WriteHeader(http.StatusOK)
		w.Write(serial)
		return
	} else if r.Method == "PUT" { // Mettre à jour les personnalités
		// décodage de la requête
		req, err := rsa.decodeUpdatePersonalityRequest(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, err.Error())
			return
		}

		// Réponse à renvoyer
		var resp []types.UpdatePersonalityInfo

		// Parcours des équipes et des pilotes
		for _, team := range rsa.pointTabTeam {
			for i := 0; i < 2; i++ {
				for _, updateReq := range req {
					if updateReq.IdDriver == team.Drivers[i].Id {
						// Test des valeurs des différentes personnalités
						if updateReq.Personality["Aggressivity"] > 5 || updateReq.Personality["Aggressivity"] < 1 ||
							updateReq.Personality["Confidence"] > 5 || updateReq.Personality["Confidence"] < 1 ||
							updateReq.Personality["Docility"] > 5 || updateReq.Personality["Docility"] < 1 ||
							updateReq.Personality["Concentration"] > 5 || updateReq.Personality["Concentration"] < 1 {
							msg := "Une des valeurs des personnalités entrée est supérieur à 5 ou inférieur à 1"
							w.WriteHeader(http.StatusBadRequest)
							serial, _ := json.Marshal(msg)
							w.Write(serial)
							return
						} else {
							// Mise à jour des valeurs de personnalité du pilote
							team.Drivers[i].Personality.TraitsValue = map[string]int{
								"Aggressivity":  updateReq.Personality["Aggressivity"],
								"Confidence":    updateReq.Personality["Confidence"],
								"Docility":      updateReq.Personality["Docility"],
								"Concentration": updateReq.Personality["Concentration"],
							}
						}

						// Remplissage de la réponse
						resp = append(resp, types.UpdatePersonalityInfo{
							IdDriver:    team.Drivers[i].Id,
							Personality: team.Drivers[i].Personality.TraitsValue,
						})
					}
				}
			}
		}

		serial, _ := json.Marshal(resp)

		w.WriteHeader(http.StatusOK)
		w.Write(serial)
		return
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "method %q not allowed", r.Method)
		return
	}
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
	mux.HandleFunc("/personalities", rsa.getAndUpdatePersonalities)

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
