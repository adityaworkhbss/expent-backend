package service

import (
	"testing"

	"expent-backend/internal/parse_transaction/model"
	"expent-backend/internal/shared"
)

func TestFinalizeParseMapsRealIDs(t *testing.T) {
	cats := []model.CatalogItem{
		{ID: "cat-travel", Name: "Travel", Type: "expense"},
		{ID: "cat-food", Name: "Food", Type: "expense"},
	}
	accs := []model.CatalogItem{
		{ID: "acc-hdfc", Name: "HDFC", Type: "BANK"},
	}

	got := finalizeParse(model.GeminiParsed{
		Type:         "EXPENSE",
		Amount:       shared.FlexFloat64(250),
		CategoryName: "Transport",
		AccountName:  "HDFC",
		Date:         "today",
	}, cats, accs, "2026-10-18")

	if got.RequiresUserInput {
		t.Fatalf("unexpected requires_user_input: %+v", got.MissingRequiredFields)
	}
	if got.Data.Type != "EXPENSE" {
		t.Fatalf("type %s", got.Data.Type)
	}
	if got.Data.CategoryID != "cat-travel" || got.Data.CategoryName != "Travel" {
		t.Fatalf("category %+v", got.Data)
	}
	if got.Data.AccountID == nil || *got.Data.AccountID != "acc-hdfc" {
		t.Fatalf("account_id %+v", got.Data.AccountID)
	}
	if got.Data.AccountName != "HDFC" {
		t.Fatalf("account_name %s", got.Data.AccountName)
	}
	if got.Data.Date != "2026-10-18" {
		t.Fatalf("date %s", got.Data.Date)
	}
}

func TestFinalizeParseMissingAmountAndType(t *testing.T) {
	got := finalizeParse(model.GeminiParsed{
		CategoryName: "something happened",
	}, nil, nil, "2026-10-18")

	if !got.RequiresUserInput {
		t.Fatal("expected requires_user_input")
	}
	if len(got.MissingRequiredFields) != 2 {
		t.Fatalf("missing %v", got.MissingRequiredFields)
	}
	if got.Data.Date != "2026-10-18" {
		t.Fatalf("date should default to today, got %s", got.Data.Date)
	}
}

func TestFinalizeParseClampsFutureDate(t *testing.T) {
	got := finalizeParse(model.GeminiParsed{
		Type:   "expense",
		Amount: shared.FlexFloat64(10),
		Date:   "2026-12-01",
	}, nil, nil, "2026-10-18")
	if got.Data.Date != "2026-10-18" {
		t.Fatalf("future date not clamped: %s", got.Data.Date)
	}
	if got.Data.Type != "EXPENSE" {
		t.Fatalf("type %s", got.Data.Type)
	}
}

func TestFinalizeParseUnknownAccountIsNull(t *testing.T) {
	got := finalizeParse(model.GeminiParsed{
		Type:        "EXPENSE",
		Amount:      shared.FlexFloat64(10),
		AccountName: "Mystery Bank",
	}, nil, []model.CatalogItem{{ID: "acc-1", Name: "HDFC", Type: "BANK"}}, "2026-10-18")
	if got.Data.AccountID != nil {
		t.Fatalf("expected null account_id, got %v", *got.Data.AccountID)
	}
}
