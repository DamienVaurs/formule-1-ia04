package main

import (
	"fmt"

	"gitlab.utc.fr/vaursdam/formule-1-ia04/utils"
)

func main() {
	c, err := utils.ReadCircuit()
	if err != nil {
		panic(err)
	}
	ind := len(c) - 1

	fmt.Println(c[ind].Name)
	fmt.Println(c[ind].Country)
	fmt.Println(c[ind].MeteoDistribution)
	fmt.Println(c[ind].GenerateMeteo())
	fmt.Println(c[ind].Portions)

	t, err := utils.ReadTeams()
	if err != nil {
		panic(err)
	}

	fmt.Println(t[0].Name)
	/**
	fmt.Println(t[0].Level)
	fmt.Println(t[0].CalcChampionshipPoints())
	fmt.Println(t[0].Drivers[0].Firstname)
	fmt.Println(t[0].Drivers[0].Lastname)
	fmt.Println(t[0].Drivers[0].Level)
	fmt.Println(t[0].Drivers[0].Country)
	fmt.Println(t[0].Drivers[0].Personnality)
	*/

}
