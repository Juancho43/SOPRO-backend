package usecases_test

import (
	"context"
	"testing"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	usecases "github.com/Juancho43/SOPRO-backend/pkg/internal/use_cases"
)

// TestUpdateMonthlyGoalUseCase verifies updating content and linked annual goal of an existing monthly goal.
func TestUpdateMonthlyGoalUseCase(t *testing.T) {
	repo := newMockFinanceRepo()
	repo.goals["g1"] = &entities.MonthlyGoal{
		ID:          "g1",
		UserID:      "u1",
		Year:        2026,
		Month:       11,
		Title:       "Original Title",
		Description: "Original Description",
		Status:      entities.MonthlyGoalStatusPending,
	}

	uc := usecases.NewUpdateMonthlyGoalUseCase(repo)

	t.Run("Updates title, description, and annual goal reference", func(t *testing.T) {
		newRef := "annual-5"
		resp, err := uc.Execute(context.Background(), usecases.UpdateMonthlyGoalCommand{
			ID:                    "g1",
			UserID:                "u1",
			Title:                 "Updated Tactical Goal",
			Description:           "Updated Description",
			AnnualGoalReferenceID: &newRef,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.Title != "Updated Tactical Goal" || resp.Description != "Updated Description" {
			t.Errorf("unexpected updated fields: %+v", resp.MonthlyGoal)
		}
		if resp.AnnualGoalReferenceID == nil || *resp.AnnualGoalReferenceID != "annual-5" {
			t.Errorf("expected annual_goal_reference_id 'annual-5', got %v", resp.AnnualGoalReferenceID)
		}
	})

	t.Run("Returns ErrMonthlyGoalNotFound for non-existent goal", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), usecases.UpdateMonthlyGoalCommand{
			ID:     "non-existent",
			UserID: "u1",
			Title:  "Valid Title",
		})
		if err != usecases.ErrMonthlyGoalNotFound {
			t.Errorf("expected ErrMonthlyGoalNotFound, got %v", err)
		}
	})
}
