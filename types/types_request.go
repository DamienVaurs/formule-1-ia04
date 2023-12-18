package types

type DriverRank struct {
	Rank         int            `json:"rank"`
	FirstName    string         `json:"firstname"`
	Lastname     string         `json:"lastname"`
	Points       int            `json:"points"`
	Personnality map[string]int `json:"personnality"`
}

func NewDriverRank(rank int, firstname, lastname string, pts int, personnality map[string]int) *DriverRank {
	return &DriverRank{Rank: rank, FirstName: firstname, Lastname: lastname, Points: pts, Personnality: personnality}
}
