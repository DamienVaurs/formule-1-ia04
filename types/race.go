package types

import "time"

type Race struct {
	CourseId       string      // Course ID
	Circuit        *Circuit    // Circuit
	Date           time.Time   // Date
	Drivers        []*Driver   // Drivers
	MeteoCondition Meteo       // Meteo condition
	FinalResult    []*Driver   // Final result, drivers rank from 1st to last
	HighLigths     []Highlight // Containes all what happend during the race
}
