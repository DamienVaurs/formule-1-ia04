package restserver

import (
	"log"
	"net/http"
	"sync"
	"time"

	"gitlab.utc.fr/vaursdam/formule-1-ia04/types"
)

type RestServer struct {
	sync.Mutex
	addr              string
	pointTabCircuit   []*types.Circuit //circuits
	pointTabTeam      []*types.Team    //current teams
	initPersonalities map[string]types.Personality
}

func NewRestServer(addr string, pointTabCircuit []*types.Circuit, pointTabTeam []*types.Team, personalities map[string]types.Personality) *RestServer {
	return &RestServer{addr: addr, pointTabCircuit: pointTabCircuit, pointTabTeam: pointTabTeam, initPersonalities: personalities}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (rsa *RestServer) Start() {
	//initialise statistics
	statistics = &types.SimulateChampionship{}
	for _, team := range rsa.pointTabTeam {
		for _, driver := range team.Drivers {
			statistics.TotalStatistics.DriversTotalPoints = append(statistics.TotalStatistics.DriversTotalPoints, &types.DriverTotalPoints{Driver: driver.Lastname, TotalPoints: 0})
			statistics.LastChampionshipStatistics.DriversTotalPoints = append(statistics.LastChampionshipStatistics.DriversTotalPoints, &types.DriverTotalPoints{Driver: driver.Lastname, TotalPoints: 0})
		}
		statistics.TotalStatistics.TeamsTotalPoints = append(statistics.TotalStatistics.TeamsTotalPoints, &types.TeamTotalPoints{Team: team.Name, TotalPoints: 0})
		statistics.LastChampionshipStatistics.TeamsTotalPoints = append(statistics.LastChampionshipStatistics.TeamsTotalPoints, &types.TeamTotalPoints{Team: team.Name, TotalPoints: 0})
	}

	// création du multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/simulateChampionship", rsa.startSimulation)
	mux.HandleFunc("/personalities", rsa.getAndUpdatePersonalities)
	mux.HandleFunc("/statisticsChampionship", rsa.statisticsChampionship)
	mux.HandleFunc("/reset", rsa.reset)
	corsHandler := corsMiddleware(mux)

	// création du serveur http
	s := &http.Server{
		Addr:           rsa.addr,
		Handler:        corsHandler,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20}

	// lancement du serveur
	log.Println("Listening on", rsa.addr)
	go log.Fatal(s.ListenAndServe())

}
