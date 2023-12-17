package types

type DriverRank struct {
	Rank      int    `json:"rank"`
	FirstName string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Points    int    `json:"points"`
}

func NewDriverRank(rank int, firstname, lastname string, pts int) *DriverRank {
	return &DriverRank{Rank: rank, FirstName: firstname, Lastname: lastname, Points: pts}
}
