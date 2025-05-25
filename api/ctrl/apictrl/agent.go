package apictrl

import (
	"api/config"
	"api/repos/agentsRepo"
	"encoding/json"
	"log"
	"net/http"
)

func GetAgents(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != *config.ApiSecretKey {
		log.Println("api: unauthorized: GetAgents")
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	agents, err := agentsRepo.GetMultiple()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(agents) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if err := json.NewEncoder(w).Encode(&agents); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
