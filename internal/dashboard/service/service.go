package service

import (
	"context"
	"time"

	"expent-backend/internal/dashboard/model"
	"expent-backend/internal/dashboard/repository"

	"go.uber.org/zap"
)

type Service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

// GetDashboard aggregates all dashboard data for the given user.
func (s *Service) GetDashboard(ctx context.Context, userID string) (*model.DashboardResponse, error) {
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	monthEnd := monthStart.AddDate(0, 1, 0).Add(-time.Nanosecond)

	// 1. Monthly transactions (with category relation fetched)
	monthlyTxs, err := s.repo.GetMonthlyTransactions(userID, monthStart, monthEnd)
	if err != nil {
		zap.L().Error("failed to get monthly transactions", zap.Error(err))
		return nil, err
	}

	var totalIncome, totalExpense float64
	categoryTotals := make(map[string]*model.CategoryBreakdown)

	for _, t := range monthlyTxs {
		switch t.Type {
		case "income":
			totalIncome += t.Amount
		case "expense":
			totalExpense += t.Amount
			cat := t.RelationsTransaction.Category
			if cat != nil {
				if cb, ok := categoryTotals[cat.ID]; ok {
					cb.TotalAmount += t.Amount
				} else {
					categoryTotals[cat.ID] = &model.CategoryBreakdown{
						CategoryID:   cat.ID,
						CategoryName: cat.Name,
						TotalAmount:  t.Amount,
					}
				}
			}
		}
	}

	var categoryBreakdown []model.CategoryBreakdown
	for _, cb := range categoryTotals {
		categoryBreakdown = append(categoryBreakdown, *cb)
	}

	// 2. Recent transactions (last 5)
	recentTxs, err := s.repo.GetRecentTransactions(userID, 5)
	if err != nil {
		zap.L().Error("failed to get recent transactions", zap.Error(err))
		return nil, err
	}
	recentTransactions := repository.MapRecentTransactions(recentTxs)

	// 3. Budget utilization
	budgets, err := s.repo.GetBudgets(userID)
	if err != nil {
		zap.L().Error("failed to get budgets", zap.Error(err))
		return nil, err
	}

	var budgetUtilization []model.BudgetUtilization
	for _, b := range budgets {
		// Get expense transactions for this budget's category and date range
		budgetTxs, err := s.repo.GetTransactionsByCategoryAndDateRange(
			userID, b.CategoryID, b.StartDate, b.EndDate,
		)
		if err != nil {
			zap.L().Error("failed to get budget transactions",
				zap.String("budgetId", b.ID),
				zap.Error(err),
			)
			continue
		}

		var spent float64
		for _, t := range budgetTxs {
			spent += t.Amount
		}

		categoryName := ""
		cat := b.RelationsBudget.Category
		if cat != nil {
			categoryName = cat.Name
		}

		budgetUtilization = append(budgetUtilization, model.BudgetUtilization{
			BudgetID:     b.ID,
			CategoryID:   b.CategoryID,
			CategoryName: categoryName,
			PeriodType:   b.PeriodType,
			LimitAmount:  b.LimitAmount,
			SpentAmount:  spent,
			StartDate:    b.StartDate.Format("2006-01-02"),
			EndDate:      b.EndDate.Format("2006-01-02"),
		})
	}

	return &model.DashboardResponse{
		TotalIncome:        totalIncome,
		TotalExpense:       totalExpense,
		Balance:            totalIncome - totalExpense,
		RecentTransactions: recentTransactions,
		CategoryBreakdown:  categoryBreakdown,
		BudgetUtilization:  budgetUtilization,
	}, nil
}
