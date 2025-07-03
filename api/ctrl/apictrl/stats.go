package apictrl

import (
	"api/models"
	"api/repos/agentsRepo"
	"encoding/json"
	"log"
	"net/http"
)

func GetStatsAgents(w http.ResponseWriter, r *http.Request) {
	_, _, reject := authenticateAndAuthorize(w, r, models.PERMISSION_AGENTS_STATS, nil)
	if reject {
		return
	}

	count, err := agentsRepo.GetMultipleUnknownOriginsCount()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]int64{
		"unknownOriginCount": count,
	}

	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		log.Println("api: error: serializing response:", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
