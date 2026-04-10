package ipinfo

import (
	"testing"
)

func TestNew(t *testing.T) {
	client := NewClient("test-token")

	if client == nil {
		t.Fatal("got nil, want ipinfo client to be initialized")
	}
}
