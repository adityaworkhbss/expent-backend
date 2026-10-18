package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	accountrepo "expent-backend/internal/account/repository"
	categoryrepo "expent-backend/internal/category/repository"
	"expent-backend/internal/parse_transaction/model"

	"google.golang.org/genai"
)

type Service struct {
	geminiKey   string
	geminiModel string
	categories  categoryrepo.Repository
	accounts    accountrepo.Repository
}

func NewService(geminiKey, geminiModel string, categories categoryrepo.Repository, accounts accountrepo.Repository) *Service {
	return &Service{
		geminiKey:   geminiKey,
		geminiModel: geminiModel,
		categories:  categories,
		accounts:    accounts,
	}
}

func (s *Service) ParseTransaction(ctx context.Context, userID, rawText, today string, reqCats, reqAccs []model.CatalogItem) (*model.ParseResult, error) {
	categories := reqCats
	if len(categories) == 0 && s.categories != nil {
		dbCats, err := s.categories.ListCategories(userID)
		if err != nil {
			return nil, fmt.Errorf("failed to load categories: %w", err)
		}
		for _, c := range dbCats {
			categories = append(categories, model.CatalogItem{ID: c.ID, Name: c.Name, Type: c.Type})
		}
	}

	accounts := reqAccs
	if len(accounts) == 0 && s.accounts != nil {
		dbAccs, err := s.accounts.ListAccounts(userID)
		if err != nil {
			return nil, fmt.Errorf("failed to load accounts: %w", err)
		}
		for _, a := range dbAccs {
			accounts = append(accounts, model.CatalogItem{ID: a.ID, Name: a.Name, Type: a.Type})
		}
	}

	raw, err := s.callGemini(ctx, rawText, today, categories, accounts)
	if err != nil {
		return nil, err
	}

	result := finalizeParse(*raw, categories, accounts, today)
	return &result, nil
}

func (s *Service) callGemini(ctx context.Context, rawText, today string, categories, accounts []model.CatalogItem) (*model.GeminiParsed, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  s.geminiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	prompt := buildPrompt(rawText, today, categories, accounts)

	modelsToTry := []string{s.geminiModel}
	for _, fallback := range []string{"gemini-3.5-flash", "gemini-2.5-flash", "gemini-3.1-flash-lite"} {
		if fallback != s.geminiModel {
			modelsToTry = append(modelsToTry, fallback)
		}
	}

	var result *genai.GenerateContentResponse
	var genErr error

	for i, modelName := range modelsToTry {
		result, genErr = client.Models.GenerateContent(
			ctx,
			modelName,
			genai.Text(prompt),
			nil,
		)
		if genErr == nil {
			break
		}

		log.Printf("Gemini call failed with model %s (attempt %d/%d): %v", modelName, i+1, len(modelsToTry), genErr)
		if i < len(modelsToTry)-1 {
			time.Sleep(1 * time.Second)
		}
	}

	if genErr != nil {
		return nil, fmt.Errorf("gemini request failed (tried %d models): %w", len(modelsToTry), genErr)
	}

	rawJSON := extractJSON(result.Text())
	if rawJSON == "" {
		return nil, fmt.Errorf("gemini returned empty or non-JSON response")
	}

	var parsed model.GeminiParsed
	if err := json.Unmarshal([]byte(rawJSON), &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse gemini response as JSON: %w\nRaw: %s", err, rawJSON)
	}

	return &parsed, nil
}

func buildPrompt(rawText, today string, categories, accounts []model.CatalogItem) string {
	if strings.TrimSpace(today) == "" {
		today = time.Now().Format("2006-01-02")
	}
	catsJSON, _ := json.Marshal(categories)
	accsJSON, _ := json.Marshal(accounts)
	if string(catsJSON) == "null" {
		catsJSON = []byte("[]")
	}
	if string(accsJSON) == "null" {
		accsJSON = []byte("[]")
	}

	return fmt.Sprintf(`You are a financial transaction parser. Given a user's spoken or typed note, extract structured transaction data.

Today (user's local calendar date) is: %s

The user has these categories. You MUST pick category_id from this list when possible. Never invent an id:
%s

The user has these accounts. You MUST pick account_id from this list when possible, or null if unknown. Never invent an id:
%s

Respond ONLY with a valid JSON object in exactly this format (no markdown, no explanation):
{
  "type": "EXPENSE or INCOME",
  "amount": <positive number, or 0 if unknown>,
  "category_id": "<id from the category list, or empty string>",
  "category_name": "<name from the category list>",
  "account_id": "<id from the account list, or null>",
  "account_name": "<name from the account list, or empty string>",
  "date": "<yyyy-MM-dd, never empty, never in the future>"
}

Rules:
- type must be EXPENSE or INCOME using those exact uppercase values.
- Resolve relative dates like "today" and "yesterday" against the user's today (%s). Do not return a future date.
- amount must be a positive number when a spend/credit amount is present.
- If amount or type cannot be determined, still return JSON with amount 0 and/or empty type.
- Pick the closest real category (e.g. Uber -> Travel) using the provided list. Do not make up category names.
- Pick a real account if the text names one (e.g. HDFC). Otherwise account_id should be null.

Raw text to parse:
"%s"`, today, string(catsJSON), string(accsJSON), today, rawText)
}

func extractJSON(text string) string {
	text = strings.TrimSpace(text)
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
