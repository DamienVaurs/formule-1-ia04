package simulator

import (
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
	for _, championship := range s.Championships {
		//On simule chaque championnat
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
			new_Race.SimulateRace()

			//Etape 3 : on ajoute la course au championnat
			championship.Races[i] = new_Race
		}
	}
}
