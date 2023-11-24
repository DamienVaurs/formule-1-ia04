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
	fmt.Println("Les pilotes sont sur la ligne de départ")
	//On lance les agents pilotes
	fmt.Println("Début de la course...")
	for _, driver := range drivers {
		go driver.Start(driver.Position, 1) //todo : mettre le bon paramètre r.Circuit.NbLaps à la place de 1
	}
	var nbFinish = 0
	var nbDrivers = len(r.Teams) * 2
	decisionMap := make(map[string]Action, nbDrivers)

	//On simule tant que tous les pilotes n'ont pas fini la course
	for nbFinish < nbDrivers {
		fmt.Printf("============ NOUVELLE BOUCLE : %d ont fini =================\n\n", nbFinish)
		time.Sleep(1 * time.Second)
		//Chaque pilote, dans un ordre aléatoire, réalise les tests sur la proba de dépasser etc...
		drivers = ShuffleDrivers(drivers)
		//fmt.Println("LLLA")
		//fmt.Println(drivers[0].Position.Id)

		//fmt.Println("Débloquage des go routines...")

		for i := range drivers {
			//On débloque les pilotes en course pour qu'ils prennent une décision
			if drivers[i].Status == CRASHED || drivers[i].Status == ARRIVED {
				continue
			}
			//fmt.Println("Envoie de déblocage à : " + driver.Driver.Lastname)
			drivers[i].ChanEnv <- 1
		}
		//fmt.Println("Les go routines sont débloquées")
		// On récupère les décisions des pilotes
		for i := range drivers {
			if drivers[i].Status == CRASHED || drivers[i].Status == ARRIVED {
				continue
			}
			decisionMap[drivers[i].Driver.Id] = <-drivers[i].ChanEnv
		}
		//fmt.Println("On a toutes les décisions")

		//On traite les décisions et on met à jour les positions des pilotes
		for i := range drivers {
			decision := decisionMap[drivers[i].Driver.Id]
			switch decision {
			case TRY_OVERTAKE:
				//On vérifie si le pilote peut bien dépasser
				driverPortion := drivers[i].Position
				fmt.Printf("Portion pilote %s : %s\n", drivers[i].Driver.Lastname, driverPortion.Id) //TODO : n'est pas ok, fixé à straight_1
				driverToOvertake, err := drivers[i].DriverToOvertake()                               //TODO: pb ici, cherche toujours sur straight_1
				if err != nil {
					log.Printf("Error while getting driver to overtake: %s\n", err)
				}
				if driverToOvertake != nil {
					//On vérifie si le pilote a réussi son dépassement
					success, crashedDrivers := drivers[i].Overtake(driverToOvertake)
					if crashedDrivers != nil {
						//On supprime les pilotes crashés
						for _, crashedDriver := range crashedDrivers {
							crashedDriver.Status = CRASHED
							fmt.Println("Le pilote " + crashedDriver.Driver.Lastname + " a crashé")
							r.FinalResult = append(r.FinalResult, crashedDriver.Driver) //on l'ajoute au tableau
							drivers[i].Position.RemoveDriverOn(crashedDriver)
							/*fmt.Print("Après remove : ")
							driver.Position.DisplayDriversOn()*/
							nbFinish++
						}

						if success {
							//On met à jour les positions
							fmt.Println("Le pilote " + drivers[i].Driver.Lastname + " a réussi son dépassement")
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
							fmt.Printf("\nPilote %s a fini !!! \n", driver.Driver.Lastname)
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
		fmt.Println("Portion des pilote après maj : ") //est ok
		for i := range drivers {
			fmt.Println(drivers[i].Position.Id, drivers[i].Position.DriversOn)
		}

		//On met à jour les positions des pilotes
		for i := range r.Circuit.Portions {
			//fmt.Printf("On remplace %s par %s\n", portion.DriversOn, newDriversOnPortion[i]) semble bon
			r.Circuit.Portions[i].DriversOn = make([]*DriverInRace, len(newDriversOnPortion[i])) //on écrase l'ancien slice
			copy(r.Circuit.Portions[i].DriversOn, newDriversOnPortion[i])                        //on remplace par le nouveau
			//fmt.Printf("%s après update : %s \n", portion.Id, portion.DriversOn) semble ok
			//fmt.Println("UUUUU", r.Circuit.Portions[i].Id, r.Circuit.Portions[i].DriversOn)
		}
		fmt.Println("Portion des pilote après maj de ler position + maj des portions : ") //est ok
		for i := range drivers {
			fmt.Println(drivers[i].Position.Id, drivers[i].Position.DriversOn)
		}
		fmt.Println("Classement intermédiraire : ", r.FinalResult)

	}
	//On affiche le classement
	fmt.Println("Classement final :")
	for i := range r.FinalResult {
		fmt.Printf("%d : %s %s\n", len(r.FinalResult)-i, r.FinalResult[i].Firstname, r.FinalResult[i].Lastname)
	}
	time.Sleep(5 * time.Second)
	return nil
}
