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
	} /*
		ind := len(c) - 1

		fmt.Println(c[ind].Name)
		fmt.Println(c[ind].Country)
		fmt.Println(c[ind].MeteoDistribution)
		fmt.Println(c[ind].GenerateMeteo())
		fmt.Println(c[ind].Portions)
	*/
	t, err := utils.ReadTeams()
	if err != nil {
		panic(err)
	}
	/*
		for _, team := range t {
			for _, driver := range team.Drivers {
				fmt.Printf("%s %s\n", driver.Firstname, driver.Lastname)
			}
		}
	*/
	//fmt.Println(t[0].Name)
	/**
	fmt.Println(t[0].Level)
	fmt.Println(t[0].CalcChampionshipPoints())
	fmt.Println(t[0].Drivers[0].Firstname)
	fmt.Println(t[0].Drivers[0].Lastname)
	fmt.Println(t[0].Drivers[0].Level)
	fmt.Println(t[0].Drivers[0].Country)
	fmt.Println(t[0].Drivers[0].Personnality)
	*/

	//Lancement simulation
	pointTabCircuit := make([]*types.Circuit, len(c))
	for i, circuit := range c {
		tempCircuit := circuit //sans tampon, tous les éléments du tableau contiendront la même adresse
		pointTabCircuit[i] = &tempCircuit
	}
	/*
		for _, circuit := range pointTabCircuit {
			fmt.Println(circuit.Name)
		}*/
	pointTabTeam := make([]*types.Team, len(t))
	for i, team := range t {
		tempTeam := team //sans tampon, tous les éléments du tableau contiendront la même adresse
		pointTabTeam[i] = &tempTeam
	}
	/*
		for _, team := range pointTabTeam {
			fmt.Printf("%s %d\n", team.Name, team.Level)
			for _, driver := range team.Drivers {
				fmt.Printf("	%s %s\n", driver.Firstname, driver.Lastname)
			}
		}*/

	//On a les équipes et les circuits, on lance la simulation
	championship := types.NewChampionship("2023", "Championship 1", pointTabCircuit, pointTabTeam)
	s := simulator.NewSimulator([]types.Championship{*championship})
	s.LaunchSimulation()
}
