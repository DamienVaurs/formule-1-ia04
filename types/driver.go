package types

import (
	"log"
	"math/rand"
)

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
	Driver   *Driver      //Pilote lui même
	Position *Portion     //Position du pilote
	NbLaps   int          //Nombre de tours effectués
	Status   DriverStatus //Status du pilote
	//Pour l'implémentation:
	// - On a un channel pour recevoir et envoyer les actions & l'environnement
	ChanEnv chan Action
}

//Actions d'un pilote

type Action int

const (
	TRY_OVERTAKE Action = iota
	NOOP
)

type DriverStatus int

const (
	RACING DriverStatus = iota
	CRASHED
	ARRIVED
	PITSTOP
)

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
	c := make(chan Action)
	return &DriverInRace{
		Driver:   driver,
		Position: position,
		NbLaps:   0,
		ChanEnv:  c,
		Status:   RACING,
	}
}

func SliceOfDriversInRace(teams []*Team, portionDepart *Portion) []*DriverInRace {
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

func (d *DriverInRace) Overtake(otherDriver *DriverInRace) (reussite bool, crashedDrivers []*DriverInRace) {

	probaDoubler := 75

	bonus := 0

	if d.Driver.Level > otherDriver.Driver.Level {
		bonus = 10
	} else if d.Driver.Level < otherDriver.Driver.Level {
		bonus = -10
	} else {
		bonus = 0
	}

	portion := d.Position

	// Pour le moment on prend en compte le niveau des pilotes et la "difficulté" de la portion
	probaDoubler += bonus
	probaDoubler -= portion.Difficulty * 7

	var dice int = rand.Intn(99) + 1

	// Si on est en dessous de probaDoubler, on double
	if dice < probaDoubler {
		return true, []*DriverInRace{}
	}

	// Sinon, on regarde si on crash

	// Ici on a un échec critique, les deux pilotes crashent
	if dice > 95 {
		return false, []*DriverInRace{d, otherDriver}
	}

	// Ici, un seul pilote crash, on tire au sort lequel
	if dice > 90 {
		if dice%2 == 0 {
			return false, []*DriverInRace{d}
		} else {
			return false, []*DriverInRace{otherDriver}
		}
	}

	// Dans le cas par défaut, le doublement est échoué mais aucun crash n'a lieu
	return false, []*DriverInRace{}

}

// Fonction pourdécider si on veut ESSAYER de doubler ou non
func (d *DriverInRace) OvertakeDecision(driverToOvertake *DriverInRace) (bool, error) {

	toOvertake, err := d.Position.DriverToOvertake(d)
	if err != nil {
		return false, err
	}
	if toOvertake != nil {
		//On décide si on veut doubler
		//TODO modifier :
		var dice = rand.Int() % 2
		if dice == 0 {
			return true, nil
		} else {
			return false, nil
		}
	}
	return false, nil
}

func (d *DriverInRace) Start(raceChan chan Action, position *Portion, nbLaps int) {
	log.Printf("Lancement du pilote driver %s %s\n", d.Driver.Firstname, d.Driver.Lastname)
	//On stocke le chanel
	d.ChanEnv = raceChan

	for {
		//On attend que l'environnement nous dise qu'on peut prendre une décision
		<-d.ChanEnv

		//On décide
		//On regarde si on peut doubler
		toOvertake, err := position.DriverToOvertake(d)
		if err != nil {
			log.Printf("Error while getting the driver to overtake : %s\n", err)
		}
		if toOvertake != nil {
			//On décide si on veut doubler
			decision, err := d.OvertakeDecision(toOvertake)
			if err != nil {
				log.Printf("Error while getting the decision to overtake : %s\n", err)
			}
			if decision {
				//On envoie la décision à l'environnement
				d.ChanEnv <- TRY_OVERTAKE
			} else {
				d.ChanEnv <- NOOP
			}

			//On vérifie si on a fini la course
			if d.NbLaps == nbLaps {
				return
			}

		}

	}
}
