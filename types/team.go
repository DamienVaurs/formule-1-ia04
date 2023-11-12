package types

type Team struct {
	Id      string   // Team ID
	Name    string   // Name
	Drivers []Driver // Pilotes
	Level   int      // Level of the team car, in [1, 10]

}

func NewTeam(id string, name string, drivers []Driver, level int) *Team {
	return &Team{
		Id:      id,
		Name:    name,
		Drivers: drivers,
		Level:   level,
	}
}

func (t *Team) CalcChampionshipPoints() int {
	var res int
	for _, driver := range t.Drivers {
		res += driver.ChampionshipPoints
	}
	return res
}
