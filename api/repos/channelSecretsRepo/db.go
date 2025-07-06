package channelSecretsRepo

import (
	"api/db"
	"api/models"
)

func Get(id string) (channelSecret models.ChannelSecret, err error) {
	return channelSecret, db.ORM.Where("id = ?", id).First(&channelSecret).Error
}

func GetByChannelAndOperatorNames(channelName, operatorUsername string) (channelSecret models.ChannelSecret, err error) {
	return channelSecret, db.ORM.
		Where("recipient_operator_username = ?", operatorUsername).
		Where("channel = ?", channelName).
		First(&channelSecret).Error
}

func GetMultiple() (channelSecrets []models.ChannelSecret, err error) {
	return channelSecrets, db.ORM.Find(&channelSecrets).Error
}

func Insert(channelSecret *models.ChannelSecret) (err error) {
	return db.ORM.Create(&channelSecret).Error
}

func Delete(id string) (err error) {
	return db.ORM.Where("id = ?", id).Delete(models.ChannelSecret{}).Error
}
