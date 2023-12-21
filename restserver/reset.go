package restserver

import (
	"encoding/json"
	"net/http"

	"gitlab.utc.fr/vaursdam/formule-1-ia04/types"
)

func (rsa *RestServer) reset(w http.ResponseWriter, r *http.Request) {
	// vérification de la méthode de la requête
	if r.Method != "GET" {
		return
	}
	// reset des variables globales
	nextChampionship = "2023/2024"
	nbSimulation = 0
	statistics = &types.SimulateChampionship{}
	//On remet les personnalités à 0
	for indTeam := range rsa.pointTabTeam {
		for indDriver := 0; indDriver < 2; indDriver++ {
			rsa.pointTabTeam[indTeam].Drivers[indDriver].Personality = rsa.initPersonalities[rsa.pointTabTeam[indTeam].Drivers[indDriver].Id]
		}
	}
	serial, err := json.Marshal(statistics) //statistics is defined in simulateChampionship.go
	if err != nil {
		panic("Error /reset : can't marshal statistics" + err.Error())
	}
	w.Write(serial)
}
