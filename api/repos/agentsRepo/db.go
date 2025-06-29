package agentsRepo

import (
	"api/db"
	"api/models"
	"time"
)

func Get(agentID string) (agent models.Agent, err error) {
	return agent, db.ORM.Where("id = ?", agentID).First(&agent).Error
}

func GetMultiple() (agents []models.Agent, err error) {
	return agents, db.ORM.Find(&agents).Error
}

func Insert(agent *models.Agent) (err error) {
	return db.ORM.Create(&agent).Error
}

func UpdateLastSeen(agentID string) (err error) {
	agent, err := Get(agentID)
	if err != nil {
		return err
	}
	agent.LastSeen = time.Now().UnixNano()
	return db.ORM.Save(agent).Error
}
