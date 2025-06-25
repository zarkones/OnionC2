package models

import (
	"time"

	"gorm.io/gorm"
)

type Operator struct {
	Username               string `json:"username" gorm:"primaryKey"`
	PublicKeyHex           string `json:"publicKeyHex"`
	EncryptedPrivateKeyHex string `json:"encryptedPrivateKeyHex"`
	CreatedAt              int64  `json:"createdAt"`
}

func (o *Operator) BeforeCreate(tx *gorm.DB) (err error) {
	o.CreatedAt = time.Now().UnixNano()
	return nil
}
