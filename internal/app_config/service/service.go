package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"expent-backend/internal/infrastructure/cloudflare"
)

var (
	// ErrNotConfigured means the Cloudflare KV credentials aren't set on the server.
	ErrNotConfigured = errors.New("app config storage isn't configured: set CF_ACCOUNT_ID, CF_KV_NAMESPACE_ID and CF_API_TOKEN")
	// ErrVersionConflict means the stored config changed after the admin loaded it.
	ErrVersionConflict = errors.New("the config was changed by someone else since you loaded it; reload and apply your edits again")
)

// AppConfig is the config as the panel sees it.
type AppConfig struct {
	Config json.RawMessage `json:"config"`
	// Version identifies the stored value; send it back as baseVersion when saving.
	Version string `json:"version"`
	// StoredInKV is false when KV is still empty and the config was read from the live public URL instead.
	StoredInKV bool `json:"storedInKv"`
}

type Service struct {
	kv        *cloudflare.KVClient
	key       string
	publicURL string
	http      *http.Client
}

func NewService(kv *cloudflare.KVClient, key, publicURL string) *Service {
	return &Service{kv: kv, key: key, publicURL: publicURL, http: &http.Client{Timeout: 20 * time.Second}}
}

func (s *Service) Get(ctx context.Context) (*AppConfig, error) {
	if !s.kv.Configured() {
		return nil, ErrNotConfigured
	}
	value, err := s.kv.Get(ctx, s.key)
	if errors.Is(err, cloudflare.ErrKeyNotFound) {
		// First run: KV hasn't been seeded yet, so start from what the app currently receives
		live, liveErr := s.fetchPublic(ctx)
		if liveErr != nil {
			return nil, fmt.Errorf("no config in KV yet, and the live config couldn't be read: %w", liveErr)
		}
		return &AppConfig{Config: live, Version: "", StoredInKV: false}, nil
	}
	if err != nil {
		return nil, err
	}
	return &AppConfig{Config: value, Version: version(value), StoredInKV: true}, nil
}

// Update validates and stores a new config. baseVersion must match the stored value
// (empty when KV is still empty), so two admins can't silently overwrite each other.
func (s *Service) Update(ctx context.Context, raw []byte, baseVersion string) (*AppConfig, error) {
	if !s.kv.Configured() {
		return nil, ErrNotConfigured
	}
	formatted, err := ValidateConfig(raw)
	if err != nil {
		return nil, &ValidationError{err}
	}

	current, err := s.kv.Get(ctx, s.key)
	switch {
	case errors.Is(err, cloudflare.ErrKeyNotFound):
		current = nil
	case err != nil:
		return nil, err
	}
	if version(current) != baseVersion {
		return nil, ErrVersionConflict
	}

	// Keep the previous value so a bad publish can be rolled back from the Cloudflare dashboard
	if current != nil {
		if err := s.kv.Put(ctx, s.key+"_previous", current); err != nil {
			return nil, fmt.Errorf("couldn't back up the current config: %w", err)
		}
	}
	if err := s.kv.Put(ctx, s.key, formatted); err != nil {
		return nil, err
	}
	return &AppConfig{Config: formatted, Version: version(formatted), StoredInKV: true}, nil
}

func (s *Service) fetchPublic(ctx context.Context) (json.RawMessage, error) {
	if s.publicURL == "" {
		return nil, errors.New("CONFIG_PUBLIC_URL isn't set")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.publicURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d from %s", resp.StatusCode, s.publicURL)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, MaxConfigBytes+1))
	if err != nil {
		return nil, err
	}
	if !json.Valid(body) {
		return nil, fmt.Errorf("%s didn't return valid JSON", s.publicURL)
	}
	return body, nil
}

// version is a content hash of the stored value; "" means nothing is stored.
func version(value []byte) string {
	if value == nil {
		return ""
	}
	sum := sha256.Sum256(value)
	return hex.EncodeToString(sum[:])
}

// ValidationError marks errors caused by the submitted config itself.
type ValidationError struct{ Err error }

func (e *ValidationError) Error() string { return e.Err.Error() }
func (e *ValidationError) Unwrap() error { return e.Err }
