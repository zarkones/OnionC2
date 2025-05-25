package messagesRepo

import (
	"api/db"
	"api/models"
	"errors"
)

var ErrMsgRespPopulated = errors.New("message's response is already populated")

func GetMultiple(agentID string) (messages []models.Message, err error) {
	return messages, db.ORM.Where("agent_id = ?", agentID).Find(&messages).Error
}

func GetMultipleForAgent(agentID string) (messages []models.Message, err error) {
	return messages, db.ORM.Where("agent_id = ?", agentID).Where("response = ?", "").Find(&messages).Error
}

func Insert(message *models.Message) (err error) {
	return db.ORM.Create(&message).Error
}

func UpdateResponse(messageID, response string) (err error) {
	var message models.Message
	if err := db.ORM.Table("messages").Where("id = ?", messageID).First(&message).Error; err != nil {
		return err
	}
	if message.Response != "" {
		return ErrMsgRespPopulated
	}
	message.Response = response
	return db.ORM.Save(message).Error
}
