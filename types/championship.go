package types

import (
	"log"
	"sort"
)

type Championship struct {
	Id       string     // Championship ID
	Name     string     // Name
	Circuits []*Circuit // Set of circuits that compose the championship. Defined at the creation of the championship
	Races    []Race     // Array of Races, filled during the championship
	Teams    []*Team    // Set of teams
}

func NewChampionship(id string, name string, circuits []*Circuit, teams []*Team) *Championship {

	c := make([]*Circuit, len(circuits))
	copy(c, circuits)

	t := make([]*Team, len(teams))
	copy(t, teams)

	r := make([]Race, len(circuits))

	return &Championship{
		Id:       id,
		Name:     name,
		Circuits: c,
		Races:    r,
		Teams:    t,
	}
}

//Remarque : on utilise des pointeurs quand l'objet ne gère pas le cycle de vie des instances

func (c *Championship) CalcTeamRank() []*Team {
	res := make([]*Team, len(c.Teams))
	copy(res, c.Teams)
	sort.Slice(res, func(i, j int) bool {
		return res[i].CalcChampionshipPoints() > res[j].CalcChampionshipPoints()
	})

	return res
}

func (c *Championship) DisplayTeamRank() {
	log.Printf("\n\n====Classement constructeur ====\n")
	teamRank := c.CalcTeamRank()
	for i, team := range teamRank {
		log.Printf("%d : %s : %d points\n", i+1, team.Name, team.CalcChampionshipPoints())
	}
}

func (c *Championship) CalcDriverRank() []*Driver {

	res := make([]*Driver, 0)
	for indT := range c.Teams {
		for indD := range c.Teams[indT].Drivers {
			res = append(res, &c.Teams[indT].Drivers[indD])
		}
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].ChampionshipPoints > res[j].ChampionshipPoints
	})

	return res
}

func (c *Championship) DisplayDriverRank() []*DriverRank {
	log.Printf("\n\n====Classement pilotes ====\n")
	driversRank := c.CalcDriverRank()
	driversRankTab := make([]*DriverRank, 20)
	for i, driver := range driversRank {
		driverRank := NewDriverRank(i+1, driver.Firstname, driver.Lastname, driver.ChampionshipPoints, driver.Personnality.TraitsValue)
		log.Printf("%d : %s %s : %d points\n", i+1, driver.Firstname, driver.Lastname, driver.ChampionshipPoints)
		log.Printf("%v", driver.Personnality.TraitsValue)

		driversRankTab = append(driversRankTab, driverRank)

	}
	return driversRankTab[20:] // Les 20 premiers indices sont nulles
}

func (c *Championship) DisplayPersonalityRepartition() {
	log.Printf("\n\n====Répartition des personnalités ====\n")
	driverRank := c.CalcDriverRank()

	aggressivity_value_5 := 0
	aggressivity_value_4 := 0
	aggressivity_value_3 := 0
	aggressivity_value_2 := 0
	aggressivity_value_1 := 0
	aggressivity_value_0 := 0
	for i, driver := range driverRank {
		if i < 15 {
			switch driver.Personnality.TraitsValue["Aggressivity"] {
			case 0:
				aggressivity_value_0 += 1
			case 1:
				aggressivity_value_1 += 1
			case 2:
				aggressivity_value_2 += 1
			case 3:
				aggressivity_value_3 += 1
			case 4:
				aggressivity_value_4 += 1
			case 5:
				aggressivity_value_5 += 1
			default:
				log.Printf("Value of aggressivity out of range : %d", driver.Personnality.TraitsValue["Aggressivity"])
			}
		}
		if i == 4 {
			log.Printf("Répartition du niveau agressivité du top 5 : \n")
			log.Printf("Agressivité 5 : %d", aggressivity_value_5)
			log.Printf("Agressivité 4 : %d", aggressivity_value_4)
			log.Printf("Agressivité 3 : %d", aggressivity_value_3)
			log.Printf("Agressivité 2 : %d", aggressivity_value_2)
			log.Printf("Agressivité 1 : %d", aggressivity_value_1)
			log.Printf("Agressivité 0 : %d", aggressivity_value_0)
		}
		if i == 9 {
			log.Printf("Répartition du niveau agressivité du top 10 : \n")
			log.Printf("Agressivité 5 : %d", aggressivity_value_5)
			log.Printf("Agressivité 4 : %d", aggressivity_value_4)
			log.Printf("Agressivité 3 : %d", aggressivity_value_3)
			log.Printf("Agressivité 2 : %d", aggressivity_value_2)
			log.Printf("Agressivité 1 : %d", aggressivity_value_1)
			log.Printf("Agressivité 0 : %d", aggressivity_value_0)
		}
		if i == 14 {
			log.Printf("Répartition du niveau agressivité du top 15 : \n")
			log.Printf("Agressivité 5 : %d", aggressivity_value_5)
			log.Printf("Agressivité 4 : %d", aggressivity_value_4)
			log.Printf("Agressivité 3 : %d", aggressivity_value_3)
			log.Printf("Agressivité 2 : %d", aggressivity_value_2)
			log.Printf("Agressivité 1 : %d", aggressivity_value_1)
			log.Printf("Agressivité 0 : %d", aggressivity_value_0)
		}
	}
}
