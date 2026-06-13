package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"expent-backend/internal/parse_transaction/model"

	"google.golang.org/genai"
)

type Service struct {
	geminiKey string
}

func NewService(geminiKey string) *Service {
	return &Service{geminiKey: geminiKey}
}

// ParseTransaction calls Gemini to extract structured transaction data from a raw SMS or text.
func (s *Service) ParseTransaction(ctx context.Context, rawText string) (*model.ParsedTransaction, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  s.geminiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	prompt := buildPrompt(rawText)

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash",
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("gemini request failed: %w", err)
	}

	rawJSON := extractJSON(result.Text())
	if rawJSON == "" {
		return nil, fmt.Errorf("gemini returned empty or non-JSON response")
	}

	var parsed model.ParsedTransaction
	if err := json.Unmarshal([]byte(rawJSON), &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse gemini response as JSON: %w\nRaw: %s", err, rawJSON)
	}

	return &parsed, nil
}

// buildPrompt constructs the prompt sent to Gemini.
func buildPrompt(rawText string) string {
	return fmt.Sprintf(`You are a financial transaction parser. Given a raw bank SMS, email, or text, extract structured transaction data.

Respond ONLY with a valid JSON object in exactly this format (no markdown, no explanation):
{
  "amount": <float>,
  "type": "<expense|income>",
  "description": "<merchant or narration>",
  "date": "<ISO 8601 date if detected, else empty string>",
  "currency": "<currency code, e.g. INR>",
  "confidence": <0.0 to 1.0>
}

Rules:
- "type" must be either "expense" (money deducted/debited) or "income" (money credited/received).
- "amount" must be a positive number.
- If you cannot determine a field, use an empty string or 0.
- "confidence" represents how confident you are in the extraction (1.0 = very confident).

Raw text to parse:
"%s"`, rawText)
}

// extractJSON tries to extract a JSON object from the model's response.
func extractJSON(text string) string {
	text = strings.TrimSpace(text)
	// Strip markdown code fences if present
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimPrefix(text, "```")
	text = strings.TrimSuffix(text, "```")
	text = strings.TrimSpace(text)

	start := strings.Index(text, "{")
	end := strings.LastIndex(text, "}")
	if start == -1 || end == -1 || end < start {
		return ""
	}
	return text[start : end+1]
}
