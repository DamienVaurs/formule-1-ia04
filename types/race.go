package types

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Race struct {
	Id             string      // Race ID
	Circuit        *Circuit    // Circuit
	Date           time.Time   // Date
	Teams          []*Team     // Set of teams
	MeteoCondition Meteo       // Meteo condition
	FinalResult    []*Driver   // Final result, drivers rank from 1st to last
	HighLigths     []Highlight // Containes all what happend during the race

	//Pour l'implémentation:
	MapChan sync.Map //Map qui contient les channels de communication entre les pilotes et l'environnement
}

func NewRace(id string, circuit *Circuit, date time.Time, teams []*Team, meteo Meteo) *Race {

	d := make([]*Team, len(teams))
	copy(d, teams)

	f := make([]*Driver, len(teams)*2) //car 2 drivers par team

	h := make([]Highlight, 0)

	m := sync.Map{}
	for _, t := range teams {
		for _, d := range t.Drivers {
			m.Store(d.Id, make(chan Action))
		}
	}

	return &Race{
		Id:             id,
		Circuit:        circuit,
		Date:           date,
		Teams:          d,
		MeteoCondition: meteo,
		FinalResult:    f,
		HighLigths:     h,
		MapChan:        m,
	}
}

func (r *Race) SimulateRace() error {
	log.Printf("	Lancement d'une nouvelle course : %s...\n", r.Id)
	//On crée les instances des pilotes en course

	drivers, err := MakeSliceOfDriversInRace(r.Teams, &(r.Circuit.Portions[0]), r.MapChan)
	if err != nil {
		return err
	}

	//On lance les agents pilotes
	for _, driver := range drivers {
		go driver.Start(driver.Position, driver.NbLaps)
	}
	var nbFinish = 0
	var nbDrivers = len(r.Teams) * 2
	decisionMap := make(map[*DriverInRace]Action, nbDrivers)

	//On simule tant que tous les pilotes n'ont pas fini la course
	for nbFinish < nbDrivers {
		//Chaque pilote, dans un ordre aléatoire, réalise les tests sur la proba de dépasser etc...
		drivers = ShuffleDrivers(drivers)
		fmt.Println("Débloquage des go routines...")

		for _, driver := range drivers {
			//On débloque le pilote pour qu'il prenne une décision
			fmt.Println("Envoie de déblocage à : " + driver.Driver.Lastname)
			driver.ChanEnv <- 1
		}
		fmt.Println("Les go routines sont débloquées")
		// On récupère les décisions des pilotes
		for _, driver := range drivers {
			decisionMap[driver] = <-driver.ChanEnv
		}

		//On traite les décisions et on met à jour les positions des pilotes
		for driver, decision := range decisionMap {
			switch decision {
			case TRY_OVERTAKE:
				//On vérifie si le pilote peut bien dépasser
				driverToOvertake, err := driver.Position.DriverToOvertake(driver)
				if err != nil {
					log.Printf("Error while getting driver to overtake: %s\n", err)
				}
				if driverToOvertake != nil {
					//On vérifie si le pilote a réussi son dépassement
					success, crashedDrivers := driver.Overtake(driverToOvertake)
					if crashedDrivers != nil {
						//On supprime les pilotes crashés
						for _, crashedDriver := range crashedDrivers {
							crashedDriver.Status = CRASHED
							driver.Position.RemoveDriverOn(crashedDriver)
							nbFinish++
						}

						if success {
							//On met à jour les positions
							driver.Position.SwapDrivers(driver, driverToOvertake)
						}
					}
				}
			case NOOP:
				//On ne fait rien

			}
		}

		//On fait avancer tout les pilotes n'ayant pas fini la course et n'étant pas crashés
		newDriversOnPortion := make([][]*DriverInRace, len(r.Circuit.Portions)) //stocke les nouvelles positions des pilotes
		for i, portion := range r.Circuit.Portions {
			for _, driver := range portion.DriversOn {
				if driver.Status != CRASHED && driver.Status != ARRIVED {
					driver.Position = &(r.Circuit.Portions[(i+1)%len(r.Circuit.Portions)]) //Usage du modulo pour gérer le cas où on est sur la dernière portion
					if i == len(r.Circuit.Portions)-1 {
						//Si on a fait un tour
						driver.NbLaps++
						if driver.NbLaps == r.Circuit.NbLaps {
							//Si on a fini la course, on enlève le pilote du circuit et on le met dans le classement
							driver.Status = ARRIVED
							nbFinish++
							r.FinalResult = append(r.FinalResult, driver.Driver)
						}
					}
				}
				if driver.Status != CRASHED && driver.Status != ARRIVED {
					newDriversOnPortion[(i+1)%len(r.Circuit.Portions)] = append(newDriversOnPortion[(i+1)%len(r.Circuit.Portions)], driver)
				}
			}
		}
		//On met à jour les positions des pilotes
		for i, portion := range r.Circuit.Portions {
			portion.DriversOn = newDriversOnPortion[i]
		}
	}
	return nil
}
