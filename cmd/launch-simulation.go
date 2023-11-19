package main

import (
	"gitlab.utc.fr/vaursdam/formule-1-ia04/utils"
)

func main() {
	c, err := utils.ReadCircuit()
	if err != nil {
		panic(err)
	}
	/**
	fmt.Println(c[0].Name)
	fmt.Println(c[0].Country)
	fmt.Println(c[0].MeteoDistribution)
	fmt.Println(c[0].GenerateMeteo())
	fmt.Println(c[0].Portions)
	*/
	t, err := utils.ReadTeams()
	if err != nil {
		panic(err)
	}
	/**
	fmt.Println(t[0].Name)
	fmt.Println(t[0].Level)
	fmt.Println(t[0].CalcChampionshipPoints())
	fmt.Println(t[0].Drivers[0].Firstname)
	fmt.Println(t[0].Drivers[0].Lastname)
	fmt.Println(t[0].Drivers[0].Level)
	fmt.Println(t[0].Drivers[0].Country)
	fmt.Println(t[0].Drivers[0].Personnality)
	*/

}
