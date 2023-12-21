package restserver

import (
	"encoding/json"
	"fmt"
	"net/http"
)

/**
*
* Return statistics about simulations WITHOUT launching a new simulation
*
 */
func (rsa *RestServer) statisticsChampionship(w http.ResponseWriter, r *http.Request) {
	// vérification de la méthode de la requête
	if r.Method != "GET" {
		return
	}
	fmt.Println("GET /statisticsChampionship")
	w.WriteHeader(http.StatusOK)
	serial, _ := json.Marshal(statistics) //statistics is defined in simulateChampionship.go
	w.Write(serial)
}
