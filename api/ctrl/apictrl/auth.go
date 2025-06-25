package apictrl

import (
	"api/core/crypto"
	"api/models"
	"net/http"
)

func authenticate(w http.ResponseWriter, r *http.Request) (username string, permissions []models.Permission, reject bool) {
	token := r.Header.Get("Authorization")
	username, permissions, err := crypto.VerifyToken(token)
	if err != nil {
		http.Error(w, "", http.StatusUnauthorized)
		return "", nil, true
	}
	return username, permissions, false
}
