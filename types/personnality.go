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
	Type TraitType // Trait name
	//Description string // Trait description
}

type TraitType int

const (
	AGRESSIVITY   TraitType = iota //statique -> impacte les proba de tentatives de dépassement
	CONFIDENCE                     // dynamique -> impatce un peu la proba de tenter et la proba de réussir un dépassement
	DOCILITY                       // dynamique -> impacte la proba d'écouter la stratégie de l'équipe
	CONCENTRATION                  //statique -> impacte la proba de réussir un dépassement

	//TODO compléter avec les autres traits
)
