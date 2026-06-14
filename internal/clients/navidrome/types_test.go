package navidrome_test

import (
	"encoding/json"
	"encoding/xml"
	"testing"

	"github.com/ccrsxx/api/internal/clients/navidrome"
)

func TestReplayGain_MarshalXML(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		rg := navidrome.ReplayGain{}

		data, err := xml.Marshal(rg)

		if err != nil {
			t.Fatalf("got err: %v, want success", err)
		}

		if string(data) != "" {
			t.Errorf("got %q, want empty string for zero-value ReplayGain", string(data))
		}
	})

	t.Run("With Values", func(t *testing.T) {
		gain := 1.5

		rg := navidrome.ReplayGain{TrackGain: &gain}

		data, err := xml.Marshal(rg)

		if err != nil {
			t.Fatalf("got err: %v, want success", err)
		}

		if len(data) == 0 {
			t.Error("want non-empty XML for populated ReplayGain")
		}
	})
}

func TestItemDate_MarshalXML(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		d := navidrome.ItemDate{}

		data, err := xml.Marshal(d)

		if err != nil {
			t.Fatalf("got err: %v, want success", err)
		}

		if string(data) != "" {
			t.Errorf("got %q, want empty string for zero-value ItemDate", string(data))
		}
	})

	t.Run("With Values", func(t *testing.T) {
		d := navidrome.ItemDate{Year: 2024, Month: 6, Day: 14}

		data, err := xml.Marshal(d)

		if err != nil {
			t.Fatalf("got err: %v, want success", err)
		}

		if len(data) == 0 {
			t.Error("want non-empty XML for populated ItemDate")
		}
	})
}

func TestArray_MarshalJSON(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		a := navidrome.Array[string]{}

		data, err := json.Marshal(a)

		if err != nil {
			t.Fatalf("got err: %v, want success", err)
		}

		if string(data) != "[]" {
			t.Errorf("got %q, want [] for empty Array", string(data))
		}
	})

	t.Run("With Values", func(t *testing.T) {
		a := navidrome.Array[string]{"a", "b"}

		data, err := json.Marshal(a)

		if err != nil {
			t.Fatalf("got err: %v, want success", err)
		}

		if string(data) != `["a","b"]` {
			t.Errorf("got %q, want %q", string(data), `["a","b"]`)
		}
	})
}
