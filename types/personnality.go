package types

type Personnality struct {
	TraitsValue map[*Trait]int // Dictionnary of traits
}

func NewPersonnality(traitsValue map[*Trait]int) *Personnality {
	m := make(map[*Trait]int)
	return &Personnality{
		TraitsValue: m,
	}
}

// Les traits sont à définir au lancement du programme. La personnalité utilisera les différents traits
type Trait struct {
	Id   string    // Trait ID
	Name TraitType // Trait name
	//Description string // Trait description
}

type TraitType int

const (
	AGRESSIVITY TraitType = iota
	CONFIDENCE
	//TODO compléter avec les autres traits
)
