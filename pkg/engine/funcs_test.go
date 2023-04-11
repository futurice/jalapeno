package engine

import (
	"testing"

	"github.com/gofrs/uuid"
)

func TestStableRandomAlpha(t *testing.T) {
	u, err := uuid.NewV4()
	if err != nil {
		t.Error(err)
	}

	val1 := stableRandomAlphanumeric(6, u)
	if len(val1) != 6 {
		t.Errorf("Expected 6 alphanumeric, got '%s'", val1)
	}

	val2 := stableRandomAlphanumeric(6, u)
	if len(val2) != 6 {
		t.Errorf("Expected 6 alphanumeric, got '%s'", val2)
	}
	if val1 == val2 {
		t.Errorf("Expected different values, got '%s'", val1)
	}

	resetRngs()

	val3 := stableRandomAlphanumeric(6, u)
	if val3 != val1 {
		t.Error("Expected the same sequence after reset")
	}
}
