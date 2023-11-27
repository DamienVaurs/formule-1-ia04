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

}

func NewRace(id string, circuit *Circuit, date time.Time, teams []*Team, meteo Meteo) *Race {

	d := make([]*Team, len(teams))
	copy(d, teams)

	f := make([]*Driver, 0) //car 2 drivers par team

	h := make([]Highlight, 0)

	return &Race{
		Id:             id,
		Circuit:        circuit,
		Date:           date,
		Teams:          d,
		MeteoCondition: meteo,
		FinalResult:    f,
		HighLigths:     h,
	}
}

func (r *Race) SimulateRace() error {
	log.Printf("	Lancement d'une nouvelle course : %s...\n", r.Id)

	//Création du map partagé
	mapChan := sync.Map{}
	for _, t := range r.Teams {
		for _, d := range t.Drivers {
			mapChan.Store(d.Id, make(chan Action))
		}
	}
	//On crée les instances des pilotes en course

	drivers, err := MakeSliceOfDriversInRace(r.Teams, &(r.Circuit.Portions[0]), mapChan)
	if err != nil {
		return err
	}

	//On met tous les agents sur la ligne de départ :
	for _, driver := range drivers {
		driver.Position.AddDriverOn(driver)
	}
	log.Println("Les pilotes sont sur la ligne de départ")
	//On lance les agents pilotes
	log.Println("Début de la course...")
	for _, driver := range drivers {
		go driver.Start(driver.Position, r.Circuit.NbLaps)
	}
	var nbFinish = 0
	var nbDrivers = len(r.Teams) * 2
	decisionMap := make(map[string]Action, nbDrivers)

	//On simule tant que tous les pilotes n'ont pas fini la course
	for nbFinish < nbDrivers {
		log.Printf("\n\n============ NOUVELLE BOUCLE %s : %d ont fini =================\n\n", r.Circuit.Name, nbFinish)
		//time.Sleep(5 * time.Second)
		//Chaque pilote, dans un ordre aléatoire, réalise les tests sur la proba de dépasser etc...
		drivers = ShuffleDrivers(drivers)

		for i := range drivers {
			//On débloque les pilotes en course pour qu'ils prennent une décision
			if drivers[i].Status == CRASHED || drivers[i].Status == ARRIVED {
				continue //Obligatoire car il ne faut attendre que les pilotes qui courent encore
			}
			drivers[i].ChanEnv <- 1
		}
		// On récupère les décisions des pilotes
		for i := range drivers {
			if drivers[i].Status == CRASHED || drivers[i].Status == ARRIVED {
				continue //Obligatoire car il ne faut attendre que les pilotes qui courent encore
			}
			decisionMap[drivers[i].Driver.Id] = <-drivers[i].ChanEnv
		}

		//On traite les décisions et on met à jour les positions des pilotes
		for i := range drivers {
			if drivers[i].Status == CRASHED || drivers[i].Status == ARRIVED {
				continue
			}
			decision := decisionMap[drivers[i].Driver.Id]
			switch decision {
			case TRY_OVERTAKE:
				//On vérifie si le pilote peut bien dépasser
				driverToOvertake, err := drivers[i].DriverToOvertake()
				if err != nil {
					log.Printf("Error while getting driver to overtake: %s\n", err)
				}
				if driverToOvertake != nil {
					//On vérifie si le pilote a réussi son dépassement
					success, crashedDrivers := drivers[i].Overtake(driverToOvertake)
					if len(crashedDrivers) > 0 {
						//On supprime les pilotes crashés
						if len(crashedDrivers) > 1 {
							log.Println("CRASH : Plusieurs pilotes sont rentrés en accident : ", crashedDrivers[0].Driver.Lastname, " et ", crashedDrivers[1].Driver.Lastname)
						} else {
							log.Println("CRASH : Le pilote " + crashedDrivers[0].Driver.Lastname + " a crashé")
						}
						for ind := range crashedDrivers {
							crashedDrivers[ind].Status = CRASHED
							r.FinalResult = append(r.FinalResult, crashedDrivers[ind].Driver) //on l'ajoute au tableau
							drivers[i].Position.RemoveDriverOn(crashedDrivers[ind])
							nbFinish++
						}

						if success {
							//On met à jour les positions
							log.Println("Overtake")
							log.Println("OVERTAKE : Le pilote " + drivers[i].Driver.Lastname + " a réussi son dépassement sur " + driverToOvertake.Driver.Lastname)
							drivers[i].Position.SwapDrivers(drivers[i], driverToOvertake)
						}
					}
				}
			case NOOP:
				//On ne fait rien

			}
		}

		//On fait avancer tout les pilotes n'ayant pas fini la course et n'étant pas crashés
		newDriversOnPortion := make([][]*DriverInRace, len(r.Circuit.Portions)) //stocke les nouvelles positions des pilotes
		for i := range r.Circuit.Portions {
			newDriversOnPortion[(i+1)%len(r.Circuit.Portions)] = make([]*DriverInRace, 0)
			for _, driver := range r.Circuit.Portions[i].DriversOn {
				if driver.Status != CRASHED && driver.Status != ARRIVED {
					//On met à jour le champ position du pilote
					driver.Position = driver.Position.NextPortion
					if i == len(r.Circuit.Portions)-1 {
						//Si on a fait un tour
						driver.NbLaps += 1
						if driver.NbLaps == r.Circuit.NbLaps {
							//Si on a fini la course, on enlève le pilote du circuit et on le met dans le classement
							log.Printf("ARRIVEE : Pilote %s est arrivé!\n", driver.Driver.Lastname)
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
		for i := range r.Circuit.Portions {
			r.Circuit.Portions[i].DriversOn = make([]*DriverInRace, len(newDriversOnPortion[i])) //on écrase l'ancien slice
			copy(r.Circuit.Portions[i].DriversOn, newDriversOnPortion[i])                        //on remplace par le nouveau
		}

	}
	//On affiche le classement
	log.Println("\n\nClassement final :")
	for i := range r.FinalResult {
		log.Printf("%d : %s %s\n", len(r.FinalResult)-i, r.FinalResult[i].Firstname, r.FinalResult[i].Lastname)
	}
	//time.Sleep(1 * time.Second)
	return nil
}
