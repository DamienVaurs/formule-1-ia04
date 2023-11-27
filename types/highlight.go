package types

import (
	"fmt"
)

type Highlight struct {
	Description string          // Describe the highlight
	Drivers     []*DriverInRace // Drivers involved in the highlight
	Type        HighlightType   // Type of highlight
}

type HighlightType int

const (
	CRASH HighlightType = iota
	OVERTAKE
	FINISH
	DRIVER_PITSTOP
	//DRIVER_PENALTY
	//DRIVER_FASTEST_LAP
)

func NewHighlight(drivers []*DriverInRace, highlightType HighlightType) (*Highlight, error) {

	d := make([]*DriverInRace, len(drivers))
	copy(d, drivers)
	var desc string

	switch highlightType {
	case CRASH:
		if len(drivers) > 1 {
			desc = fmt.Sprintf("CRASH au tour %d: Plusieurs pilotes sont rentrés en accident : %s et %s", drivers[0].NbLaps, drivers[0].Driver.Lastname, drivers[1].Driver.Lastname)
		} else if len(drivers) == 1 {
			desc = fmt.Sprintf("CRASH au tour %d: Le pilote %s a crashé", drivers[0].NbLaps, drivers[0].Driver.Lastname)
		} else {
			return nil, fmt.Errorf("CRASH highlight must include 1 or 2 drivers")
		}
	case OVERTAKE:
		if len(drivers) != 2 {
			return nil, fmt.Errorf("OVERTAKE highlight must have 2 drivers")
		}
		desc = fmt.Sprintf("DEPASSEMENT au tour %d: Le pilote %s a réussi son dépassement sur %s", drivers[0].NbLaps, drivers[0].Driver.Lastname, drivers[1].Driver.Lastname)
	case FINISH:
		if len(drivers) != 1 {
			return nil, fmt.Errorf("FINISH highlight must include exactly 1 driver")
		} else {
			desc = fmt.Sprintf("ARRIVEE: Le pilote %s est arrivé!", drivers[0].Driver.Lastname)
		}
	}
	return &Highlight{
		Description: desc,
		Drivers:     drivers,
		Type:        highlightType,
	}, nil
}

func (h *Highlight) String() string {
	return h.Description
}
