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

}

type DriverInRace struct {
	Driver   *Driver  //Pilote lui même
	Position *Portion //Position du pilote
	NbLaps   int      //Nombre de tours effectués

	//Pour l'implémentation:
	// - On a un channel pour recevoir les messages de l'environnement
	// - On a un channel pour envoyer les messages à l'environnement
	ChanEnvIn  chan int
	ChanEnvOut chan int
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

func NewDriverInRace(driver *Driver, position *Portion) *DriverInRace {
	cout := make(chan int)
	return &DriverInRace{
		Driver:     driver,
		Position:   position,
		NbLaps:     0,
		ChanEnvOut: cout,
	}
}

func SliceOfDrivers(teams []*Team, portionDepart *Portion) []*DriverInRace {
	res := make([]*DriverInRace, 0)
	for _, team := range teams {
		for _, driver := range team.Drivers {
			d := NewDriverInRace(&driver, portionDepart)
			res = append(res, d)
		}
	}
	return res
}

func ShuffleDrivers(drivers []*DriverInRace) []*DriverInRace {
	rand.Shuffle(len(drivers), func(i, j int) {
		drivers[i], drivers[j] = drivers[j], drivers[i]
	})
	return drivers
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

func (d *DriverInRace) Start(canVersRace chan int, canDepuisRace chan int, position *Portion, nbLaps int) {

	//On stocke le chanel entrant
	d.ChanEnvIn = canDepuisRace
	d.ChanEnvOut = canVersRace

	for {
		//On attend que l'environnement nous dise qu'on peut jouer
		<-d.ChanEnvIn

		//On joue
		//TOD

		//On vérifie si on a fini la course
		if d.NbLaps == nbLaps {
			//On envoie la fin à l'environnement
			d.ChanEnvOut <- 1
			return
		} else {
			d.ChanEnvOut <- 0
		}

	}

}
