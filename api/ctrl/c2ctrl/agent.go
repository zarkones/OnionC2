package c2ctrl

import (
	"api/models"
	"api/repos/agentsRepo"
	"encoding/json"
	"log"
	"net/http"
)

func InsertAgent(w http.ResponseWriter, r *http.Request) {
	var newAgent models.Agent

	if err := json.NewDecoder(r.Body).Decode(&newAgent); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(newAgent.Hostname) == 0 {
		http.Error(w, "hostname is empty", http.StatusUnprocessableEntity)
		return
	}

	if err := agentsRepo.Insert(&newAgent); err != nil {
		log.Println("error: InsertAgent: agentsRepo.Insert:", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	resp := map[string]string{
		"id": newAgent.ID,
	}

	serializedResp, err := json.Marshal(&resp)
	if err != nil {
		log.Println("error: InsertAgent: serializing the response:", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(serializedResp)
}
