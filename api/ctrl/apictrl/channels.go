package apictrl

import (
	"api/models"
	"api/repos/channelSecretsRepo"
	"api/repos/channelsRepo"
	"api/repos/permissionsRepo"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

func GetChannels(w http.ResponseWriter, r *http.Request) {
	_, _, reject := authenticateAndAuthorize(w, r, models.PERMISSION_CHAT_LIST_CHANNELS, nil)
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
	type Ctx struct {
		models.Channel
		InvitedOperators []struct {
			Username                     string `json:"username"`
			HexEncodedRsaEncryptedAesKey string `json:"secret"`
		} `json:"invitedUsernames"`
	}

	currentUsername, _, reject := authenticateAndAuthorize(w, r, models.PERMISSION_CHAT_INSERT_CHANNEL, nil)
	if reject {
		return
	}

	var ctx Ctx

	if err := json.NewDecoder(r.Body).Decode(&ctx); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if ctx.InvitedOperators == nil {
		http.Error(w, "invited operators is nil", http.StatusBadRequest)
		return
	}

	channel := models.Channel{
		Name:        ctx.Name,
		Description: ctx.Description,
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

	issues := []error{}
	for _, invitedOperator := range ctx.InvitedOperators {
		if err := permissionsRepo.Insert(&models.Permission{
			Key:      models.PERMISSION_CHAT_LIST_CHANNEL_MESSAGES,
			Username: invitedOperator.Username,
			Metadata: channel.Name,
		}); err != nil {
			issues = append(issues, err)
			log.Println("failed to give access to channel to user:", invitedOperator.Username, channel.Name)
			continue
		}

		if err := channelSecretsRepo.Insert(&models.ChannelSecret{
			RecipientOperatorUsername:    invitedOperator.Username,
			HexEncodedRsaEncryptedAesKey: invitedOperator.HexEncodedRsaEncryptedAesKey,
			CreatedBy:                    currentUsername,
		}); err != nil {
			issues = append(issues, err)
			log.Println("failed to insert channel secret for user:", invitedOperator.Username, channel.Name)
			continue
		}
	}
	if len(issues) != 0 {
		http.Error(w, errors.Join(issues...).Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func UpdateChannel(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func DeleteChannels(w http.ResponseWriter, r *http.Request) {
	_, _, reject := authenticateAndAuthorize(w, r, models.PERMISSION_CHAT_DELETE_CHANNEL, nil)
	if reject {
		return
	}

	channelName := r.PathValue("channelName")

	if err := channelsRepo.Delete(channelName); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: Trigger message deletion somehow.
	// TODO: Trigger channel secrets somehow.
}
