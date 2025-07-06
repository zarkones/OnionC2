package filesRepo

import (
	"api/db"
	"api/models"
	"time"
)

func Get(fileID string) (file models.File, err error) {
	return file, db.ORM.Where("id = ?", fileID).First(&file).Error
}

func GetMultiple() (files []models.File, err error) {
	return files, db.ORM.Find(&files).Error
}

func GetMultipleByAgentID(agentID string) (files []models.File, err error) {
	return files, db.ORM.Where("agent_id = ?", agentID).Find(&files).Error
}

func Insert(file *models.File) (err error) {
	return db.ORM.Create(&file).Error
}

func SetCompleted(id string) (err error) {
	file, err := Get(id)
	if err != nil {
		return err
	}
	if file.CompletedAt != 0 {
		return nil
	}
	file.CompletedAt = time.Now().UnixNano()
	return db.ORM.Save(&file).Error
}
