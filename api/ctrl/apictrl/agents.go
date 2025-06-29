package apictrl

import (
	"api/models"
	"api/repos/agentsRepo"
	"encoding/json"
	"net/http"
)

// GetAgents returns list of registered agents.
func GetAgents(w http.ResponseWriter, r *http.Request) {
	_, _, reject := authenticateAndAuthorize(w, r, models.PERMISSION_AGENTS_LIST, nil)
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
