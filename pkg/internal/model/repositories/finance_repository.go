package repositories

import (
	"context"
	"time"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
)

// FinanceRepository defines the persistence port for Monthly Goals (HU 4.1) and Monthly Budget items (HU 4.2).
type FinanceRepository interface {
	// Monthly Goals (HU 4.1)
	CreateMonthlyGoal(ctx context.Context, goal *entities.MonthlyGoal) error
	GetMonthlyGoalByID(ctx context.Context, id string, userID string) (*entities.MonthlyGoal, error)
	ListMonthlyGoalsByPeriod(ctx context.Context, userID string, year int, month int) ([]*entities.MonthlyGoal, error)
	UpdateMonthlyGoal(ctx context.Context, goal *entities.MonthlyGoal) error
	UpdateMonthlyGoalStatus(ctx context.Context, id string, userID string, status entities.MonthlyGoalStatus, achievedAt *time.Time) error
	DeleteMonthlyGoal(ctx context.Context, id string, userID string) error

	// Monthly Budget & Break-Even (HU 4.2)
	CreateBudgetItem(ctx context.Context, item *entities.BudgetItem) error
	GetBudgetItemByID(ctx context.Context, id string, userID string) (*entities.BudgetItem, error)
	ListBudgetItemsByPeriod(ctx context.Context, userID string, year int, month int) ([]*entities.BudgetItem, error)
	UpdateBudgetItem(ctx context.Context, item *entities.BudgetItem) error
	DeleteBudgetItem(ctx context.Context, id string, userID string) error
}
