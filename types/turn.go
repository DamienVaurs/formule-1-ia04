package types

type Turn struct {
	TurnId           string  // Turn ID
	CrashProbability float64 // Crash probability
	OverTakeZone
}
