package restserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gitlab.utc.fr/vaursdam/formule-1-ia04/simulator"
	"gitlab.utc.fr/vaursdam/formule-1-ia04/types"
)

var driverTotalPoints []*types.DriverTotalPoints
var teamTotalPoints []*types.TeamTotalPoints
var personalityAveragePoints []*types.PersonalityAveragePoints
var nextChampionship = "2023/2024"

func getNextChampionshipName(currChampionship string) (string, error) {
	years := strings.Split(currChampionship, "/")
	newFirstYear, err := time.Parse("2006", years[0]) //on souhaite récupérer la première année
	if err != nil {
		return "", err
	}

	newFirstYear = newFirstYear.AddDate(1, 0, 0)
	newLastYear := newFirstYear.AddDate(1, 0, 0)
	return fmt.Sprintf("%d/%d", newFirstYear.Year(), newLastYear.Year()), nil

}

// Lancement d'une simulation d'un championnat
func (rsa *RestServer) startSimulation(w http.ResponseWriter, r *http.Request) {

	// vérification de la méthode de la requête
	if r.Method != "GET" {
		return
	}

	championship := types.NewChampionship(nextChampionship, nextChampionship, rsa.pointTabCircuit, rsa.pointTabTeam)
	ch, err := getNextChampionshipName(nextChampionship)
	if err != nil {
		panic("Error /simulateChampionship : can't create new Dates" + err.Error())
	}
	nextChampionship = ch

	s := simulator.NewSimulator([]types.Championship{*championship})

	// Lancement de la simulation
	driverTotalPoints, teamTotalPoints, personalityAveragePoints = s.LaunchSimulation()
	lastChampionshipStatistics := types.NewLastChampionshipStatistics(driverTotalPoints, teamTotalPoints, personalityAveragePoints, nil)
	simulateChampionship := types.NewSimulateChampionship(championship.Name, types.TotalStatistics{}, *lastChampionshipStatistics)
	w.WriteHeader(http.StatusOK)
	serial, _ := json.Marshal(simulateChampionship)
	w.Write(serial)
}
