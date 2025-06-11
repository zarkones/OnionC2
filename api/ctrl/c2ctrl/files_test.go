package c2ctrl

import (
	"testing"
)

func TestXxx(t *testing.T) {
	paths := []string{
		"something/somewhere/xyz.id",
		"C:\\Users\\blablabla\\xyz.id",
		"xyz.id",
	}

	for _, path := range paths {
		if getFileName(path) != "xyz.id" {
			t.FailNow()
		}
	}
}
