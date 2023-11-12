package types

type Driver struct {
	Id                 string       // Driver ID
	Firstname          string       // Firstname
	Lastname           string       // Lastname
	Level              int          // Level of the driver, in [1, 10]
	Country            string       // Country
	Team               *Team        // Team
	Personnality       Personnality // Personnality
	ChampionshipPoints int          // Points in the current champonship

}

func NewDriver(id string, firstname string, lastname string, level int, country string, team *Team, personnality Personnality) *Driver {
	return &Driver{
		Id:                 id,
		Firstname:          firstname,
		Lastname:           lastname,
		Level:              level,
		Country:            country,
		Team:               team,
		Personnality:       personnality,
		ChampionshipPoints: 0,
	}
}
