package shared

import (
	"encoding/json"
	"testing"
	"time"
)

func TestFlexFloat64_StringAndNumber(t *testing.T) {
	var fromString FlexFloat64
	if err := json.Unmarshal([]byte(`"9000"`), &fromString); err != nil {
		t.Fatalf("string amount: %v", err)
	}
	if fromString.Float64() != 9000 {
		t.Fatalf("got %v", fromString.Float64())
	}

	var fromNumber FlexFloat64
	if err := json.Unmarshal([]byte(`9000.5`), &fromNumber); err != nil {
		t.Fatalf("number amount: %v", err)
	}
	if fromNumber.Float64() != 9000.5 {
		t.Fatalf("got %v", fromNumber.Float64())
	}
}

func TestFlexInt_String(t *testing.T) {
	var n FlexInt
	if err := json.Unmarshal([]byte(`"9"`), &n); err != nil {
		t.Fatalf("%v", err)
	}
	if n.Int() != 9 {
		t.Fatalf("got %d", n.Int())
	}
}

func TestParseFlexibleDate_DDMMYYYY(t *testing.T) {
	got, err := ParseFlexibleDate("07/10/2026")
	if err != nil {
		t.Fatal(err)
	}
	want := time.Date(2026, 10, 7, 0, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Fatalf("got %v want %v", got, want)
	}
}
