package service

import (
	"strings"
	"testing"
)

func TestValidateConfig(t *testing.T) {
	cases := []struct {
		name    string
		raw     string
		wantErr string
	}{
		{"valid", `{"serverUrl":"https://api.example.com","content":{"auth":{"title":"Hi","list":["a","b"]}}}`, ""},
		{"valid trailing slash", `{"serverUrl":"https://api.example.com/","content":{}}`, ""},
		{"not an object", `[1,2]`, "must be a JSON object"},
		{"missing serverUrl", `{"content":{}}`, `"serverUrl" must be a non-empty string`},
		{"relative serverUrl", `{"serverUrl":"api.example.com","content":{}}`, "absolute http(s) URL"},
		{"serverUrl with path", `{"serverUrl":"https://api.example.com/api/v1","content":{}}`, "no path"},
		{"missing content", `{"serverUrl":"https://api.example.com"}`, `"content" is missing`},
		{"content not object", `{"serverUrl":"https://api.example.com","content":"x"}`, `"content" must be an object`},
		{"number value", `{"serverUrl":"https://api.example.com","content":{"a":{"b":1}}}`, "content.a.b must be text"},
		{"list of objects", `{"serverUrl":"https://api.example.com","content":{"a":[{"x":"y"}]}}`, "content.a[0] must be text"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ValidateConfig([]byte(tc.raw))
			if tc.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("want error containing %q, got %v", tc.wantErr, err)
			}
		})
	}
}

func TestValidateConfigKeepsKeyOrder(t *testing.T) {
	out, err := ValidateConfig([]byte(`{"serverUrl":"https://a.com","content":{"z":"1","a":"2"}}`))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Index(string(out), `"z"`) > strings.Index(string(out), `"a"`) {
		t.Fatalf("key order changed:\n%s", out)
	}
}
