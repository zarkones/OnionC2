package apictrl

import (
	"api/models"
	"api/repos"
	"api/repos/filesRepo"
	"api/repos/messagesRepo"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type GetMessagesRespCtx struct {
	Messages []models.Message `json:"messages"`
	Since    string           `json:"since"`
}

// GetMessages returns messages exchanged with a specific agent.
func GetMessages(w http.ResponseWriter, r *http.Request) {
	_, _, reject := authenticateAndAuthorize(w, r, models.PERMISSION_AGENTS_LIST_MESSAGES, nil)
	if reject {
		return
	}

	agentID := r.PathValue("agentID")

	q := r.URL.Query()
	since, errSince := strconv.ParseInt(q.Get("since"), 10, 64)
	page, _ := strconv.Atoi(q.Get("page"))
	limit, err := strconv.Atoi(q.Get("limit"))
	if err != nil {
		limit = repos.DEFAULT_LIMIT
	}
	offset := page * limit

	var messages []models.Message

	if errSince == nil {
		messages, err = messagesRepo.GetMultipleSince(agentID, int64(since), limit)
	} else {
		messages, err = messagesRepo.GetMultiple(agentID, offset, limit)
	}

	if err != nil {
		log.Println("api: error: GetMessages:", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if len(messages) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	resp := GetMessagesRespCtx{
		Messages: messages,
		Since:    fmt.Sprint(messages[len(messages)-1].CreatedAt),
	}

	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		log.Println("api: error: serializing response:", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

// InsertMessage sends a message to a specific agent.
func InsertMessage(w http.ResponseWriter, r *http.Request) {
	_, _, reject := authenticateAndAuthorize(w, r, models.PERMISSION_AGENTS_INSERT_MESSAGE, nil)
	if reject {
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

	if strings.HasPrefix(newMsg.Request, "/upload-file|") {
		filePath := strings.TrimPrefix(newMsg.Request, "/upload-file|")
		file := &models.File{
			AgentID:      newMsg.AgentID,
			OriginalPath: filePath,
			Order:        models.ORDER_UPLOAD,
		}
		if err := filesRepo.Insert(file); err != nil {
			fmt.Println("api: error: filesRepo.Insert:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Uses the ID set by BeforeCreate
		newMsg.Request = "/upload-file|" + filePath + "|" + file.ID
	}

	if strings.HasPrefix(newMsg.Request, "/download-file|") {
		filePath := strings.TrimPrefix(newMsg.Request, "/download-file|")
		file := &models.File{
			AgentID:      newMsg.AgentID,
			OriginalPath: filePath,
			Order:        models.ORDER_DOWNLOAD,
		}
		if err := filesRepo.Insert(file); err != nil {
			fmt.Println("api: error: filesRepo.Insert:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Uses the ID set by BeforeCreate
		newMsg.Request = "/download-file|" + filePath + "|" + file.ID
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
