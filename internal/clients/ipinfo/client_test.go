package ipinfo_test

import (
	"testing"

	"github.com/ccrsxx/api/internal/clients/ipinfo"
)

func TestNew(t *testing.T) {
	client := ipinfo.NewClient("test-token")

	if client == nil {
		t.Fatal("got nil, want ipinfo client to be initialized")
	}
}
