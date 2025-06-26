package apictrl

import (
	"api/models"
	"api/repos"
	"api/repos/channelMessagesRepo"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

func GetChannelMessages(w http.ResponseWriter, r *http.Request) {
	currentUsername, permissions, reject := authenticate(w, r)
	if reject {
		return
	}
	if authorize(models.PERMISSION_CHAT_LIST_CHANNEL_MESSAGES, &currentUsername, permissions) {
		return
	}

	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	limit, err := strconv.Atoi(q.Get("limit"))
	if err != nil {
		limit = repos.DEFAULT_LIMIT
	}
	offset := page * limit

	channelName := r.PathValue("channelName")

	messages, err := channelMessagesRepo.GetMultiple(channelName, offset, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(messages) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if err := json.NewEncoder(w).Encode(&messages); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func InsertChannelMessage(w http.ResponseWriter, r *http.Request) {
	currentUsername, permissions, reject := authenticate(w, r)
	if reject {
		return
	}
	if authorize(models.PERMISSION_CHAT_INSERT_CHANNEL_MESSAGE, &currentUsername, permissions) {
		return
	}

	channelName := r.PathValue("channelName")

	messageContent, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(messageContent) == 0 || len(messageContent) > 40_960 {
		http.Error(w, "invalid message", http.StatusBadRequest)
		return
	}

	if err := channelMessagesRepo.Insert(&models.ChannelMessage{
		Channel: channelName,
		Data:    string(messageContent),
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func DeleteChannelMessages(w http.ResponseWriter, r *http.Request) {
	// TODO: Delete all messages of a specific channel.
}
