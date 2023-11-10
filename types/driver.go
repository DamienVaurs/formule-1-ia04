package types

type Driver struct {
	DriverId           string       // Driver ID
	Name               string       // Name
	Country            string       // Country
	Team               *Team        // Team
	Personnality       Personnality // Personnality
	ChampionshipPoints int          // Points in the current champonship

}
