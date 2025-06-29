package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PermissionKey int

const (
	// Agents permissions.
	PERMISSION_NOT_SPECIFIED = PermissionKey(iota)
	PERMISSION_AGENTS_LIST
	PERMISSION_AGENTS_LIST_MESSAGES
	PERMISSION_AGENTS_INSERT_MESSAGE

	// Chat permissions. (if uncommented they need to go down no to mess the key int order)
	PERMISSION_CHAT_LIST_CHANNELS
	PERMISSION_CHAT_LIST_CHANNEL_MESSAGES
	PERMISSION_CHAT_INSERT_CHANNEL
	PERMISSION_CHAT_INSERT_CHANNEL_MESSAGE
	PERMISSION_CHAT_DELETE_CHANNEL
	PERMISSION_CHAT_DELETE_CHANNEL_MESSAGE

	// Operator management permissions.
	PERMISSION_OPERATORS_LIST
	PERMISSION_OPERATORS_INSERT
	PERMISSION_OPERATORS_DELETE

	// Permissions management permissions.
	PERMISSION_INSERT
	PERMISSION_DELETE
)

type Permission struct {
	ID        string        `json:"id" gorm:"primaryKey"`
	Key       PermissionKey `json:"key"`
	Username  string        `json:"username"`
	Metadata  string        `json:"metadata"`
	CreatedAt int64         `json:"createdAt"`
}

func (o *Permission) BeforeCreate(tx *gorm.DB) (err error) {
	if o.ID == "" {
		o.ID = uuid.New().String()
	}
	o.CreatedAt = time.Now().UnixNano()
	return nil
}
