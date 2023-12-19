package restserver

import (
	"encoding/json"
	"net/http"

	"gitlab.utc.fr/vaursdam/formule-1-ia04/simulator"
	"gitlab.utc.fr/vaursdam/formule-1-ia04/types"
)

var driverTotalPoints []*types.DriverTotalPoints
var teamTotalPoints []*types.TeamTotalPoints
var personalityTotalPoints []*types.PersonalityTotalPoints

// Lancement d'une simulation
func (rsa *RestServer) startSimulation(w http.ResponseWriter, r *http.Request) {

	// vérification de la méthode de la requête
	if r.Method != "GET" {
		return
	}

	championship := types.NewChampionship("2023", "Championship 1", rsa.pointTabCircuit, rsa.pointTabTeam)
	s := simulator.NewSimulator([]types.Championship{*championship})

	// Lancement de la simulation
	driverTotalPoints, teamTotalPoints, personalityTotalPoints = s.LaunchSimulation()
	lastChampionshipStatistics := types.NewLastChampionshipStatistics(driverTotalPoints, teamTotalPoints, personalityTotalPoints, nil)
	simulateChampionship := types.NewSimulateChampionship(championship.Name, *lastChampionshipStatistics)
	w.WriteHeader(http.StatusOK)
	serial, _ := json.Marshal(simulateChampionship)
	w.Write(serial)
}
