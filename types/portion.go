package types

type PortionType int

const (
	TURN PortionType = iota
	STRAIGHT
)

// L'interface est utilisée pour pouvoir faire un slice de portions mélageant lignes droites et virages
type PortionInt interface {
	GetType() PortionType
}

type Portion struct {
	Id        string    // Portion ID
	Diffculty int       // Difficulty of the portion in [0,5]
	DriversOn []*Driver // Drivers on the portion
}

/****** TURN ******/
type Turn struct {
	Portion
}

func (t *Turn) GetType() PortionType {
	return TURN
}

func NewTurn(id string, difficulty int, driversOn []*Driver) *Turn {

	d := make([]*Driver, len(driversOn))
	copy(d, driversOn)

	return &Turn{
		Portion: Portion{
			Id:        id,
			Diffculty: difficulty,
			DriversOn: d,
		},
	}
}

func NewStraight(id string, difficulty int, driversOn []*Driver, isDRSZone bool) *Straight {

	d := make([]*Driver, len(driversOn))
	copy(d, driversOn)

	return &Straight{
		Portion: Portion{
			Id:        id,
			Diffculty: difficulty,
			DriversOn: d,
		},
		IsDRSZone: isDRSZone,
	}
}

/****** STRAIGHT ******/
type Straight struct {
	Portion
	IsDRSZone bool // True if is a DRS Zone. -> increases chances of overtaking
}

func (s *Straight) GetType() PortionType {
	return STRAIGHT
}
