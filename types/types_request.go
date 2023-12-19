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
	Personality map[string]int `json:"personality"`
}

type SimulateChampionship struct {
	LastChampionship           string                     `json:"lastChampionship"`
	LastChampionshipStatistics LastChampionshipStatistics `json:"lastChampionshipStatistics"`
}

type LastChampionshipStatistics struct {
	DriversTotalPoints       []DriverTotalPoints      `json:"driversTotalPoints"`
	TeamsTotalPoints         []TeamTotalPoints        `json:"teamsTotalPoints"`
	PersonalitiesTotalPoints []PersonalityTotalPoints `json:"personnalityTotalPoints"`
	NbCrashsPersonnality     []NbCrashsPersonnality   `json:"nbCrashsPersonnality"`
}

type DriverTotalPoints struct {
	Driver      string `json:"driver"`
	TotalPoints int    `json:"totalPoints"`
}

type TeamTotalPoints struct {
	Team        string `json:"team"`
	TotalPoints int    `json:"totalPoints"`
}

type PersonalityTotalPoints struct {
	Personality map[string]int `json:"personnality"`
	TotalPoints int            `json:"totalPoints"`
}

type NbCrashsPersonnality struct {
	Personality map[string]int `json:"personnality"`
	NbCrash     int            `json:"nbCrash"`
}
