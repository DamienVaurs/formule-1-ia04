package restserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"gitlab.utc.fr/vaursdam/formule-1-ia04/types"
)

var i = 0
var raceStatistics *types.SimulateRace = &types.SimulateRace{}
var championship *types.Championship
var firstSimulation = true

func (rsa *RestServer) startRaceSimulation(w http.ResponseWriter, r *http.Request) {
	// var driverTotalPoints []*types.DriverTotalPoints
	// var teamTotalPoints []*types.TeamTotalPoints
	// var personalityAveragePoints []*types.PersonalityAveragePoints
	// var personnalityAverage map[string]map[int]float64

	// driversRankTab := make([]*types.DriverTotalPoints, 0)

	if firstSimulation { //Initialise le championnat si premier lancement de la simulation course-par-course
		championship = types.NewChampionship(nextChampionship, nextChampionship, rsa.pointTabCircuit, rsa.pointTabTeam)
	}

	if r.Method != "GET" {
		return
	}
	fmt.Println("GET /simulateRace")

	//On simule la course i

	//Création de la course
	var id = championship.Circuits[i].Name + " " + championship.Name
	raceStatistics.Championship = championship.Name
	raceStatistics.Race = championship.Circuits[i].Name

	var date = time.Now()
	if i != 0 {
		date = championship.Races[i-1].Date.AddDate(0, 0, 14)
	}
	var meteo = championship.Circuits[i].GenerateMeteo()
	new_Race := types.NewRace(id, championship.Circuits[i], date, championship.Teams, meteo)

	//Simulation de la course
	pointsMap, err := new_Race.SimulateRace()
	if err != nil {
		log.Printf("Erreur simulation cours %s : %s\n", new_Race.Id, err.Error())
	}

	// Ajout points gagnés au points du championnat
	for indT := range championship.Teams {
		for indD := range championship.Teams[indT].Drivers {
			championship.Teams[indT].Drivers[indD].ChampionshipPoints += pointsMap[championship.Teams[indT].Drivers[indD].Id]
		}
	}

	// Points des pilotes pour la course
	driversRankTab := make([]*types.DriverTotalPoints, 0)
	for _, driver := range new_Race.FinalResult {
		driverRank := types.NewDriverTotalPoints(driver.Lastname, pointsMap[driver.Id])
		driversRankTab = append(driversRankTab, driverRank)
	}
	raceStatistics.RaceStatistics.DriversTotalPoints = driversRankTab

	// Team points for current race
	teamsRankTab := make([]*types.TeamTotalPoints, 0)
	for _, team := range new_Race.Teams {
		teamPoints := 0
		for _, driver := range team.Drivers {
			teamPoints += pointsMap[driver.Id]
		}
		teamRank := types.NewTeamTotalPoints(team.Name, teamPoints)
		teamsRankTab = append(teamsRankTab, teamRank)
	}
	raceStatistics.RaceStatistics.TeamsTotalPoints = teamsRankTab

	w.WriteHeader(http.StatusOK)
	serial, _ := json.Marshal(raceStatistics)
	w.Write(serial)
}
