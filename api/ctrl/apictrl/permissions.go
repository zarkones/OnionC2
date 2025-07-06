package apictrl

import (
	"api/models"
	"api/repos/permissionsRepo"
	"encoding/json"
	"net/http"
)

type ExtendedPermission struct {
	models.Permission
	Acquired bool `json:"acquired"`
}

type GetPermissionsRespCtx map[models.PermissionKey]ExtendedPermission

func GetPermissions(w http.ResponseWriter, r *http.Request) {
	_, _, reject := authenticateAndAuthorize(w, r, models.PERMISSION_LIST, nil)
	if reject {
		return
	}

	operatorUsername := r.PathValue("operatorUsername")

	permissions, err := permissionsRepo.GetMultipleByUsername(operatorUsername)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(permissions) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	operatorPermissionMap := make(map[models.PermissionKey]models.Permission, len(permissions))
	for _, permission := range permissions {
		operatorPermissionMap[permission.Key] = permission
	}

	permissionMap := make(GetPermissionsRespCtx, len(models.AllPermissions))
	for _, permissionKey := range models.AllPermissions {
		if _, ok := operatorPermissionMap[permissionKey]; ok {
			permissionMap[permissionKey] = ExtendedPermission{
				Permission: operatorPermissionMap[permissionKey],
				Acquired:   true,
			}
			continue
		}
		permissionMap[permissionKey] = ExtendedPermission{
			Permission: models.Permission{
				Key:      permissionKey,
				Username: operatorUsername,
			},
			Acquired: false,
		}
	}

	if err := json.NewEncoder(w).Encode(&permissionMap); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func InsertPermission(w http.ResponseWriter, r *http.Request) {
	_, _, reject := authenticateAndAuthorize(w, r, models.PERMISSION_INSERT, nil)
	if reject {
		return
	}

	var newPermission models.Permission

	if err := json.NewDecoder(r.Body).Decode(&newPermission); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(newPermission.ID) != 0 {
		http.Error(w, "id cannot be provided", http.StatusUnprocessableEntity)
		return
	}
	if len(newPermission.Username) == 0 {
		http.Error(w, "username must be provided", http.StatusUnprocessableEntity)
		return
	}
	if newPermission.Key == models.PERMISSION_NOT_SPECIFIED {
		http.Error(w, "permission key cannot be PERMISSION_NOT_SPECIFIED", http.StatusUnprocessableEntity)
		return
	}

	if err := permissionsRepo.Insert(&newPermission); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func DeletePermission(w http.ResponseWriter, r *http.Request) {
	_, _, reject := authenticateAndAuthorize(w, r, models.PERMISSION_DELETE, nil)
	if reject {
		return
	}

	permissionID := r.PathValue("permissionID")

	if err := permissionsRepo.Delete(permissionID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
