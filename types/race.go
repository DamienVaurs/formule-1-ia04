package types

import "time"

type Race struct {
	Id             string      // Race ID
	Circuit        *Circuit    // Circuit
	Date           time.Time   // Date
	Team           []*Team     // Set of teams
	MeteoCondition Meteo       // Meteo condition
	FinalResult    []*Driver   // Final result, drivers rank from 1st to last
	HighLigths     []Highlight // Containes all what happend during the race
}

func NewRace(id string, circuit *Circuit, date time.Time, teams []*Team, meteo Meteo) *Race {

	d := make([]*Team, len(teams))
	copy(d, teams)

	f := make([]*Driver, len(teams)*2) //car 2 drievrs par team

	h := make([]Highlight, 0)

	return &Race{
		Id:             id,
		Circuit:        circuit,
		Date:           date,
		Team:           d,
		MeteoCondition: meteo,
		FinalResult:    f,
		HighLigths:     h,
	}
}

/*
func (r *Race) SimulateRace() {
	//On lance les agents pilotes
	for _, team := range r.Team {
		for _, driver := range team.Drivers {
			go driver.SimulateDriver(r)
		}
	}
	var nbFinish = 0
	//On simule tant que tous les pilotes n'ont pas fini la course
	for nbFinish < len(r.Team)*2 {

	}
}
*/
