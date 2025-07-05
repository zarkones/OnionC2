package messagesRepo

import (
	"api/db"
	"api/models"
	"errors"
	"time"
)

var ErrMsgRespPopulated = errors.New("message's response is already populated")

func GetLatestFS(agentID string) (messages models.Message, err error) {
	return messages, db.ORM.
		Where("agent_id = ?", agentID).
		Where("request = ? OR request LIKE ?", "/ls", "/ls|%").
		Where("response != ?", "").
		Order("created_at DESC").
		First(&messages).Error
}

func GetLatestRequestedFS(agentID string) (messages models.Message, err error) {
	return messages, db.ORM.
		Where("agent_id = ?", agentID).
		Where("request = ? OR request LIKE ?", "/ls", "/ls|%").
		Order("created_at DESC").
		First(&messages).Error
}

func GetMultiple(agentID string, offset, limit int) (messages []models.Message, err error) {
	return messages, db.ORM.
		Where("agent_id = ?", agentID).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&messages).Error
}

func GetMultipleBefore(agentID string, before int64, limit int) (messages []models.Message, err error) {
	return messages, db.ORM.
		Where("agent_id = ?", agentID).
		Where("created_at < ?", before).
		Order("created_at DESC").
		Limit(limit).
		Find(&messages).Error
}

func GetMultipleAfter(agentID string, after int64, limit int) (messages []models.Message, err error) {
	return messages, db.ORM.
		Where("agent_id = ?", agentID).
		Where("created_at > ?", after).
		Order("created_at ASC").
		Limit(limit).
		Find(&messages).Error
}

func Get(messageID string) (messages models.Message, err error) {
	return messages, db.ORM.Where("id = ?", messageID).First(&messages).Error
}

func GetMultipleByIDs(messageIDs []string) (messages []models.Message, err error) {
	return messages, db.ORM.
		Find(&messages, messageIDs).Error
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
