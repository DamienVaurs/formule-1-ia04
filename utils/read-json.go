package utils

import (
	"encoding/json"
	"fmt"
	"os"

	"gitlab.utc.fr/vaursdam/formule-1-ia04/types"
)

const (
	// Path to the JSON file containing the circuits
	CIRCUITS_PATH = "instances/circuits/inst-circuits.json"
	// Path to the JSON file containing the teams
	TEAMS_PATH = "instances/teams/inst-teams.json"
)

func ReadCircuit() ([]types.Circuit, error) {
	// Ouvrir et lire le fichier JSON
	file, err := os.Open(CIRCUITS_PATH)
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture du fichier :", err)
		return nil, err
	}
	defer file.Close()

	var circuits []types.Circuit

	// Décoder le fichier JSON dans la structure de données
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&circuits); err != nil {
		fmt.Println("Erreur lors de la lecture du fichier JSON :", err)
		return nil, err
	}

	// On spécifie le type de chaque portion
	for i := 0; i < len(circuits); i++ {
		circuits[i].Id = fmt.Sprintf("circuit-%d", i)
		for j := 0; j < len(circuits[i].Portions); j++ {
			if len(circuits[i].Portions[j].Id) < len("turn") {
			} else if circuits[i].Portions[j].Id[:len("turn")] == "turn" {
				circuits[i].Portions[j].Type = types.TURN
			} else {
				circuits[i].Portions[j].Type = types.STRAIGHT
			}
		}
	}

	return circuits, nil
}

func ReadTeams() ([]types.Team, error) {
	// Ouvrir et lire le fichier JSON
	file, err := os.Open(TEAMS_PATH)
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture du fichier :", err)
		return nil, err
	}
	defer file.Close()

	teams := make([]types.Team, 0)

	// Décoder le fichier JSON dans la structure de données
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&teams); err != nil {
		fmt.Println("Erreur lors de la lecture du fichier JSON :", err)
		return nil, err
	}
	//Ajout d'Id aux pilotes et aux team
	for i, team := range teams {
		teams[i].Id = fmt.Sprintf("team-%d", i)
		for j := range team.Drivers {
			teams[i].Drivers[j].Id = fmt.Sprintf("driver-%d-%d", i, j)
		}
	}

	return teams, nil
}
