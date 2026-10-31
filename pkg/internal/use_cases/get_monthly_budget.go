package usecases

import (
	"context"
	"time"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/repositories"
)

// GetMonthlyBudgetQuery represents the query parameters to retrieve the monthly budget items and break-even summary (HU 4.2).
type GetMonthlyBudgetQuery struct {
	UserID string `json:"user_id"`
	Year   int    `json:"year"`
	Month  int    `json:"month"`
}

// GetMonthlyBudgetResponse represents the monthly budget items along with the calculated financial summary.
type GetMonthlyBudgetResponse struct {
	entities.BudgetOverview
}

// GetMonthlyBudgetUseCase retrieves a user's budget items for a given month and computes break-even and solvency metrics.
type GetMonthlyBudgetUseCase struct {
	repo     repositories.FinanceRepository
	timezone *time.Location
}

// NewGetMonthlyBudgetUseCase constructs a new GetMonthlyBudgetUseCase instance.
func NewGetMonthlyBudgetUseCase(repo repositories.FinanceRepository, tz string) (*GetMonthlyBudgetUseCase, error) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		loc = time.Local
	}
	return &GetMonthlyBudgetUseCase{
		repo:     repo,
		timezone: loc,
	}, nil
}

// Execute fetches the budget items for the requested month and calculates the financial break-even summary.
func (uc *GetMonthlyBudgetUseCase) Execute(ctx context.Context, query GetMonthlyBudgetQuery) (*GetMonthlyBudgetResponse, error) {
	year, month, err := normalizeFinancePeriod(query.Year, query.Month, uc.timezone)
	if err != nil {
		return nil, err
	}

	items, err := uc.repo.ListBudgetItemsByPeriod(ctx, query.UserID, year, month)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []*entities.BudgetItem{}
	}

	summary := CalculateBudgetSummary(year, month, items)

	return &GetMonthlyBudgetResponse{
		BudgetOverview: entities.BudgetOverview{
			BudgetSummary: summary,
			Items:         items,
			Summary:       summary,
		},
	}, nil
}

// CalculateBudgetSummary computes total income, fixed expenses (break-even), variable expenses, net balance, and solvency status.
func CalculateBudgetSummary(year int, month int, items []*entities.BudgetItem) entities.BudgetSummary {
	var totalIncome, totalFixed, totalVariable float64

	for _, item := range items {
		if item == nil {
			continue
		}
		switch item.Type {
		case entities.BudgetItemTypeIncome:
			totalIncome += item.Amount
		case entities.BudgetItemTypeExpense:
			if item.ExpenseCategory != nil && *item.ExpenseCategory == entities.ExpenseCategoryFixed {
				totalFixed += item.Amount
			} else {
				totalVariable += item.Amount
			}
		}
	}

	totalIncome = roundMoney(totalIncome)
	totalFixed = roundMoney(totalFixed)
	totalVariable = roundMoney(totalVariable)
	totalExpenses := roundMoney(totalFixed + totalVariable)
	netBalance := roundMoney(totalIncome - totalExpenses)
	isSolvent := netBalance >= 0

	deficitGap := 0.0
	if netBalance < 0 {
		deficitGap = roundMoney(-netBalance)
	}

	breakEvenGap := 0.0
	if totalIncome < totalFixed {
		breakEvenGap = roundMoney(totalFixed - totalIncome)
	}

	return entities.BudgetSummary{
		Year:                  year,
		Month:                 month,
		TotalIncome:           totalIncome,
		TotalFixedExpenses:    totalFixed,
		TotalVariableExpenses: totalVariable,
		TotalExpenses:         totalExpenses,
		NetBalance:            netBalance,
		IsSolvent:             isSolvent,
		DeficitGap:            deficitGap,
		BreakEvenGap:          breakEvenGap,
	}
}
