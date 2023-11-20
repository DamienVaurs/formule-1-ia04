package types

type PortionType int

const (
	TURN PortionType = iota
	STRAIGHT
)

type Portion struct {
	Id         string      // Portion ID
	Difficulty int         // Difficulty of the portion in [0,5]
	DriversOn  []*Driver   // Drivers on the portion
	Type       PortionType // Type of the portion
	IsDRSZone  bool        // True if is a DRS Zone. -> increases chances of overtaking
}

func NewPortion(id string, difficulty int, driversOn []*Driver, isDRSZone bool) *Portion {

	d := make([]*Driver, len(driversOn))
	copy(d, driversOn)

	var t PortionType
	if len(id) < len("turn") {
	} else if id[:len("turn")] == "turn" {
		t = TURN
	} else {
		t = STRAIGHT
	}
	return &Portion{
		Id:         id,
		Difficulty: difficulty,
		DriversOn:  d,
		Type:       t,
		IsDRSZone:  isDRSZone,
	}
}
