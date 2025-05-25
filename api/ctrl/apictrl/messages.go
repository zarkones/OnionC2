package apictrl

import (
	"api/config"
	"api/models"
	"api/repos/messagesRepo"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func GetMessages(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != *config.ApiSecretKey {
		log.Println("api: unauthorized: GetMessages")
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	agentID := r.PathValue("agentID")

	messages, err := messagesRepo.GetMultiple(agentID)
	if err != nil {
		log.Println("api: error: GetMessages:", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if len(messages) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if err := json.NewEncoder(w).Encode(&messages); err != nil {
		log.Println("api: error: serializing response:", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func InsertMessage(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != *config.ApiSecretKey {
		log.Println("api: unauthorized: InsertMessage")
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	var newMsg models.Message

	if err := json.NewDecoder(r.Body).Decode(&newMsg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if newMsg.Response != "" {
		http.Error(w, `{"error":"cannot insert message with response"}`, http.StatusBadRequest)
		return
	}

	if err := messagesRepo.Insert(&newMsg); err != nil {
		fmt.Println("api: error: messagesRepo.Insert:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := map[string]string{
		"id": newMsg.ID,
	}

	serializedResp, err := json.Marshal(&resp)
	if err != nil {
		log.Println("api: error: InsertMessage: serializing the response:", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(serializedResp)
}
