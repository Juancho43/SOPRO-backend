package usecases_test

import (
	"context"
	"testing"

	usecases "github.com/Juancho43/SOPRO-backend/pkg/internal/use_cases"
)

// TestGetMonthlyBudgetUseCase verifies break-even calculation, net balance, and solvency status.
func TestGetMonthlyBudgetUseCase(t *testing.T) {
	repo := newMockFinanceRepo()
	createUC, _ := usecases.NewCreateBudgetItemUseCase(repo, "America/Buenos_Aires")
	getUC, _ := usecases.NewGetMonthlyBudgetUseCase(repo, "America/Buenos_Aires")
	ctx := context.Background()

	fixedCat := "FIXED"
	varCat := "VARIABLE"

	// Projected Income: 3000 USD
	_, _ = createUC.Execute(ctx, usecases.CreateBudgetItemCommand{
		UserID:  "u1",
		Year:    2026,
		Month:   11,
		Type:    "INCOME",
		Concept: "Primary Salary",
		Amount:  3000,
	})

	// Fixed Expenses: Rent (1200) + Servers (300) = 1500 (Break-even)
	_, _ = createUC.Execute(ctx, usecases.CreateBudgetItemCommand{
		UserID:          "u1",
		Year:            2026,
		Month:           11,
		Type:            "EXPENSE",
		ExpenseCategory: &fixedCat,
		Concept:         "Rent",
		Amount:          1200,
	})
	_, _ = createUC.Execute(ctx, usecases.CreateBudgetItemCommand{
		UserID:          "u1",
		Year:            2026,
		Month:           11,
		Type:            "EXPENSE",
		ExpenseCategory: &fixedCat,
		Concept:         "Servers",
		Amount:          300,
	})

	// Variable Expenses: Lifestyle (500) => Total Expenses = 2000
	_, _ = createUC.Execute(ctx, usecases.CreateBudgetItemCommand{
		UserID:          "u1",
		Year:            2026,
		Month:           11,
		Type:            "EXPENSE",
		ExpenseCategory: &varCat,
		Concept:         "Lifestyle",
		Amount:          500,
	})

	t.Run("Computes break-even, net balance, and solvent status (Green)", func(t *testing.T) {
		resp, err := getUC.Execute(ctx, usecases.GetMonthlyBudgetQuery{
			UserID: "u1",
			Year:   2026,
			Month:  11,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.TotalIncome != 3000 {
			t.Errorf("expected TotalIncome 3000, got %.2f", resp.TotalIncome)
		}
		if resp.TotalFixedExpenses != 1500 {
			t.Errorf("expected break-even (TotalFixedExpenses) 1500, got %.2f", resp.TotalFixedExpenses)
		}
		if resp.TotalVariableExpenses != 500 {
			t.Errorf("expected TotalVariableExpenses 500, got %.2f", resp.TotalVariableExpenses)
		}
		if resp.TotalExpenses != 2000 {
			t.Errorf("expected TotalExpenses 2000, got %.2f", resp.TotalExpenses)
		}
		if resp.NetBalance != 1000 || !resp.IsSolvent {
			t.Errorf("expected NetBalance 1000 and IsSolvent=true, got %.2f / %v", resp.NetBalance, resp.IsSolvent)
		}
	})

	t.Run("Detects financial deficit (Red) when total expenses exceed income", func(t *testing.T) {
		_, _ = createUC.Execute(ctx, usecases.CreateBudgetItemCommand{
			UserID:          "u1",
			Year:            2026,
			Month:           11,
			Type:            "EXPENSE",
			ExpenseCategory: &varCat,
			Concept:         "Extraordinary Purchase",
			Amount:          1500,
		})

		resp, err := getUC.Execute(ctx, usecases.GetMonthlyBudgetQuery{
			UserID: "u1",
			Year:   2026,
			Month:  11,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.IsSolvent {
			t.Errorf("expected IsSolvent=false when expenses (3500) exceed income (3000)")
		}
		if resp.NetBalance != -500 || resp.DeficitGap != 500 {
			t.Errorf("expected NetBalance -500 and DeficitGap 500, got %.2f / %.2f", resp.NetBalance, resp.DeficitGap)
		}
	})
}
