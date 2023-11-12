package types

type Personnality struct {
	TraitsValue map[*Trait]int // Dictionnary of traits
}

func NewPersonnality(traitsValue map[*Trait]int) *Personnality {
	return &Personnality{
		TraitsValue: traitsValue,
	}
}

// Les traits sont à définir au lancement du programme. La personnalité utilisera les différents traits
type Trait struct {
	Id   string // Trait ID
	Name string // Trait name
	//Description string // Trait description
}
