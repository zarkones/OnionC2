package apictrl

import (
	"api/repos/agentsRepo"
	"encoding/json"
	"net/http"
)

// GetAgents returns list of registered agents.
func GetAgents(w http.ResponseWriter, r *http.Request) {
	_, _, reject := authenticate(w, r)
	if reject {
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
