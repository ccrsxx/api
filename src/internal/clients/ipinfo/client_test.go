package ipinfo

import (
	"testing"
)

func TestDefaultClient(t *testing.T) {
	c1 := DefaultClient()

	if c1 == nil {
		t.Fatal("expected ipinfo client to be initialized, got nil")
	}

	c2 := DefaultClient()

	if c1 != c2 {
		t.Error("Client() did not return the singleton instance")
	}
}
