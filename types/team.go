package types

type Team struct {
	Id      string   // Team ID
	Name    string   // Name
	Drivers []Driver // Pilotes

}

func (t *Team) CalcChampionshipPoints() int {
	var res int
	for _, driver := range t.Drivers {
		res += driver.ChampionshipPoints
	}
	return res
}
