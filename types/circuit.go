package types

type Circuit struct {
	CircuitId   string    // Circuit ID
	CircuitName string    // Circuit name
	TurnList    []Turn    // Turn list
	DRSZoneList []DRSZone // DRS zone list
	Country     string    // Country

}
