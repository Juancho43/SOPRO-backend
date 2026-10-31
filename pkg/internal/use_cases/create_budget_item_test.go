package usecases_test

import (
	"context"
	"testing"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	usecases "github.com/Juancho43/SOPRO-backend/pkg/internal/use_cases"
)

// TestCreateBudgetItemUseCase verifies creation and domain validation of income and expense budget items.
func TestCreateBudgetItemUseCase(t *testing.T) {
	repo := newMockFinanceRepo()
	uc, err := usecases.NewCreateBudgetItemUseCase(repo, "America/Buenos_Aires")
	if err != nil {
		t.Fatalf("unexpected error constructing use case: %v", err)
	}

	t.Run("Creates an INCOME budget item with default currency USD", func(t *testing.T) {
		resp, err := uc.Execute(context.Background(), usecases.CreateBudgetItemCommand{
			UserID:  "u1",
			Year:    2026,
			Month:   11,
			Type:    "INCOME",
			Concept: "Monthly Retainer",
			Amount:  3500.556,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.Type != entities.BudgetItemTypeIncome {
			t.Errorf("expected type INCOME, got %s", resp.Type)
		}
		if resp.Amount != 3500.56 {
			t.Errorf("expected rounded amount 3500.56, got %.2f", resp.Amount)
		}
		if resp.Currency != "USD" {
			t.Errorf("expected default currency USD, got %s", resp.Currency)
		}
		if resp.ExpenseCategory != nil {
			t.Errorf("expected nil ExpenseCategory for INCOME, got %v", resp.ExpenseCategory)
		}
	})

	t.Run("Requires FIXED or VARIABLE expense_category when type is EXPENSE", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), usecases.CreateBudgetItemCommand{
			UserID:  "u1",
			Year:    2026,
			Month:   11,
			Type:    "EXPENSE",
			Concept: "Cloud Infrastructure",
			Amount:  150,
		})
		if err == nil {
			t.Errorf("expected error when creating EXPENSE without expense_category")
		}

		fixedCat := "FIXED"
		resp, err := uc.Execute(context.Background(), usecases.CreateBudgetItemCommand{
			UserID:          "u1",
			Year:            2026,
			Month:           11,
			Type:            "EXPENSE",
			ExpenseCategory: &fixedCat,
			Concept:         "Cloud Infrastructure",
			Amount:          150,
			Currency:        "EUR",
		})
		if err != nil {
			t.Fatalf("unexpected error creating EXPENSE with FIXED category: %v", err)
		}
		if resp.ExpenseCategory == nil || *resp.ExpenseCategory != entities.ExpenseCategoryFixed {
			t.Errorf("expected FIXED expense category, got %v", resp.ExpenseCategory)
		}
	})
}
