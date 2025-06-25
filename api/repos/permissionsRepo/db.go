package permissionsRepo

import (
	"api/db"
	"api/models"
)

func Get(permissionID string) (permission models.Permission, err error) {
	return permission, db.ORM.Where("id = ?", permissionID).First(&permission).Error
}

func GetMultipleByUsername(username string) (permissions []models.Permission, err error) {
	return permissions, db.ORM.Where("username = ?", username).Find(&permissions).Error
}

func GetMultipleByChannel(channel string) (permissions []models.Permission, err error) {
	return permissions, db.ORM.
		Where("metadata = ?", channel).
		Where("key = ?", models.PERMISSION_CHAT_LIST_CHANNEL_MESSAGES).
		Find(&permissions).Error
}

func GetMultiple() (permissions []models.Permission, err error) {
	return permissions, db.ORM.Find(&permissions).Error
}

func Insert(permission *models.Permission) (err error) {
	return db.ORM.Create(&permission).Error
}
