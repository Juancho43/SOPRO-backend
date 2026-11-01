package usecases_test

import (
	"context"
	"testing"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	usecases "github.com/Juancho43/SOPRO-backend/pkg/internal/use_cases"
)

// TestUpdateMonthlyGoalStatusUseCase verifies status transitions and achieved_at timestamp handling.
func TestUpdateMonthlyGoalStatusUseCase(t *testing.T) {
	repo := newMockFinanceRepo()
	repo.goals["g1"] = &entities.MonthlyGoal{
		ID:     "g1",
		UserID: "u1",
		Year:   2026,
		Month:  11,
		Title:  "Close 3 enterprise deals",
		Status: entities.MonthlyGoalStatusPending,
	}

	uc, err := usecases.NewUpdateMonthlyGoalStatusUseCase(repo, "America/Buenos_Aires")
	if err != nil {
		t.Fatalf("unexpected error constructing use case: %v", err)
	}

	t.Run("Transitions to ACHIEVED and populates achieved_at", func(t *testing.T) {
		resp, err := uc.Execute(context.Background(), usecases.UpdateMonthlyGoalStatusCommand{
			ID:     "g1",
			UserID: "u1",
			Status: "ACHIEVED",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.Status != entities.MonthlyGoalStatusAchieved || resp.AchievedAt == nil {
			t.Errorf("expected ACHIEVED status and non-nil achieved_at")
		}
	})

	t.Run("Transitions back to IN_PROGRESS and clears achieved_at", func(t *testing.T) {
		resp, err := uc.Execute(context.Background(), usecases.UpdateMonthlyGoalStatusCommand{
			ID:     "g1",
			UserID: "u1",
			Status: "IN_PROGRESS",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.Status != entities.MonthlyGoalStatusInProgress || resp.AchievedAt != nil {
			t.Errorf("expected IN_PROGRESS status and nil achieved_at")
		}
	})

	t.Run("Rejects invalid status values", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), usecases.UpdateMonthlyGoalStatusCommand{
			ID:     "g1",
			UserID: "u1",
			Status: "INVALID_STATE",
		})
		if err == nil {
			t.Errorf("expected error for invalid status")
		}
	})
}
