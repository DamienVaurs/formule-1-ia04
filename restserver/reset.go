package restserver

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.utc.fr/vaursdam/formule-1-ia04/types"
)

func (rsa *RestServer) reset(w http.ResponseWriter, r *http.Request) {
	// vérification de la méthode de la requête
	if r.Method != "GET" {
		return
	}
	fmt.Println("GET /reset")
	// reset des variables globales
	nextChampionship = "2023/2024"
	nbSimulation = 0
	statistics = &types.SimulateChampionship{}
	for _, team := range rsa.pointTabTeam {
		for _, driver := range team.Drivers {
			statistics.TotalStatistics.DriversTotalPoints = append(statistics.TotalStatistics.DriversTotalPoints, &types.DriverTotalPoints{Driver: driver.Lastname, TotalPoints: 0})
			statistics.LastChampionshipStatistics.DriversTotalPoints = append(statistics.LastChampionshipStatistics.DriversTotalPoints, &types.DriverTotalPoints{Driver: driver.Lastname, TotalPoints: 0})
		}
		statistics.TotalStatistics.TeamsTotalPoints = append(statistics.TotalStatistics.TeamsTotalPoints, &types.TeamTotalPoints{Team: team.Name, TotalPoints: 0})
		statistics.LastChampionshipStatistics.TeamsTotalPoints = append(statistics.LastChampionshipStatistics.TeamsTotalPoints, &types.TeamTotalPoints{Team: team.Name, TotalPoints: 0})
	}
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
