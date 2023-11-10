package types

import "sort"

type Championship struct {
	Id       string     // Championship ID
	Name     string     // Name
	Circuits []*Circuit // Set of circuits that compose the championship. Defined at the creation of the championship
	Races    []*Race    // Array of Races, filled during the championship
	Teams    []*Team    // Set of teams
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
