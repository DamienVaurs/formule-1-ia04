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
	Id                     string    // Portion ID
	CrashProbability       int       // Crash probability in [1,10]
	TryOvertakeProbability float64   //Probability of intenting an overtake in [0,1]
	DriversOn              []*Driver // Drivers on the portion
}

func NewPortion(id string, crashProbability int, tryOvertakeProbability float64, driversOn []*Driver) *Portion {

	d := make([]*Driver, len(driversOn))
	copy(d, driversOn)

	return &Portion{
		Id:                     id,
		CrashProbability:       crashProbability,
		TryOvertakeProbability: tryOvertakeProbability,
		DriversOn:              d,
	}
}

/****** TURN ******/
type Turn struct {
	Portion
}

func NewTurn(id string, crashProbability int, tryOvertakeProbability float64, driversOn []*Driver) *Turn {

	d := make([]*Driver, len(driversOn))
	copy(d, driversOn)

	return &Turn{
		Portion: Portion{
			Id:                     id,
			CrashProbability:       crashProbability,
			TryOvertakeProbability: tryOvertakeProbability,
			DriversOn:              d,
		},
	}
}

func (t *Turn) GetType() PortionType {
	return TURN
}

/****** STRAIGHT ******/
type Straight struct {
	Portion
	IsDRSZone bool // True if is a DRS Zone. -> increases chances of overtaking
}

func NewStraight(id string, crashProbability int, tryOvertakeProbability float64, driversOn []*Driver, isDRSZone bool) *Straight {

	d := make([]*Driver, len(driversOn))
	copy(d, driversOn)

	return &Straight{
		Portion: Portion{
			Id:                     id,
			CrashProbability:       crashProbability,
			TryOvertakeProbability: tryOvertakeProbability,
			DriversOn:              d,
		},
		IsDRSZone: isDRSZone,
	}
}

func (s *Straight) GetType() PortionType {
	return STRAIGHT
}
