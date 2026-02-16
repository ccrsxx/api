package ipinfo

import (
	"testing"
)

func TestDefaultClient(t *testing.T) {
	client := DefaultClient()

	if client == nil {
		t.Fatal("expected ipinfo client to be initialized, got nil")
	}
}
