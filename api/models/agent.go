package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Agent struct {
	ID       string `json:"id" gorm:"primaryKey" sql:"type:text"`
	Hostname string `json:"hostname"`
	OS       string `json:"os"`
	Arch     string `json:"arch"`
}

func (agent *Agent) BeforeCreate(tx *gorm.DB) (err error) {
	agent.ID = uuid.New().String()
	return
}
