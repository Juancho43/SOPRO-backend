package usecases_test

import (
	"context"
	"testing"
	"time"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	usecases "github.com/Juancho43/SOPRO-backend/pkg/internal/use_cases"
)

// mockFinanceRepo implements repositories.FinanceRepository in memory for unit testing Epic 4 use cases.
type mockFinanceRepo struct {
	goals       map[string]*entities.MonthlyGoal
	budgetItems map[string]*entities.BudgetItem
}

// newMockFinanceRepo initializes an in-memory mock FinanceRepository.
func newMockFinanceRepo() *mockFinanceRepo {
	return &mockFinanceRepo{
		goals:       make(map[string]*entities.MonthlyGoal),
		budgetItems: make(map[string]*entities.BudgetItem),
	}
}

func (m *mockFinanceRepo) CreateMonthlyGoal(ctx context.Context, goal *entities.MonthlyGoal) error {
	if goal.ID == "" {
		goal.ID = "goal-1"
	}
	goal.CreatedAt = time.Now()
	goal.UpdatedAt = time.Now()
	m.goals[goal.ID] = goal
	return nil
}

func (m *mockFinanceRepo) GetMonthlyGoalByID(ctx context.Context, id string, userID string) (*entities.MonthlyGoal, error) {
	if g, ok := m.goals[id]; ok && g.UserID == userID {
		return g, nil
	}
	return nil, nil
}

func (m *mockFinanceRepo) ListMonthlyGoalsByPeriod(ctx context.Context, userID string, year int, month int) ([]*entities.MonthlyGoal, error) {
	var res []*entities.MonthlyGoal
	for _, g := range m.goals {
		if g.UserID == userID && g.Year == year && g.Month == month {
			res = append(res, g)
		}
	}
	return res, nil
}

func (m *mockFinanceRepo) UpdateMonthlyGoal(ctx context.Context, goal *entities.MonthlyGoal) error {
	goal.UpdatedAt = time.Now()
	m.goals[goal.ID] = goal
	return nil
}

func (m *mockFinanceRepo) UpdateMonthlyGoalStatus(ctx context.Context, id string, userID string, status entities.MonthlyGoalStatus, achievedAt *time.Time) error {
	if g, ok := m.goals[id]; ok && g.UserID == userID {
		g.Status = status
		g.AchievedAt = achievedAt
		g.UpdatedAt = time.Now()
	}
	return nil
}

func (m *mockFinanceRepo) DeleteMonthlyGoal(ctx context.Context, id string, userID string) error {
	delete(m.goals, id)
	return nil
}

func (m *mockFinanceRepo) CreateBudgetItem(ctx context.Context, item *entities.BudgetItem) error {
	if item.ID == "" {
		item.ID = "item-" + item.Concept
	}
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()
	m.budgetItems[item.ID] = item
	return nil
}

func (m *mockFinanceRepo) GetBudgetItemByID(ctx context.Context, id string, userID string) (*entities.BudgetItem, error) {
	if it, ok := m.budgetItems[id]; ok && it.UserID == userID {
		return it, nil
	}
	return nil, nil
}

func (m *mockFinanceRepo) ListBudgetItemsByPeriod(ctx context.Context, userID string, year int, month int) ([]*entities.BudgetItem, error) {
	var res []*entities.BudgetItem
	for _, it := range m.budgetItems {
		if it.UserID == userID && it.Year == year && it.Month == month {
			res = append(res, it)
		}
	}
	return res, nil
}

func (m *mockFinanceRepo) UpdateBudgetItem(ctx context.Context, item *entities.BudgetItem) error {
	item.UpdatedAt = time.Now()
	m.budgetItems[item.ID] = item
	return nil
}

func (m *mockFinanceRepo) DeleteBudgetItem(ctx context.Context, id string, userID string) error {
	delete(m.budgetItems, id)
	return nil
}

// TestCreateMonthlyGoalUseCase verifies monthly goal creation and domain validations.
func TestCreateMonthlyGoalUseCase(t *testing.T) {
	repo := newMockFinanceRepo()
	uc, err := usecases.NewCreateMonthlyGoalUseCase(repo, "America/Buenos_Aires")
	if err != nil {
		t.Fatalf("unexpected error constructing use case: %v", err)
	}

	t.Run("Creates a monthly goal with PENDING initial status and linked annual goal", func(t *testing.T) {
		ref := "annual-goal-2"
		resp, err := uc.Execute(context.Background(), usecases.CreateMonthlyGoalCommand{
			UserID:                "u1",
			Year:                  2026,
			Month:                 11,
			Title:                 "Launch v1 of the SaaS platform",
			Description:           "Complete onboarding and billing flow",
			AnnualGoalReferenceID: &ref,
		})
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if resp.Status != entities.MonthlyGoalStatusPending {
			t.Errorf("expected initial status PENDING, got %s", resp.Status)
		}
		if resp.AnnualGoalReferenceID == nil || *resp.AnnualGoalReferenceID != ref {
			t.Errorf("expected annual_goal_reference_id %q, got %v", ref, resp.AnnualGoalReferenceID)
		}
	})

	t.Run("Rejects empty title or invalid month", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), usecases.CreateMonthlyGoalCommand{
			UserID: "u1",
			Year:   2026,
			Month:  11,
			Title:  "   ",
		})
		if err == nil {
			t.Errorf("expected validation error for empty title")
		}

		_, err = uc.Execute(context.Background(), usecases.CreateMonthlyGoalCommand{
			UserID: "u1",
			Year:   2026,
			Month:  13,
			Title:  "Valid title",
		})
		if err == nil {
			t.Errorf("expected validation error for month > 12")
		}
	})
}
