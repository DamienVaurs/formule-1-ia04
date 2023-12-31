package restserver

import (
	"fmt"
	"net/http"
)

func (rsa *RestServer) startRaceSimulation(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		return
	}
	fmt.Println("GET /simulateRace")

	race := NewRace()

}
