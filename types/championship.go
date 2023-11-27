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

func (c *Championship) DisplayDriverRank() {
	log.Printf("\n\n====Classement pilotes ====\n")
	driverRank := c.CalcDriverRank()
	for i, driver := range driverRank {
		log.Printf("%d : %s %s : %d points\n", i+1, driver.Firstname, driver.Lastname, driver.ChampionshipPoints)
	}
}
