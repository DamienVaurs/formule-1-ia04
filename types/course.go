package types

import "time"

type Course struct {
	CourseId string    // Course ID
	Circuit  Circuit   // Circuit
	Date     time.Time // Date
	Drivers  []Driver  // Drivers
	Meteo    Meteo     // Meteo (TODO : définir type Meteo)

}
