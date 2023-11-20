package types

import (
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
	CIn  chan int   //Canal pour recevoir les messages des pilotes
	COut []chan int //Canal pour envoyer les messages aux pilotes (1 par pilote)
}

func NewRace(id string, circuit *Circuit, date time.Time, teams []*Team, meteo Meteo) *Race {

	d := make([]*Team, len(teams))
	copy(d, teams)

	f := make([]*Driver, len(teams)*2) //car 2 drivers par team

	h := make([]Highlight, 0)

	cIn := make(chan int)
	cOut := make([]chan int, len(teams)*2)

	return &Race{
		Id:             id,
		Circuit:        circuit,
		Date:           date,
		Teams:          d,
		MeteoCondition: meteo,
		FinalResult:    f,
		HighLigths:     h,
		CIn:            cIn,
		COut:           cOut,
	}
}

func (r *Race) SimulateRace() {
	//On crée les instances des pilotes en course

	var drivers = SliceOfDrivers(r.Teams, &(r.Circuit.Portions[0]))
	//On lance les agents pilotes
	for i, driver := range drivers {
		go driver.Start(r.CIn, r.COut[i], driver.Position, driver.NbLaps)
	}

	var nbFinish = 0
	var nbDrivers = len(r.Teams) * 2

	//On simule tant que tous les pilotes n'ont pas fini la course
	for nbFinish < nbDrivers {
		//Chaque pilote, dans un ordre aléatoire, réalise les tests sur la proba de dépasser etc...
		drivers = ShuffleDrivers(drivers)
		for _, driver := range drivers {
			//On dit au pilote qu'il peut jouer
			driver.ChanEnvIn <- 1
		}
		//On attend que tous les pilotes aient fini de jouer
		for i := 0; i < nbDrivers; i++ {
			ret := <-r.CIn
			nbFinish += ret // on dit que le pilote envoie 1 s'il a fini (ou crash), 0 sinon
		}
	}
}
