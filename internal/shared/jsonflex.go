package shared

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// FlexFloat64 unmarshals JSON numbers or numeric strings (Android sends strings).
type FlexFloat64 float64

func (f FlexFloat64) Float64() float64 { return float64(f) }

func (f *FlexFloat64) UnmarshalJSON(data []byte) error {
	n, err := parseJSONFloat(data)
	if err != nil {
		return err
	}
	*f = FlexFloat64(n)
	return nil
}

// FlexInt unmarshals JSON numbers or numeric strings (Android sends strings).
type FlexInt int

func (i FlexInt) Int() int { return int(i) }

func (i *FlexInt) UnmarshalJSON(data []byte) error {
	n, err := parseJSONFloat(data)
	if err != nil {
		return err
	}
	*i = FlexInt(int(n))
	return nil
}

func parseJSONFloat(data []byte) (float64, error) {
	data = bytes.TrimSpace(data)
	if len(data) == 0 || bytes.Equal(data, []byte("null")) {
		return 0, nil
	}
	if data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return 0, err
		}
		s = strings.TrimSpace(s)
		if s == "" {
			return 0, nil
		}
		n, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid number %q", s)
		}
		return n, nil
	}
	var n float64
	if err := json.Unmarshal(data, &n); err != nil {
		return 0, fmt.Errorf("invalid number %s", string(data))
	}
	return n, nil
}

func FlexIntPtr(v *FlexInt) *int {
	if v == nil {
		return nil
	}
	n := v.Int()
	return &n
}

func ParseFlexibleDate(dStr string) (time.Time, error) {
	dStr = strings.TrimSpace(dStr)
	if dStr == "" {
		return time.Time{}, fmt.Errorf("date is required")
	}

	layouts := []string{
		time.RFC3339,
		"2006-01-02",
		"2006/01/02",
		"02/01/2006",
		"02-01-2006",
		"2006-01-02T15:04:05",
		"02/01/2006 15:04:05",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, dStr); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid date %q", dStr)
}
