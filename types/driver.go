package types

import (
	"fmt"
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
	IsPitStop          bool         // PitStop --> true if the driver is in pitstop
	TimeWoPitStop      int          // Time without pitstop --> increments at each step

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

// Fonction pour tester si un pilote réussit une portion sans se crasher
func (d *Driver) PortionSuccess(portion *Portion) bool {
	// Pour le moment on prend en compte le niveau du pilote et la "difficulté" de la portion
	probaReussite := 80
	probaReussite += d.Level * 2
	probaReussite -= portion.Difficulty * 2

	var dice int = rand.Intn(99) + 1

	return dice <= probaReussite
}

func (d *Driver) PitStop() {
	d.IsPitStop = true
	// Envoyer sur un channel vers le contrôleur de jeu que le pilote est en pitstop pendant x steps
	// On attends ensuite de recevoir un message sur le channel comme quoi le pitstop est terminé ?
	d.IsPitStop = false
	d.TimeWoPitStop = 0
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

func (d *DriverInRace) DriverToOvertake() (*DriverInRace, error) {
	p := d.Position
	//fmt.Println("GAGA", p.Id)
	for i := range p.DriversOn {
		if p.DriversOn[i] == d {
			if len(p.DriversOn) > i+1 && p.DriversOn[i+1] != nil {
				return p.DriversOn[i+1], nil
			} else {
				return nil, nil
			}
		}
	}
	// TODO : vérifier
	return nil, fmt.Errorf("Driver %s (%s) not found on portion %s", d.Driver.Id, d.Driver.Lastname, p.Id)
}

// Fonction pourdécider si on veut ESSAYER de doubler ou non
func (d *DriverInRace) OvertakeDecision(driverToOvertake *DriverInRace) (bool, error) {

	toOvertake, err := d.DriverToOvertake()
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

func (d *DriverInRace) Start(position *Portion, nbLaps int) {
	log.Printf("		Lancement du pilote %s %s...\n", d.Driver.Firstname, d.Driver.Lastname)

	for {
		if d.Status == ARRIVED || d.Status == CRASHED {
			return
		}
		//fmt.Printf("Nb tours fait pour %s : %d\n", d.Driver.Lastname, d.NbLaps) //semble ok
		//On attend que l'environnement nous dise qu'on peut prendre une décision
		//fmt.Println("Attente de l'env : " + d.Driver.Lastname)
		<-d.ChanEnv
		//fmt.Println("Réception de l'env : " + d.Driver.Lastname)

		//On décide
		//On regarde si on peut doubler
		toOvertake, err := d.DriverToOvertake()
		if err != nil {
			log.Printf("Error while getting the driver to overtake : %s\n", err)
		}
		if toOvertake != nil {
			//On décide si on veut doubler
			//fmt.Printf("%s peut essayer de dépasser %s sur %s\n", d.Driver.Lastname, toOvertake.Driver.Lastname, position.Id)

			decision, err := d.OvertakeDecision(toOvertake)
			if err != nil {
				log.Printf("Error while getting the decision to overtake : %s\n", err)
			}
			if decision {
				//fmt.Printf("%s décide de dépasser %s\n", d.Driver.Lastname, toOvertake.Driver.Lastname)
				//On envoie la décision à l'environnement
				d.ChanEnv <- TRY_OVERTAKE
			} else {
				//fmt.Printf("%s décide de NE PAS dépasser %s\n", d.Driver.Lastname, toOvertake.Driver.Lastname)
				d.ChanEnv <- NOOP
			}

		} else {
			//Si pas de possibilité de doubler, on ne fait rien
			//fmt.Printf("%s ne peut dépasser personne sur %s\n", d.Driver.Lastname, position.Id)

			d.ChanEnv <- NOOP
		}
		//On vérifie si on a fini la course
		if d.NbLaps == nbLaps {
			return
		}
		//fmt.Printf("Fin de décision pour %s\n", d.Driver.Lastname)

	}
}
