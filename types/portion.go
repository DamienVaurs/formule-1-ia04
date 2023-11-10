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
	Id                     string  // Portion ID
	CrashProbability       int     // Crash probability in [1,10]
	TryOvertakeProbability float64 //Probability of intenting an overtake in [0,1]
}

/****** TURN ******/
type Turn struct {
	Portion
}

func (t *Turn) GetType() PortionType {
	return TURN
}

/****** STRAIGHT ******/
type Straight struct {
	Portion
	IsDRSZone bool // True if is a DRS Zone. -> increases chances of overtaking
}

func (s *Straight) GetType() PortionType {
	return STRAIGHT
}
