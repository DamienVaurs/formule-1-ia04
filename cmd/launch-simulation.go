package main

import (
	"log"
	"net/http"
	"time"

	"gitlab.utc.fr/vaursdam/formule-1-ia04/simulator"
	"gitlab.utc.fr/vaursdam/formule-1-ia04/types"
	"gitlab.utc.fr/vaursdam/formule-1-ia04/utils"
)

func main() {
	c, err := utils.ReadCircuit()
	if err != nil {
		panic(err)
	}

	t, err := utils.ReadTeams()
	if err != nil {
		panic(err)
	}

	//On crée des pointeurs vers les équipes et les circuits
	pointTabCircuit := make([]*types.Circuit, len(c))
	for i, circuit := range c {
		tempCircuit := circuit //sans tampon, tous les éléments du tableau contiendront la même adresse
		pointTabCircuit[i] = &tempCircuit
	}

	pointTabTeam := make([]*types.Team, len(t))
	for i, team := range t {
		tempTeam := team //sans tampon, tous les éléments du tableau contiendront la même adresse
		pointTabTeam[i] = &tempTeam
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Bienvenue sur le serveur"))
	})

	mux.HandleFunc("/api/startSimulation", func(w http.ResponseWriter, r *http.Request) {
		//On a les équipes et les circuits, on lance la simulation
		championship := types.NewChampionship("2023", "Championship 1", pointTabCircuit, pointTabTeam)
		s := simulator.NewSimulator([]types.Championship{*championship})
		s.LaunchSimulation()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Simulation démarrée"))
	})

	mux.HandleFunc("/api/start100Simulation", func(w http.ResponseWriter, r *http.Request) {
		//On a les équipes et les circuits, on lance la simulation
		championship := make([]types.Championship, 0)
		for i := 0; i < 100; i++ {
			championship = append(championship, *types.NewChampionship("2023", "Championship 1", pointTabCircuit, pointTabTeam))
		}
		s := simulator.NewSimulator(championship)
		s.LaunchSimulation()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Simulation démarrée"))
	})

	server := &http.Server{
		Addr:           ":8080",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// lancement du serveur
	log.Println("Listening on", server.Addr)
	log.Fatal(server.ListenAndServe())
}
