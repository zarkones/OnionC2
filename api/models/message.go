package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Message struct {
	ID                  string `json:"id" gorm:"primaryKey"`
	AgentID             string `json:"agentId"`
	Request             string `json:"request"`
	Response            string `json:"response"`
	CreatedAt           int64  `json:"createdAt"`
}

func (m *Message) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	m.CreatedAt = time.Now().UnixNano()
	return nil
}
