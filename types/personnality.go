package types

type Personality struct {
	TraitsValue map[string]int `json:"personality"` // Dictionnaire de traits
}

func NewPersonality(traitsValue map[string]int) *Personality {
	return &Personality{
		TraitsValue: traitsValue,
	}
}
