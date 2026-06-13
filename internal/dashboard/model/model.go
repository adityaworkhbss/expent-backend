package model

import "time"

// DashboardResponse is the top-level response returned by GET /dashboard.
type DashboardResponse struct {
	TotalIncome        float64              `json:"totalIncome"`
	TotalExpense       float64              `json:"totalExpense"`
	Balance            float64              `json:"balance"`
	RecentTransactions []RecentTransaction  `json:"recentTransactions"`
	CategoryBreakdown  []CategoryBreakdown  `json:"categoryBreakdown"`
	BudgetUtilization  []BudgetUtilization  `json:"budgetUtilization"`
}

// RecentTransaction represents a single recent transaction entry.
type RecentTransaction struct {
	ID         string    `json:"id"`
	Amount     float64   `json:"amount"`
	Type       string    `json:"type"`
	Notes      string    `json:"notes,omitempty"`
	Date       string    `json:"date"`
	CategoryID string    `json:"categoryId"`
	AccountID  string    `json:"accountId"`
	CreatedAt  time.Time `json:"createdAt"`
}

// CategoryBreakdown represents expense totals grouped by category.
type CategoryBreakdown struct {
	CategoryID   string  `json:"categoryId"`
	CategoryName string  `json:"categoryName"`
	TotalAmount  float64 `json:"totalAmount"`
}

// BudgetUtilization represents how much of a budget has been spent.
type BudgetUtilization struct {
	BudgetID     string  `json:"budgetId"`
	CategoryID   string  `json:"categoryId"`
	CategoryName string  `json:"categoryName"`
	PeriodType   string  `json:"periodType"`
	LimitAmount  float64 `json:"limitAmount"`
	SpentAmount  float64 `json:"spentAmount"`
	StartDate    string  `json:"startDate"`
	EndDate      string  `json:"endDate"`
}
