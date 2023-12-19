package restserver

import (
	"encoding/json"
	"net/http"
)

// Obtenir le classements des pilotes à la fin d'un championnat
func (rsa *RestServer) getChampionshipRank(w http.ResponseWriter, r *http.Request) {
	// vérification de la méthode de la requête
	if r.Method != "GET" {
		return
	}

	serial, _ := json.Marshal(driversRank)
	w.Write(serial)
}
