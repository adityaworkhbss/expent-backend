package cloudflare

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// ErrKeyNotFound is returned by Get when the key doesn't exist in the namespace.
var ErrKeyNotFound = errors.New("kv key not found")

// KVClient reads and writes values in one Workers KV namespace through the Cloudflare REST API.
type KVClient struct {
	accountID   string
	namespaceID string
	apiToken    string
	http        *http.Client
}

func NewKVClient(accountID, namespaceID, apiToken string) *KVClient {
	return &KVClient{
		accountID:   accountID,
		namespaceID: namespaceID,
		apiToken:    apiToken,
		http:        &http.Client{Timeout: 20 * time.Second},
	}
}

// Configured reports whether all credentials needed to reach KV are set.
func (k *KVClient) Configured() bool {
	return k.accountID != "" && k.namespaceID != "" && k.apiToken != ""
}

func (k *KVClient) valueURL(key string) string {
	return fmt.Sprintf(
		"https://api.cloudflare.com/client/v4/accounts/%s/storage/kv/namespaces/%s/values/%s",
		url.PathEscape(k.accountID), url.PathEscape(k.namespaceID), url.PathEscape(key),
	)
}

// Get returns the raw value stored under key, or ErrKeyNotFound.
func (k *KVClient) Get(ctx context.Context, key string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, k.valueURL(key), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+k.apiToken)

	resp, err := k.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("kv get %q: %w", key, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("kv get %q: reading body: %w", key, err)
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrKeyNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("kv get %q: %s", key, apiError(resp.StatusCode, body))
	}
	return body, nil
}

// Put stores value under key, replacing any existing value.
func (k *KVClient) Put(ctx context.Context, key string, value []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, k.valueURL(key), bytes.NewReader(value))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+k.apiToken)
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")

	resp, err := k.http.Do(req)
	if err != nil {
		return fmt.Errorf("kv put %q: %w", key, err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("kv put %q: %s", key, apiError(resp.StatusCode, body))
	}
	return nil
}

// apiError turns a Cloudflare API error envelope into a short message.
func apiError(status int, body []byte) string {
	var envelope struct {
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}
	if json.Unmarshal(body, &envelope) == nil && len(envelope.Errors) > 0 {
		return fmt.Sprintf("HTTP %d: %s", status, envelope.Errors[0].Message)
	}
	return fmt.Sprintf("HTTP %d", status)
}
