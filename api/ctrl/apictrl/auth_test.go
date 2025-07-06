package apictrl

import (
	"api/models"
	"testing"
)

func TestAuthorization(t *testing.T) {
	if reject := authorize(models.PERMISSION_AGENTS_LIST, nil, []models.Permission{
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
	}); reject {
		t.Log("sufficient permissions weren't authorized")
		t.FailNow()
	}

	if reject := authorize(models.PERMISSION_AGENTS_LIST, nil, []models.Permission{
		{
			Key: models.PERMISSION_CHAT_DELETE_CHANNEL,
		},
		{
			Key: models.PERMISSION_CHAT_INSERT_CHANNEL_MESSAGE,
		},
		{
			Key: models.PERMISSION_INSERT,
		},
	}); !reject {
		t.Log("insufficient permissions were authorized")
		t.FailNow()
	}

	metadata := "XYZ123"

	if reject := authorize(models.PERMISSION_CHAT_LIST_CHANNEL_MESSAGES, &metadata, []models.Permission{
		{
			Key:      models.PERMISSION_CHAT_LIST_CHANNEL_MESSAGES,
			Metadata: "XYZ123",
		},
		{
			Key: models.PERMISSION_INSERT,
		},
	}); reject {
		t.Log("sufficient permissions weren't authorized")
		t.FailNow()
	}

	if reject := authorize(models.PERMISSION_CHAT_LIST_CHANNEL_MESSAGES, &metadata, []models.Permission{
		{
			Key:      models.PERMISSION_CHAT_LIST_CHANNEL_MESSAGES,
			Metadata: "somethingdifferent",
		},
		{
			Key: models.PERMISSION_INSERT,
		},
	}); !reject {
		t.Log("insufficient permissions were authorized")
		t.FailNow()
	}
}
