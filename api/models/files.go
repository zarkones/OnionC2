package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type File struct {
	ID                string `json:"id" gorm:"primaryKey"`
	UploadedByAgentID string `json:"uploadedByAgentId"`
	OriginalPath      string `json:"originalPath"`
	UploadedAt        int64  `json:"uploadedAt"`
	CreatedAt         int64  `json:"createdAt"`
}

func (f *File) BeforeCreate(tx *gorm.DB) (err error) {
	if f.ID == "" {
		f.ID = uuid.New().String()
	}
	f.CreatedAt = time.Now().UnixNano()
	return nil
}
