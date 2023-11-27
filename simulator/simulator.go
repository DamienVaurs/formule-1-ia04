package simulator

import (
	"fmt"
	"log"
	"time"

	"gitlab.utc.fr/vaursdam/formule-1-ia04/types"
)

type Simulator struct {
	Championships []types.Championship
}

func NewSimulator(championships []types.Championship) *Simulator {
	c := make([]types.Championship, len(championships))
	copy(c, championships)

	return &Simulator{
		Championships: c,
	}
}

func (s *Simulator) LaunchSimulation() {
	log.Println("Lancement d'une nouvelle simulation...")
	for _, championship := range s.Championships {
		//On simule chaque championnat
		log.Printf("Lancement d'un nouveau championnat : %s...\n", championship.Name)
		for i, circuit := range championship.Circuits {
			//On simule chaque course
			//Etape 1 : on crée la course
			var id = circuit.Name + " " + championship.Name

			var date = time.Now()
			if i != 0 {
				date = championship.Races[i-1].Date.AddDate(0, 0, 14)
			}
			var meteo = circuit.GenerateMeteo()
			new_Race := types.NewRace(id, circuit, date, championship.Teams, meteo)

			//Etape 2 (la principale) : on joue la course
			pointsMap, err := new_Race.SimulateRace()
			if err != nil {
				log.Printf("Erreur simulation cours %s : %s\n", new_Race.Id, err.Error())
			}

			//On enregistre les points gagnés par chaque pilote
			for indT := range championship.Teams {
				for indD := range championship.Teams[indT].Drivers {
					championship.Teams[indT].Drivers[indD].ChampionshipPoints += pointsMap[championship.Teams[indT].Drivers[indD].Id]
				}
			}
			/*
				for _, team := range championship.Teams {
					log.Printf("%s : %d points\n", team.Name, team.CalcChampionshipPoints())
					for _, driver := range team.Drivers {
						log.Printf("	%s %s : %d points\n", driver.Firstname, driver.Lastname, driver.ChampionshipPoints)
					}
				}*/

			//Etape 3 : on ajoute la course au championnat
			fmt.Println("Ajout de la course au championnat...")
			championship.Races[i] = *new_Race
		}
		//On affiche le classement du championnat
		log.Printf("\n\n===== Classements du championnat %s =====\n", championship.Name)
		championship.DisplayTeamRank()
		championship.DisplayDriverRank()

	}
}
