package types

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
)

type Driver struct {
	Id                 string       // Driver ID
	Firstname          string       // Firstname
	Lastname           string       // Lastname
	Level              int          // Level of the driver, in [1, 10]
	Country            string       // Country
	Personnality       Personnality // Personnality
	ChampionshipPoints int          // Points in the current champonship
}

type DriverInRace struct {
	Driver        *Driver      //Pilote lui même
	Position      *Portion     //Position du pilote
	NbLaps        int          //Nombre de tours effectués
	Status        DriverStatus //Status du pilote
	IsPitStop     bool         // PitStop --> true if the driver is in pitstop
	TimeWoPitStop int          // Time without pitstop --> increments at each step
	//Pour l'implémentation:
	// - On a un channel pour recevoir et envoyer les actions & l'environnement
	ChanEnv      chan Action
	PitstopSteps int // Nombre de steps bloqué en pitstop
}

//Actions d'un pilote

type Action int

const (
	TRY_OVERTAKE Action = iota
	NOOP
	CONTINUE
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
		Personnality:       personnality,
		ChampionshipPoints: 0,
	}
}

func NewDriverInRace(driver *Driver, position *Portion, channel chan Action) *DriverInRace {
	return &DriverInRace{
		Driver:        driver,
		Position:      position,
		NbLaps:        0,
		ChanEnv:       channel,
		Status:        RACING,
		TimeWoPitStop: 0,
		PitstopSteps:  0,
	}
}

// Fonction pour tester si un pilote réussit une portion sans se crasher
func (d *DriverInRace) PortionSuccess() bool {
	// Pour le moment on prend en compte le niveau du pilote, la difficulté de la portion et l'usure des pneus
	portion := d.Position
	probaReussite := 80
	probaReussite += d.Driver.Level * 2
	probaReussite -= portion.Difficulty * 2
	probaReussite -= d.TimeWoPitStop

	var dice int = rand.Intn(99) + 1

	return dice <= probaReussite
}
func MakeSliceOfDriversInRace(teams []*Team, portionDepart *Portion, mapChan sync.Map) ([]*DriverInRace, error) {
	res := make([]*DriverInRace, 0)
	for _, team := range teams {
		for _, driver := range team.Drivers {
			dtamp := driver //nécessaire, sinon n'utilise l'adresse que d'un membre de l'équipe
			c, ok := mapChan.Load(dtamp.Id)
			if !ok {
				return nil, fmt.Errorf("error while creating driver in race : %s", driver.Id)
			}
			d := NewDriverInRace(&dtamp, portionDepart, c.(chan Action))
			res = append(res, d)
		}
	}
	return res, nil
}

func ShuffleDrivers(drivers []*DriverInRace) []*DriverInRace {
	rand.Shuffle(len(drivers), func(i, j int) {
		drivers[i], drivers[j] = drivers[j], drivers[i]
	})
	return drivers
}

func (d *DriverInRace) PitStop() bool {

	probaPitStop := 0

	// On regarde si on doit faire un pitstop
	probaPitStop += d.TimeWoPitStop * 2

	var dice int = rand.Intn(99) + 1

	if dice < probaPitStop {
		d.PitstopSteps = 3
		d.TimeWoPitStop = 0
		return true
	}

	return false
}

func (d *DriverInRace) Overtake(otherDriver *DriverInRace) (reussite bool, crashedDrivers []*DriverInRace) {

	// Si l'autre pilote est en pitstop, on est sûr de doubler
	if otherDriver.Status == PITSTOP {
		return true, []*DriverInRace{}
	}

	if d.Status == PITSTOP {
		return false, []*DriverInRace{}
	}

	probaDoubler := 75

	if d.Driver.Level > otherDriver.Driver.Level {
		probaDoubler += 10
	} else if d.Driver.Level < otherDriver.Driver.Level {
		probaDoubler -= 10
	}

	portion := d.Position

	// Pour le moment on prend en compte le niveaus des pilotes et la "difficulté" de la portion
	probaDoubler -= portion.Difficulty * 2

	var dice int = rand.Intn(99) + 1
	fmt.Println("Dice : ", dice, " probaDoubler : ", probaDoubler)

	// Si on est en dessous de probaDoubler, on double
	if dice <= probaDoubler {
		return true, []*DriverInRace{}
	}

	// Sinon, on regarde si on crash

	// Ici on a un échec critique, les deux pilotes crashent
	if dice >= 99 {
		return false, []*DriverInRace{d, otherDriver}
	}

	// Ici, un seul pilote crash, on tire au sort lequel
	if dice >= 95 {
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
	for i := range p.DriversOn {
		if p.DriversOn[i] == d && d.Status != PITSTOP {
			if len(p.DriversOn) > i+1 && p.DriversOn[i+1] != nil {
				return p.DriversOn[i+1], nil
			} else {
				return nil, nil
			}
		}
	}
	return nil, fmt.Errorf("Driver %s (%s, crashé si =1 : %d) who want to overtake is not found on portion %s", d.Driver.Id, d.Driver.Lastname, d.Status, p.Id)
}

// Fonction pour décider si on veut ESSAYER de doubler ou non
func (d *DriverInRace) OvertakeDecision(driverToOvertake *DriverInRace) (bool, error) {

	toOvertake, err := d.DriverToOvertake()
	if err != nil {
		return false, err
	}
	if toOvertake != nil {
		//On décide si on veut doubler
		//TODO modifier :

		// Si le pilote est en pitstop, on choisit de doubler
		if driverToOvertake.Status == PITSTOP {
			return true, nil
		}

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
		//On attend que l'environnement nous dise qu'on peut prendre une décision
		<-d.ChanEnv
		if d.Status == ARRIVED || d.Status == CRASHED {
			return
		}
		//On décide

		// On regarder si on doit faire un pitstop

		pitstop := d.PitStop()
		if pitstop {
			d.ChanEnv <- NOOP
			d.Status = PITSTOP
			continue
		}

		if d.Status == PITSTOP && d.PitstopSteps != 0 {
			d.ChanEnv <- NOOP
			d.PitstopSteps--
			continue
		} else if d.Status == PITSTOP && d.PitstopSteps == 0 {
			d.Status = RACING
		}

		//On regarde si on peut doubler
		toOvertake, err := d.DriverToOvertake()
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
				d.ChanEnv <- CONTINUE
			}

		} else {
			//Si pas de possibilité de doubler, on ne fait rien
			d.ChanEnv <- CONTINUE
		}
		//On vérifie si on a fini la course
		if d.NbLaps == nbLaps {
			return
		}
	}
}
