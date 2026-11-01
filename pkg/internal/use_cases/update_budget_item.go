package usecases

import (
	"context"
	"errors"
	"strings"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/repositories"
)

// UpdateBudgetItemCommand represents the input payload to modify an existing budget item (HU 4.2).
type UpdateBudgetItemCommand struct {
	ID                   string   `json:"id"`
	UserID               string   `json:"user_id"`
	Type                 *string  `json:"type,omitempty"`
	ExpenseCategory      *string  `json:"expense_category,omitempty"`
	ExpenseCategoryCamel *string  `json:"expenseCategory,omitempty" swaggerignore:"true"`
	Concept              *string  `json:"concept,omitempty"`
	Amount               *float64 `json:"amount,omitempty"`
	Currency             *string  `json:"currency,omitempty"`
}

// UpdateBudgetItemResponse represents the updated budget item returned by the use case.
type UpdateBudgetItemResponse struct {
	*entities.BudgetItem
}

// UpdateBudgetItemUseCase handles updating the fields of an existing monthly budget item.
type UpdateBudgetItemUseCase struct {
	repo repositories.FinanceRepository
}

// NewUpdateBudgetItemUseCase constructs a new UpdateBudgetItemUseCase instance.
func NewUpdateBudgetItemUseCase(repo repositories.FinanceRepository) *UpdateBudgetItemUseCase {
	return &UpdateBudgetItemUseCase{repo: repo}
}

// Execute validates the updated fields, persists the changes, and returns the updated budget item.
func (uc *UpdateBudgetItemUseCase) Execute(ctx context.Context, cmd UpdateBudgetItemCommand) (*UpdateBudgetItemResponse, error) {
	existing, err := uc.repo.GetBudgetItemByID(ctx, cmd.ID, cmd.UserID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrBudgetItemNotFound
	}

	if cmd.Concept != nil {
		concept := strings.TrimSpace(*cmd.Concept)
		if concept == "" {
			return nil, errors.New("el concepto de la partida presupuestaria no puede estar vacío")
		}
		if len(concept) > 255 {
			return nil, errors.New("el concepto no puede exceder 255 caracteres")
		}
		existing.Concept = concept
	}

	if cmd.Amount != nil {
		if *cmd.Amount <= 0 {
			return nil, errors.New("el monto de la partida presupuestaria debe ser mayor a cero")
		}
		existing.Amount = roundMoney(*cmd.Amount)
	}

	if cmd.Currency != nil {
		curr := strings.ToUpper(strings.TrimSpace(*cmd.Currency))
		if curr != "" {
			existing.Currency = curr
		}
	}

	if cmd.Type != nil {
		itemType, err := parseBudgetItemType(*cmd.Type)
		if err != nil {
			return nil, err
		}
		existing.Type = itemType
	}

	rawCat := normalizeOptionalString(cmd.ExpenseCategory, cmd.ExpenseCategoryCamel)
	if cmd.Type != nil || rawCat != nil {
		if rawCat == nil && existing.ExpenseCategory != nil {
			existingCat := string(*existing.ExpenseCategory)
			rawCat = &existingCat
		}
		cat, err := parseExpenseCategory(existing.Type, rawCat)
		if err != nil {
			return nil, err
		}
		existing.ExpenseCategory = cat
	}

	if err := uc.repo.UpdateBudgetItem(ctx, existing); err != nil {
		return nil, err
	}

	return &UpdateBudgetItemResponse{BudgetItem: existing}, nil
}
