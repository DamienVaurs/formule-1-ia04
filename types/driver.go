package types

import "math/rand"

type Driver struct {
	Id                 string       // Driver ID
	Firstname          string       // Firstname
	Lastname           string       // Lastname
	Level              int          // Level of the driver, in [1, 10]
	Country            string       // Country
	Team               *Team        // Team
	Personnality       Personnality // Personnality
	ChampionshipPoints int          // Points in the current champonship
	IsPitStop          bool         // PitStop --> true if the driver is in pitstop
	TimeWoPitStop      int          // Time without pitstop --> increments at each step

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

// Fonction pour tester si un pilote réussit une portion sans se crasher
func (d *Driver) PortionSuccess(portion *Portion) bool {
	// Pour le moment on prend en compte le niveau du pilote et la "difficulté" de la portion
	probaReussite := 80
	probaReussite += d.Level * 2
	probaReussite -= portion.CrashProbability * 2

	var dice int = rand.Intn(99) + 1

	return dice <= probaReussite
}

func (d *Driver) Overtake(otherDriver *Driver, portion *Portion) (reussite bool, crashedDrivers []*Driver) {

	probaDoubler := 75

	bonus := 0

	if d.Level > otherDriver.Level {
		bonus = 10
	} else if d.Level < otherDriver.Level {
		bonus = -10
	} else {
		bonus = 0
	}

	// Pour le moment on prend en compte le niveau des pilotes et la "difficulté" de la portion
	probaDoubler += bonus
	probaDoubler -= portion.CrashProbability * 7

	var dice int = rand.Intn(99) + 1

	// Si on est en dessous de probaDoubler, on double
	if dice < probaDoubler {
		return true, []*Driver{}
	}

	// Sinon, on regarde si on crash

	// Ici on a un échec critique, les deux pilotes crashent
	if dice > 95 {
		return false, []*Driver{d, otherDriver}
	}

	// Ici, un seul pilote crash, on tire au sort lequel
	if dice > 90 {
		if dice%2 == 0 {
			return false, []*Driver{d}
		} else {
			return false, []*Driver{otherDriver}
		}
	}

	// Dans le cas par défaut, le doublement est échoué mais aucun crash n'a lieu
	return false, []*Driver{}

}

func (d *Driver) PitStop() {
	d.IsPitStop = true
	// Envoyer sur un channel vers le contrôleur de jeu que le pilote est en pitstop pendant x steps
	// On attends ensuite de recevoir un message sur le channel comme quoi le pitstop est terminé ?
	d.IsPitStop = false
	d.TimeWoPitStop = 0
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

func (d *Driver) Overtake(otherDriver *Driver, portion *Portion) (reussite bool, crashedDrivers []*Driver) {

	probaDoubler := 75

	bonus := 0

	if d.Level > otherDriver.Level {
		bonus = 10
	} else if d.Level < otherDriver.Level {
		bonus = -10
	} else {
		bonus = 0
	}

	// Pour le moment on prend en compte le niveau des pilotes et la "difficulté" de la portion
	probaDoubler += bonus
	probaDoubler -= portion.Difficulty * 7

	var dice int = rand.Intn(99) + 1

	// Si on est en dessous de probaDoubler, on double
	if dice < probaDoubler {
		return true, []*Driver{}
	}

	// Sinon, on regarde si on crash

	// Ici on a un échec critique, les deux pilotes crashent
	if dice > 95 {
		return false, []*Driver{d, otherDriver}
	}

	// Ici, un seul pilote crash, on tire au sort lequel
	if dice > 90 {
		if dice%2 == 0 {
			return false, []*Driver{d}
		} else {
			return false, []*Driver{otherDriver}
		}
	}

	// Dans le cas par défaut, le doublement est échoué mais aucun crash n'a lieu
	return false, []*Driver{}

}
