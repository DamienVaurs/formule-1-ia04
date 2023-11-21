package types

import (
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

func (r *Race) SimulateRace() {
	//On crée les instances des pilotes en course

	var drivers = SliceOfDriversInRace(r.Teams, &(r.Circuit.Portions[0]))
	//On lance les agents pilotes
	for _, driver := range drivers {
		c, ok := r.MapChan.Load(driver.Driver.Id)
		if ok != true {
			log.Printf("Error while loading channel for driver %s\n", driver.Driver.Id)
		}
		go driver.Start(c.(chan Action), driver.Position, driver.NbLaps)
	}

	var nbFinish = 0
	var nbDrivers = len(r.Teams) * 2
	decisionMap := make(map[*DriverInRace]Action, nbDrivers)

	//On simule tant que tous les pilotes n'ont pas fini la course
	for nbFinish < nbDrivers {
		//Chaque pilote, dans un ordre aléatoire, réalise les tests sur la proba de dépasser etc...
		drivers = ShuffleDrivers(drivers)
		for _, driver := range drivers {
			//On débloque le pilote qu'il décide de dépasser ou non
			driver.ChanEnv <- 1
		}
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
							driver.Position.RemoveDriverOn(crashedDriver)
							nbFinish++
						}

						if success {
							//On met à jour les positions
							driver.Position.RemoveDriverOn(driver)
							driverToOvertake.Position.AddDriverOn(driver)
						}
					}

				}
			}
		}
	}
}
