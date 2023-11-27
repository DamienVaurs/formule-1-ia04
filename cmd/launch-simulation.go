package main

import (
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

	//On a les équipes et les circuits, on lance la simulation
	championship := types.NewChampionship("2023", "Championship 1", pointTabCircuit, pointTabTeam)
	s := simulator.NewSimulator([]types.Championship{*championship})
	s.LaunchSimulation()
}
