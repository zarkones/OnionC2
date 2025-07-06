package c2ctrl

import (
	"api/core"
	"testing"
)

func TestFileNameGrabbing(t *testing.T) {
	paths := []string{
		"something/somewhere/xyz.id",
		"C:\\Users\\blablabla\\xyz.id",
		"xyz.id",
	}

	for _, path := range paths {
		if core.GetFileName(path) != "xyz.id" {
			t.FailNow()
		}
	}
}
