package types

import (
	"math/rand"
	"time"
)

type Personnality struct {
	TraitsValue map[string]int // Dictionnaire de traits
}

func NewPersonnality(traitsValue map[string]int) *Personnality {
	return &Personnality{
		TraitsValue: traitsValue,
	}
}

func GenerateTraits() map[string]int {
	rand.NewSource(time.Now().UnixNano())
	// Un trait est un entier entre 1 et 5
	random_agressivity := rand.Intn(5) + 1   // AGRESSIVITY 	(statique -> impacte les proba de tentatives de dépassement)
	random_confidence := rand.Intn(5) + 1    // CONFIDENCE  	(dynamique -> impacte un peu la proba de tenter et la proba de réussir un dépassement)
	random_docility := rand.Intn(5) + 1      // DOCILITY    	(dynamique -> impacte la proba d'écouter la stratégie de l'équipe)
	random_concentration := rand.Intn(5) + 1 // CONCENTRATION 	(statique -> impacte la proba de réussir un dépassement)

	traits := make(map[string]int, 4)
	aggressivity := random_agressivity
	confidence := random_confidence
	docility := random_docility
	concentration := random_concentration
	traits["Aggressivity"] = aggressivity
	traits["Confidence"] = confidence
	traits["Docility"] = docility
	traits["Concentration"] = concentration
	return traits
}
