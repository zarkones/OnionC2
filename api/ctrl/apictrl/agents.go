package apictrl

import (
	"api/models"
	"api/repos/agentsRepo"
	"encoding/json"
	"net/http"
	"slices"
	"strings"
)

// GetAgents returns list of registered agents.
func GetAgents(w http.ResponseWriter, r *http.Request) {
	_, _, reject := authenticateAndAuthorize(w, r, models.PERMISSION_AGENTS_LIST, nil)
	if reject {
		return
	}

	q := r.URL.Query()
	originsRaw := q.Get("origins")
	origins := strings.Split(originsRaw, ",")
	if len(originsRaw) == 0 {
		origins = []string{}
	}

	var agents []models.Agent
	var err error

	if len(origins) == 0 {
		agents, err = agentsRepo.GetMultiple()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		if slices.Contains(origins, "unknown") {
			unknownOriginAgents, err := agentsRepo.GetMultipleUnknownOrigins()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			agents = append(agents, unknownOriginAgents...)
		}

		agentsByCountry, err := agentsRepo.GetMultipleByCountryCode(origins)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		agents = append(agents, agentsByCountry...)
	}

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
