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
	// Vérification de la méthode de la requête
	if r.Method != "GET" {
		return
	}
	fmt.Println("GET /statisticsChampionship")

	// Marshal des statistiques en JSON
	serial, err := json.Marshal(statistics)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// En-tête Content-Disposition pour indiquer le téléchargement du fichier
	w.Header().Set("Content-Disposition", "attachment; filename=statistics.json")
	w.Header().Set("Content-Type", "application/json")

	// Écriture des données dans la réponse
	w.Write(serial)
}
