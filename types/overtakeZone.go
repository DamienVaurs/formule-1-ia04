package types

type OverTakeZone struct {
	IntentOverTakeProbability float64 // Probabilité de tenter un dépassement
	OverTakeProbability       float64 // Probabilité de réussir le dépassement (vTODO : vraiment utile ?)
}
