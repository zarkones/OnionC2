package apictrl

import (
	"api/models"
	"api/repos/channelsRepo"
	"encoding/json"
	"net/http"
)

func GetChannels(w http.ResponseWriter, r *http.Request) {
	_, _, reject := authenticate(w, r, models.PERMISSION_CHAT_LIST_CHANNELS)
	if reject {
		return
	}

	channels, err := channelsRepo.GetMultiple()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(channels) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if err := json.NewEncoder(w).Encode(&channels); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func InsertChannel(w http.ResponseWriter, r *http.Request) {
	_, _, reject := authenticate(w, r, models.PERMISSION_CHAT_INSERT_CHANNEL)
	if reject {
		return
	}

	var channel models.Channel

	if err := json.NewDecoder(r.Body).Decode(&channel); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(channel.Name) == 0 || len(channel.Name) > 64 {
		http.Error(w, "invalid channel name", http.StatusBadRequest)
		return
	}
	if len(channel.Description) > 512 {
		http.Error(w, "invalid channel description", http.StatusBadRequest)
		return
	}

	if err := channelsRepo.Insert(&channel); err != nil {
		http.Error(w, "invalid channel description", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func UpdateChannel(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func DeleteChannels(w http.ResponseWriter, r *http.Request) {
	_, _, reject := authenticate(w, r, models.PERMISSION_CHAT_DELETE_CHANNEL)
	if reject {
		return
	}

	channelName := r.PathValue("channelName")

	if err := channelsRepo.Delete(channelName); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: Trigger message deletion somehow.
}
