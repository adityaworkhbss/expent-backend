package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// MaxConfigBytes caps the config size; KV allows far more, but the app downloads it on every change.
const MaxConfigBytes = 512 * 1024

// ValidateConfig checks that raw is a config the Android app can use:
// a JSON object with a valid http(s) "serverUrl" and a "content" object whose
// values are strings, lists of strings, or nested objects of those.
// It returns the config re-indented for storage.
func ValidateConfig(raw []byte) ([]byte, error) {
	if len(raw) > MaxConfigBytes {
		return nil, fmt.Errorf("config is too large (%d KB, max %d KB)", len(raw)/1024, MaxConfigBytes/1024)
	}

	var root map[string]json.RawMessage
	if err := json.Unmarshal(raw, &root); err != nil {
		return nil, fmt.Errorf("config must be a JSON object: %v", err)
	}

	var serverURL string
	if err := json.Unmarshal(root["serverUrl"], &serverURL); err != nil || strings.TrimSpace(serverURL) == "" {
		return nil, fmt.Errorf(`"serverUrl" must be a non-empty string`)
	}
	u, err := url.Parse(strings.TrimSpace(serverURL))
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return nil, fmt.Errorf(`"serverUrl" must be an absolute http(s) URL, got %q`, serverURL)
	}
	if u.Path != "" && u.Path != "/" {
		return nil, fmt.Errorf(`"serverUrl" must be just the server origin (no path); the app adds /api/v1 itself`)
	}

	content, ok := root["content"]
	if !ok {
		return nil, fmt.Errorf(`"content" is missing`)
	}
	var contentValue interface{}
	if err := json.Unmarshal(content, &contentValue); err != nil {
		return nil, fmt.Errorf(`"content" is not valid JSON: %v`, err)
	}
	if _, isObject := contentValue.(map[string]interface{}); !isObject {
		return nil, fmt.Errorf(`"content" must be an object`)
	}
	if err := validateText(contentValue, "content"); err != nil {
		return nil, err
	}

	var out bytes.Buffer
	if err := json.Indent(&out, raw, "", "  "); err != nil {
		return nil, err
	}
	out.WriteByte('\n')
	return out.Bytes(), nil
}

func validateText(v interface{}, path string) error {
	switch value := v.(type) {
	case string:
		return nil
	case map[string]interface{}:
		for key, child := range value {
			if err := validateText(child, path+"."+key); err != nil {
				return err
			}
		}
		return nil
	case []interface{}:
		for i, item := range value {
			if _, isString := item.(string); !isString {
				return fmt.Errorf("%s[%d] must be text; lists can only contain text", path, i)
			}
		}
		return nil
	default:
		return fmt.Errorf("%s must be text, a list of text, or a group; got %T", path, v)
	}
}
