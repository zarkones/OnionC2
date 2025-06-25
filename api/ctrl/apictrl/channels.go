package apictrl

import (
	"api/core"
	"api/models"
	"api/repos/channelSecretsRepo"
	"api/repos/channelsRepo"
	"api/repos/permissionsRepo"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"slices"
	"sort"
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
			Channel:                      channel.Name,
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

func InviteToChannel(w http.ResponseWriter, r *http.Request) {
	currentUsername, permissions, reject := authenticate(w, r)
	if reject {
		return
	}
	if authorize(models.PERMISSION_CHAT_LIST_CHANNEL_MESSAGES, &currentUsername, permissions) {
		return
	}

	channelName := r.PathValue("channelName")
	invitedOperatorUsername := r.PathValue("invitedOperatorUsername")

	type Ctx struct {
		HexEncodedRsaEncryptedAesKey string `json:"secret"`
	}

	var ctx Ctx

	if err := json.NewDecoder(r.Body).Decode(&ctx); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(ctx.HexEncodedRsaEncryptedAesKey) == 0 {
		http.Error(w, "invalid secret", http.StatusBadRequest)
		return
	}

	if err := channelSecretsRepo.Insert(&models.ChannelSecret{
		RecipientOperatorUsername:    invitedOperatorUsername,
		Channel:                      channelName,
		HexEncodedRsaEncryptedAesKey: ctx.HexEncodedRsaEncryptedAesKey,
		CreatedBy:                    currentUsername,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func RemoveFromChannel(w http.ResponseWriter, r *http.Request) {
	currentUsername, currentUserPermissions, reject := authenticate(w, r)
	if reject {
		return
	}
	if authorize(models.PERMISSION_CHAT_LIST_CHANNEL_MESSAGES, &currentUsername, currentUserPermissions) {
		return
	}

	channelName := r.PathValue("channelName")

	var ctx struct {
		RemovedOperatorUsername string `json:"removedOperatorUsername"`
		Secrets                 []struct {
			Username                     string `json:"username"`
			HexEncodedRsaEncryptedAesKey string `json:"secret"`
		} `json:"secrets"`
	}

	if err := json.NewDecoder(r.Body).Decode(&ctx); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	channelPermissions, err := permissionsRepo.GetMultipleByChannel(channelName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(channelPermissions) == 0 {
		http.Error(w, "strangely there are no operators with access to this channel", http.StatusUnprocessableEntity)
		return
	}

	operatorUsernames := []string{}
	for _, permission := range channelPermissions {
		if permission.Username == ctx.RemovedOperatorUsername {
			continue
		}
		operatorUsernames = append(operatorUsernames, permission.Username)
	}
	operatorUsernames = core.Deduplicate(operatorUsernames)

	operatorUsernamesFromCtx := make([]string, len(ctx.Secrets))
	for i, secret := range ctx.Secrets {
		operatorUsernamesFromCtx[i] = secret.Username
	}
	operatorUsernamesFromCtx = core.Deduplicate(operatorUsernamesFromCtx)

	// Sorting so they can be compared
	operatorUsernames = sort.StringSlice(operatorUsernames)
	operatorUsernamesFromCtx = sort.StringSlice(operatorUsernamesFromCtx)

	if len(operatorUsernames) != len(operatorUsernamesFromCtx) || slices.Compare(operatorUsernames, operatorUsernamesFromCtx) != 0 {
		http.Error(w, "new secret must be provided for all operators with accesss to the channel", http.StatusUnprocessableEntity)
		return
	}

	issues := []error{}
	for _, secret := range ctx.Secrets {
		if err := channelSecretsRepo.Insert(&models.ChannelSecret{
			RecipientOperatorUsername:    secret.Username,
			Channel:                      channelName,
			HexEncodedRsaEncryptedAesKey: secret.HexEncodedRsaEncryptedAesKey,
			CreatedBy:                    currentUsername,
		}); err != nil {
			log.Println("failed to insert new channel secret during removal of an operator:", err, channelName)
			issues = append(issues, err)
			continue
		}
	}
	if len(issues) != 0 {
		http.Error(w, errors.Join(issues...).Error(), http.StatusInternalServerError)
		return
	}
}
