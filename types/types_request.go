package types

type DriverRank struct {
	Rank        int            `json:"rank"`
	FirstName   string         `json:"firstname"`
	Lastname    string         `json:"lastname"`
	Points      int            `json:"points"`
	Personality map[string]int `json:"personality"`
}

func NewDriverRank(rank int, firstname, lastname string, pts int, personality map[string]int) *DriverRank {
	return &DriverRank{Rank: rank, FirstName: firstname, Lastname: lastname, Points: pts, Personality: personality}
}

type PersonalityInfo struct {
	IdDriver    string         `json:"idDriver"`
	Lastname    string         `json:"lastname"`
	Personality map[string]int `json:"personality"`
}

type UpdatePersonalityInfo struct {
	IdDriver    string         `json:"idDriver"`
	Personality map[string]int `json:"personnalities"`
}
