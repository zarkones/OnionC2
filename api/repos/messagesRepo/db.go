package messagesRepo

import (
	"api/db"
	"api/models"
	"errors"
	"time"
)

var ErrMsgRespPopulated = errors.New("message's response is already populated")

func GetMultiple(agentID string, offset, limit int) (messages []models.Message, err error) {
	return messages, db.ORM.
		Where("agent_id = ?", agentID).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&messages).Error
}

func GetMultipleSince(agentID string, since int64, limit int) (messages []models.Message, err error) {
	return messages, db.ORM.
		Where("agent_id = ?", agentID).
		Where("created_at < ?", since).
		Order("created_at DESC").
		Limit(limit).
		Find(&messages).Error
}

func GetMultipleForAgent(agentID string) (messages []models.Message, err error) {
	return messages, db.ORM.Where("agent_id = ?", agentID).Where("response = ?", "").Find(&messages).Error
}

func Insert(message *models.Message) (err error) {
	return db.ORM.Create(&message).Error
}

func UpdateResponse(messageID, response string) (agentID string, err error) {
	var message models.Message
	if err := db.ORM.Table("messages").Where("id = ?", messageID).First(&message).Error; err != nil {
		return "", err
	}
	if message.Response != "" {
		return "", ErrMsgRespPopulated
	}
	message.Response = response
	message.UpdatedAt = time.Now().UnixNano()
	return message.AgentID, db.ORM.Save(message).Error
}
