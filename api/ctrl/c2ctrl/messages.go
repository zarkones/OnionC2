package c2ctrl

import (
	"api/repos/messagesRepo"
	"encoding/json"
	"log"
	"net/http"
)

func GetMessages(w http.ResponseWriter, r *http.Request) {
	agentID := r.PathValue("agentID")

	messages, err := messagesRepo.GetMultipleForAgent(agentID)
	if err != nil {
		log.Println("c2: error: GetMessages:", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if len(messages) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if err := json.NewEncoder(w).Encode(&messages); err != nil {
		log.Println("c2: error: serializing response:", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

type AgentMsgRespCtx struct {
	MessageID string `json:"messageId"`
	Response  string `json:"response"`
}

func InsertMessageResponse(w http.ResponseWriter, r *http.Request) {
	var newMsg AgentMsgRespCtx

	if err := json.NewDecoder(r.Body).Decode(&newMsg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if newMsg.Response == "" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if newMsg.MessageID == "" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if err := messagesRepo.UpdateResponse(newMsg.MessageID, newMsg.Response); err != nil {
		log.Println("c2: error: messagesRepo.UpdateResponse:", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
