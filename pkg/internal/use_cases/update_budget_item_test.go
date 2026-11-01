package usecases_test

import (
	"context"
	"testing"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	usecases "github.com/Juancho43/SOPRO-backend/pkg/internal/use_cases"
)

// TestUpdateBudgetItemUseCase verifies updating concept, amount, and category of a budget item.
func TestUpdateBudgetItemUseCase(t *testing.T) {
	repo := newMockFinanceRepo()
	fixedCat := entities.ExpenseCategoryFixed
	repo.budgetItems["item-1"] = &entities.BudgetItem{
		ID:              "item-1",
		UserID:          "u1",
		Year:            2026,
		Month:           11,
		Type:            entities.BudgetItemTypeExpense,
		ExpenseCategory: &fixedCat,
		Concept:         "Office Rent",
		Amount:          1000,
		Currency:        "USD",
	}

	uc := usecases.NewUpdateBudgetItemUseCase(repo)

	t.Run("Updates concept and amount of an existing budget item", func(t *testing.T) {
		newConcept := "Coworking Membership"
		newAmount := 450.0
		varCat := "VARIABLE"

		resp, err := uc.Execute(context.Background(), usecases.UpdateBudgetItemCommand{
			ID:              "item-1",
			UserID:          "u1",
			Concept:         &newConcept,
			Amount:          &newAmount,
			ExpenseCategory: &varCat,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.Concept != "Coworking Membership" || resp.Amount != 450 {
			t.Errorf("unexpected updated item: %+v", resp.BudgetItem)
		}
		if resp.ExpenseCategory == nil || *resp.ExpenseCategory != entities.ExpenseCategoryVariable {
			t.Errorf("expected updated category VARIABLE, got %v", resp.ExpenseCategory)
		}
	})

	t.Run("Returns ErrBudgetItemNotFound when item does not exist", func(t *testing.T) {
		newAmount := 200.0
		_, err := uc.Execute(context.Background(), usecases.UpdateBudgetItemCommand{
			ID:     "non-existent",
			UserID: "u1",
			Amount: &newAmount,
		})
		if err != usecases.ErrBudgetItemNotFound {
			t.Errorf("expected ErrBudgetItemNotFound, got %v", err)
		}
	})
}
