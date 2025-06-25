package apictrl

import (
	"api/models"
	"testing"
)

func TestAuthorization(t *testing.T) {
	if authorized := authorize(models.PERMISSION_AGENTS_LIST, []models.Permission{
		{
			Key: models.PERMISSION_CHAT_DELETE_CHANNEL,
		},
		{
			Key: models.PERMISSION_CHAT_INSERT_CHANNEL_MESSAGE,
		},
		{
			Key: models.PERMISSION_AGENTS_LIST,
		},
		{
			Key: models.PERMISSION_INSERT,
		},
	}); !authorized {
		t.Log("sufficient permissions weren't authorized")
		t.FailNow()
	}

	if authorized := authorize(models.PERMISSION_AGENTS_LIST, []models.Permission{
		{
			Key: models.PERMISSION_CHAT_DELETE_CHANNEL,
		},
		{
			Key: models.PERMISSION_CHAT_INSERT_CHANNEL_MESSAGE,
		},
		{
			Key: models.PERMISSION_INSERT,
		},
	}); authorized {
		t.Log("insufficient permissions weren authorized")
		t.FailNow()
	}
}
