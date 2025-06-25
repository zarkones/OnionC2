package channelMessagesRepo

import (
	"api/db"
	"api/models"
)

func GetMultiple(channelName string, offset, limit int) (channelMessages []models.ChannelMessage, err error) {
	return channelMessages, db.ORM.Where("channel = ?", channelName).Offset(offset).Limit(limit).Find(&channelMessages).Error
}

func Insert(channelMessage *models.ChannelMessage) (err error) {
	return db.ORM.Create(&channelMessage).Error
}

func DeleteMultiple(channelName string) (err error) {
	return db.ORM.Where("channel = ?", channelName).Delete(models.Channel{}).Error
}
