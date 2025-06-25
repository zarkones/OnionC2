package models

import (
	"time"

	"gorm.io/gorm"
)

type Channel struct {
	Name        string `json:"name" gorm:"primaryKey"`
	Description string `json:"description"`
	CreatedAt   int64  `json:"createdAt"`
}

func (c *Channel) BeforeCreate(tx *gorm.DB) (err error) {
	c.CreatedAt = time.Now().UnixNano()
	return nil
}
