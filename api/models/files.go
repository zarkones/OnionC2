package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	ORDER_UPLOAD = iota
	ORDER_DOWNLOAD
)

type File struct {
	ID                string `json:"id" gorm:"primaryKey"`
	UploadedByAgentID string `json:"uploadedByAgentId"`
	OriginalPath      string `json:"originalPath"`
	Order             int    `json:"order"`
	CompletedAt       int64  `json:"completedAt"`
	CreatedAt         int64  `json:"createdAt"`
}

func (f *File) BeforeCreate(tx *gorm.DB) (err error) {
	if f.ID == "" {
		f.ID = uuid.New().String()
	}
	f.CreatedAt = time.Now().UnixNano()
	return nil
}
