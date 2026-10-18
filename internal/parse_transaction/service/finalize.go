package service

import (
	"strings"
	"time"
	"unicode"

	"expent-backend/internal/parse_transaction/model"
)

var nameAliases = map[string][]string{
	"travel":        {"transport", "transportation", "commute", "uber", "ola", "cab", "taxi", "flight", "fuel"},
	"food":          {"dining", "restaurant", "groceries", "grocery", "eats", "food & dining", "food and dining"},
	"shopping":      {"retail", "amazon", "flipkart", "clothes"},
	"entertainment": {"movies", "netflix", "spotify", "fun"},
	"utilities":     {"bills", "electricity", "internet", "wifi", "recharge"},
	"salary":        {"payroll", "income", "paycheck"},
}

func finalizeParse(raw model.GeminiParsed, categories, accounts []model.CatalogItem, today string) model.ParseResult {
	todayDate := parseToday(today)
	txnType := normalizeTxnType(raw.Type)

	catName := firstNonEmpty(raw.CategoryName, raw.Category)
	cat := matchCatalog(raw.CategoryID, catName, txnType, categories)
	if cat != nil {
		if txnType == "" {
			txnType = normalizeTxnType(cat.Type)
		}
	}

	accName := raw.AccountName
	var accID *string
	acc := matchCatalog(deref(raw.AccountID), accName, "", accounts)
	if acc != nil {
		id := acc.ID
		accID = &id
		accName = acc.Name
	} else {
		accID = nil
		if strings.TrimSpace(accName) == "" {
			accName = ""
		}
	}

	date := normalizeDate(raw.Date, todayDate)
	currency := strings.ToUpper(strings.TrimSpace(raw.Currency))
	if currency == "" {
		currency = "INR"
	}

	merchant := firstNonEmpty(raw.Merchant, guessMerchant(raw.Description, raw.Note))
	description := firstNonEmpty(raw.Description, merchant, catName)
	note := firstNonEmpty(raw.Note, description)

	out := model.ParsedTransaction{
		Type:          txnType,
		Amount:        raw.Amount.Float64(),
		AccountID:     accID,
		AccountName:   accName,
		PaymentMethod: strings.TrimSpace(raw.PaymentMethod),
		Description:   description,
		Merchant:      merchant,
		Note:          note,
		Date:          date,
		Currency:      currency,
		Confidence:    raw.Confidence.Float64(),
	}
	if cat != nil {
		out.CategoryID = cat.ID
		out.CategoryName = cat.Name
	} else {
		out.CategoryName = catName
	}

	var missing []string
	if out.Amount <= 0 {
		missing = append(missing, "amount")
	}
	if out.Type != "INCOME" && out.Type != "EXPENSE" {
		missing = append(missing, "type")
		out.Type = ""
	}

	if missing == nil {
		missing = []string{}
	}

	return model.ParseResult{
		RequiresUserInput:     len(missing) > 0,
		MissingRequiredFields: missing,
		Data:                  out,
	}
}

func parseToday(today string) time.Time {
	if t, err := time.Parse("2006-01-02", strings.TrimSpace(today)); err == nil {
		return t
	}
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
}

func normalizeDate(raw string, today time.Time) string {
	raw = strings.TrimSpace(raw)
	parsed, ok := parseFlexibleDate(raw, today)
	if !ok {
		return today.Format("2006-01-02")
	}
	if parsed.After(today) {
		return today.Format("2006-01-02")
	}
	return parsed.Format("2006-01-02")
}

func parseFlexibleDate(raw string, today time.Time) (time.Time, bool) {
	if raw == "" {
		return time.Time{}, false
	}
	lower := strings.ToLower(raw)
	switch lower {
	case "today":
		return today, true
	case "yesterday":
		return today.AddDate(0, 0, -1), true
	}
	layouts := []string{
		"2006-01-02",
		time.RFC3339,
		"02-01-2006",
		"02/01/2006",
		"02-Jan-06",
		"02-Jan-2006",
		"2 Jan 2006",
		"Jan 2, 2006",
		"2006/01/02",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, raw); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

func normalizeTxnType(t string) string {
	s := strings.ToUpper(strings.TrimSpace(t))
	switch s {
	case "INCOME", "CREDIT", "CREDITED", "SALARY":
		return "INCOME"
	case "EXPENSE", "DEBIT", "DEBITED", "SPENT":
		return "EXPENSE"
	default:
		return ""
	}
}

func matchCatalog(id, name, txnType string, items []model.CatalogItem) *model.CatalogItem {
	id = strings.TrimSpace(id)
	if id != "" {
		for i := range items {
			if items[i].ID == id {
				return &items[i]
			}
		}
	}

	want := normalizeTxnType(txnType)
	normName := normalizeName(name)
	if normName == "" {
		return nil
	}

	scoreBest := func(requireType bool) *model.CatalogItem {
		var best *model.CatalogItem
		bestScore := 0
		for i := range items {
			item := &items[i]
			if requireType && want != "" && normalizeTxnType(item.Type) != want {
				continue
			}
			score := nameScore(normName, item.Name)
			if score > bestScore {
				bestScore = score
				best = item
			}
		}
		if bestScore >= 2 {
			return best
		}
		return nil
	}

	if found := scoreBest(true); found != nil {
		return found
	}
	return scoreBest(false)
}

func nameScore(want, have string) int {
	h := normalizeName(have)
	if h == "" {
		return 0
	}
	if want == h {
		return 5
	}
	if strings.Contains(h, want) || strings.Contains(want, h) {
		return 4
	}
	if aliasesOverlap(want, h) {
		return 3
	}
	return 0
}

func aliasesOverlap(a, b string) bool {
	for key, list := range nameAliases {
		if relates(a, key, list) && relates(b, key, list) {
			return true
		}
	}
	return false
}

func relates(name, key string, list []string) bool {
	if name == key || strings.Contains(name, key) || strings.Contains(key, name) {
		return true
	}
	for _, alias := range list {
		an := normalizeName(alias)
		if name == an || strings.Contains(name, an) || strings.Contains(an, name) {
			return true
		}
	}
	return false
}

func normalizeName(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, "&", " and ")
	var b strings.Builder
	prevSpace := false
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			prevSpace = false
			continue
		}
		if !prevSpace {
			b.WriteByte(' ')
			prevSpace = true
		}
	}
	return strings.Join(strings.Fields(b.String()), " ")
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if s := strings.TrimSpace(v); s != "" {
			return s
		}
	}
	return ""
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func guessMerchant(desc, note string) string {
	return firstNonEmpty(desc, note)
}
