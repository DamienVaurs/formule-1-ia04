package types

type Highlight struct {
	Id          string        // Highlight ID
	Description string        // Describe the highlight
	Drivers     []*Driver     // Drivers involved in the highlight
	Type        HighlightType // Type of highlight
}

func NewHighlight(id string, description string, drivers []*Driver, highlightType HighlightType) *Highlight {

	d := make([]*Driver, len(drivers))
	copy(d, drivers)

	return &Highlight{
		Id:          id,
		Description: description,
		Drivers:     d,
		Type:        highlightType,
	}
}

type HighlightType int

const (
	CRASH HighlightType = iota
	OVERTAKE
	FINISH
	//DRIVER_PITSTOP
	//DRIVER_PENALTY
	//DRIVER_FASTEST_LAP
)
