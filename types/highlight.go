package types

type Highlight struct {
	Id          string        // Highlight ID
	Description string        // Describe the highlight
	Drivers     []*Driver     // Drivers involved in the highlight
	Type        HighlightType // Type of highlight
}

type HighlightType int

const (
	DRIVER_CRASH HighlightType = iota
	DRIVER_OVERTAKE
	//DRIVER_PITSTOP
	//DRIVER_PENALTY
	//DRIVER_FASTEST_LAP
)
