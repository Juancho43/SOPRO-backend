package usecases

import (
	"context"
	"errors"
	"math"
	"strings"
	"time"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/repositories"
)

// ErrBudgetItemNotFound is returned when a budget item does not exist or does not belong to the user.
var ErrBudgetItemNotFound = errors.New("partida presupuestaria no encontrada")

// CreateBudgetItemCommand represents the input payload to register an income or expense budget item (HU 4.2).
type CreateBudgetItemCommand struct {
	UserID               string  `json:"user_id"`
	Year                 int     `json:"year"`
	Month                int     `json:"month"`
	Type                 string  `json:"type"`
	ExpenseCategory      *string `json:"expense_category,omitempty"`
	ExpenseCategoryCamel *string `json:"expenseCategory,omitempty" swaggerignore:"true"`
	Concept              string  `json:"concept"`
	Amount               float64 `json:"amount"`
	Currency             string  `json:"currency"`
}

// CreateBudgetItemResponse represents the output returned after creating a budget item.
type CreateBudgetItemResponse struct {
	*entities.BudgetItem
}

// CreateBudgetItemUseCase handles the creation and validation of a monthly income or expense item.
type CreateBudgetItemUseCase struct {
	repo     repositories.FinanceRepository
	timezone *time.Location
}

// NewCreateBudgetItemUseCase constructs a new CreateBudgetItemUseCase instance.
func NewCreateBudgetItemUseCase(repo repositories.FinanceRepository, tz string) (*CreateBudgetItemUseCase, error) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		loc = time.Local
	}
	return &CreateBudgetItemUseCase{
		repo:     repo,
		timezone: loc,
	}, nil
}

// Execute validates the command parameters and persists a new budget item in the repository.
func (uc *CreateBudgetItemUseCase) Execute(ctx context.Context, cmd CreateBudgetItemCommand) (*CreateBudgetItemResponse, error) {
	year, month, err := normalizeFinancePeriod(cmd.Year, cmd.Month, uc.timezone)
	if err != nil {
		return nil, err
	}

	concept := strings.TrimSpace(cmd.Concept)
	if concept == "" {
		return nil, errors.New("el concepto de la partida presupuestaria no puede estar vacío")
	}
	if len(concept) > 255 {
		return nil, errors.New("el concepto no puede exceder 255 caracteres")
	}

	if cmd.Amount <= 0 {
		return nil, errors.New("el monto de la partida presupuestaria debe ser mayor a cero")
	}

	itemType, err := parseBudgetItemType(cmd.Type)
	if err != nil {
		return nil, err
	}

	rawCategory := normalizeOptionalString(cmd.ExpenseCategory, cmd.ExpenseCategoryCamel)
	category, err := parseExpenseCategory(itemType, rawCategory)
	if err != nil {
		return nil, err
	}

	currency := strings.ToUpper(strings.TrimSpace(cmd.Currency))
	if currency == "" {
		currency = "USD"
	}

	item := &entities.BudgetItem{
		UserID:          cmd.UserID,
		Year:            year,
		Month:           month,
		Type:            itemType,
		ExpenseCategory: category,
		Concept:         concept,
		Amount:          roundMoney(cmd.Amount),
		Currency:        currency,
	}

	if err := uc.repo.CreateBudgetItem(ctx, item); err != nil {
		return nil, err
	}

	return &CreateBudgetItemResponse{BudgetItem: item}, nil
}

// parseBudgetItemType validates and normalizes the budget item type (INCOME or EXPENSE).
func parseBudgetItemType(raw string) (entities.BudgetItemType, error) {
	switch strings.ToUpper(strings.TrimSpace(raw)) {
	case "INCOME", "INGRESO":
		return entities.BudgetItemTypeIncome, nil
	case "EXPENSE", "GASTO":
		return entities.BudgetItemTypeExpense, nil
	default:
		return "", errors.New("tipo de partida inválido (valores permitidos: INCOME, EXPENSE)")
	}
}

// parseExpenseCategory validates that EXPENSE items specify FIXED or VARIABLE and INCOME items have no expense category.
func parseExpenseCategory(itemType entities.BudgetItemType, rawCategory *string) (*entities.ExpenseCategory, error) {
	if itemType == entities.BudgetItemTypeIncome {
		return nil, nil
	}
	if rawCategory == nil || strings.TrimSpace(*rawCategory) == "" {
		return nil, errors.New("la categoría de gasto (expense_category: FIXED o VARIABLE) es obligatoria para partidas de tipo EXPENSE")
	}

	switch strings.ToUpper(strings.TrimSpace(*rawCategory)) {
	case "FIXED", "FIJO":
		cat := entities.ExpenseCategoryFixed
		return &cat, nil
	case "VARIABLE":
		cat := entities.ExpenseCategoryVariable
		return &cat, nil
	default:
		return nil, errors.New("categoría de gasto inválida (valores permitidos: FIXED, VARIABLE)")
	}
}

// roundMoney rounds a monetary float value to 2 decimal places.
func roundMoney(val float64) float64 {
	return math.Round(val*100) / 100
}
