package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChannelMessage struct {
	ID        string `json:"id" gorm:"primaryKey"`
	Channel   string `json:"channel"`
	Data      string `json:"data"`
	CreatedAt int64  `json:"createdAt"`
}

func (c *ChannelMessage) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	c.CreatedAt = time.Now().UnixNano()
	return nil
}
