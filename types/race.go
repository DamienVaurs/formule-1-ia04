package types

import "time"

type Race struct {
	Id             string      // Race ID
	Circuit        *Circuit    // Circuit
	Date           time.Time   // Date
	Drivers        []*Driver   // Drivers
	MeteoCondition Meteo       // Meteo condition
	FinalResult    []*Driver   // Final result, drivers rank from 1st to last
	HighLigths     []Highlight // Containes all what happend during the race
}

func NewRace(id string, circuit *Circuit, date time.Time, drivers []*Driver, meteo Meteo, finalResult []*Driver, highlights []Highlight) *Race {
	return &Race{
		Id:             id,
		Circuit:        circuit,
		Date:           date,
		Drivers:        drivers,
		MeteoCondition: meteo,
		FinalResult:    finalResult,
		HighLigths:     highlights,
	}
}
