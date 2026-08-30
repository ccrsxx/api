package pixiv_test

import (
	"encoding/json/jsontext"
	"errors"
	"strings"
	"testing"

	"github.com/ccrsxx/api/internal/clients/pixiv"
)

type errReader struct{}

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}

func TestFlexibleID_UnmarshalJSONFrom(t *testing.T) {
	t.Run("Valid String", func(t *testing.T) {
		dec := jsontext.NewDecoder(strings.NewReader(`"123"`))

		var f pixiv.FlexibleID

		if err := f.UnmarshalJSONFrom(dec); err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		if f != "123" {
			t.Errorf("got %s, want 123", f)
		}
	})

	t.Run("Valid Number", func(t *testing.T) {
		dec := jsontext.NewDecoder(strings.NewReader(`456`))

		var f pixiv.FlexibleID

		if err := f.UnmarshalJSONFrom(dec); err != nil {
			t.Fatalf("unwanted error: %v", err)
		}

		if f != "456" {
			t.Errorf("got %s, want 456", f)
		}
	})

	t.Run("Invalid Token Kind", func(t *testing.T) {
		dec := jsontext.NewDecoder(strings.NewReader(`true`))

		var f pixiv.FlexibleID

		if err := f.UnmarshalJSONFrom(dec); err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("Decoder Read Error", func(t *testing.T) {
		dec := jsontext.NewDecoder(errReader{})

		var f pixiv.FlexibleID

		err := f.UnmarshalJSONFrom(dec)

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "flexible id read token error") {
			t.Errorf("got %v, want 'flexible id read token error'", err)
		}
	})
}
