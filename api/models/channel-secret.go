package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChannelSecret struct {
	ID                           string `json:"id" gorm:"primaryKey"`
	RecipientOperatorUsername    string `json:"recipientOperatorUsername"`
	HexEncodedRsaEncryptedAesKey string `json:"hexEncodedRsaEncryptedAesKey"`
	CreatedBy                    string `json:"createdBy"`
	CreatedAt                    int64  `json:"createdAt"`
}

func (c *ChannelSecret) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	c.CreatedAt = time.Now().UnixNano()
	return nil
}
