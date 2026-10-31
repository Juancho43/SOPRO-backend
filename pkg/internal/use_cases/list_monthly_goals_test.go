package usecases_test

import (
	"context"
	"testing"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	usecases "github.com/Juancho43/SOPRO-backend/pkg/internal/use_cases"
)

// TestListMonthlyGoalsUseCase verifies period filtering when listing monthly goals.
func TestListMonthlyGoalsUseCase(t *testing.T) {
	repo := newMockFinanceRepo()
	repo.goals["g1"] = &entities.MonthlyGoal{
		ID:     "g1",
		UserID: "u1",
		Year:   2026,
		Month:  10,
		Title:  "October Goal",
		Status: entities.MonthlyGoalStatusAchieved,
	}
	repo.goals["g2"] = &entities.MonthlyGoal{
		ID:     "g2",
		UserID: "u1",
		Year:   2026,
		Month:  11,
		Title:  "November Goal",
		Status: entities.MonthlyGoalStatusPending,
	}

	uc, err := usecases.NewListMonthlyGoalsUseCase(repo, "America/Buenos_Aires")
	if err != nil {
		t.Fatalf("unexpected error constructing use case: %v", err)
	}

	resp, err := uc.Execute(context.Background(), usecases.ListMonthlyGoalsQuery{
		UserID: "u1",
		Year:   2026,
		Month:  11,
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(resp.Goals) != 1 || resp.Goals[0].ID != "g2" {
		t.Errorf("expected 1 goal (g2) for November 2026, got %d", len(resp.Goals))
	}
}
