package apictrl

import (
	"api/core/crypto"
	"api/models"
	"net/http"
)

func authenticate(w http.ResponseWriter, r *http.Request, permissionKey models.PermissionKey) (username string, permissions []models.Permission, reject bool) {
	token := r.Header.Get("Authorization")

	username, permissions, err := crypto.VerifyToken(token)
	if err != nil {
		// TODO: Log this incident.
		http.Error(w, "", http.StatusUnauthorized)
		return "", nil, true
	}

	if permissionKey != models.PERMISSION_NOT_SPECIFIED {
		if hasRequiredPermissions := authorize(permissionKey, permissions); !hasRequiredPermissions {
			// TODO: Log this incident.
			http.Error(w, "", http.StatusUnauthorized)
			return "", nil, true
		}
	}

	return username, permissions, false
}

func authorize(permissionKey models.PermissionKey, permissions []models.Permission) (hasRequiredPermissions bool) {
	for _, permission := range permissions {
		if permission.Key != permissionKey {
			continue
		}

		return true
	}

	return false
}
