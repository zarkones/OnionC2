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
		// TODO: Log this incident.
		http.Error(w, "", http.StatusUnauthorized)
		return "", nil, true
	}

	return username, permissions, false
}

func authenticateAndAuthorize(w http.ResponseWriter, r *http.Request, permissionKey models.PermissionKey, metadata *string) (username string, permissions []models.Permission, reject bool) {
	username, permissions, reject = authenticate(w, r)
	if reject {
		return username, permissions, reject
	}

	if permissionKey != models.PERMISSION_NOT_SPECIFIED {
		if reject := authorize(permissionKey, metadata, permissions); reject {
			// TODO: Log this incident.
			http.Error(w, "", http.StatusUnauthorized)
			return "", nil, true
		}
	}

	return username, permissions, false
}

func authorize(permissionKey models.PermissionKey, metadata *string, permissions []models.Permission) (reject bool) {
	for _, permission := range permissions {
		if permission.Key != permissionKey {
			continue
		}
		if metadata != nil && *metadata != permission.Metadata {
			continue
		}

		return false
	}

	return true
}
