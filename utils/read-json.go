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

	var teams []types.Team

	// Décoder le fichier JSON dans la structure de données
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&teams); err != nil {
		fmt.Println("Erreur lors de la lecture du fichier JSON :", err)
		return nil, err
	}

	return teams, nil
}
