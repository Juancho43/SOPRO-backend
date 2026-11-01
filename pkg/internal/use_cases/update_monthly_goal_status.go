package usecases

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/repositories"
)

// UpdateMonthlyGoalStatusCommand represents the input payload to transition the status of a monthly tactical goal (HU 4.1).
type UpdateMonthlyGoalStatusCommand struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Status string `json:"status"`
}

// UpdateMonthlyGoalStatusResponse represents the monthly tactical goal after its status transition.
type UpdateMonthlyGoalStatusResponse struct {
	*entities.MonthlyGoal
}

// UpdateMonthlyGoalStatusUseCase handles state transitions (PENDING, IN_PROGRESS, ACHIEVED) for monthly goals.
type UpdateMonthlyGoalStatusUseCase struct {
	repo     repositories.FinanceRepository
	timezone *time.Location
}

// NewUpdateMonthlyGoalStatusUseCase constructs a new UpdateMonthlyGoalStatusUseCase instance.
func NewUpdateMonthlyGoalStatusUseCase(repo repositories.FinanceRepository, tz string) (*UpdateMonthlyGoalStatusUseCase, error) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		loc = time.Local
	}
	return &UpdateMonthlyGoalStatusUseCase{
		repo:     repo,
		timezone: loc,
	}, nil
}

// Execute validates the requested status, updates the goal's status and achieved_at timestamp, and returns the updated goal.
func (uc *UpdateMonthlyGoalStatusUseCase) Execute(ctx context.Context, cmd UpdateMonthlyGoalStatusCommand) (*UpdateMonthlyGoalStatusResponse, error) {
	status := entities.MonthlyGoalStatus(strings.ToUpper(strings.TrimSpace(cmd.Status)))
	switch status {
	case entities.MonthlyGoalStatusPending, entities.MonthlyGoalStatusInProgress, entities.MonthlyGoalStatusAchieved:
	default:
		return nil, errors.New("estado de meta mensual no válido (valores permitidos: PENDING, IN_PROGRESS, ACHIEVED)")
	}

	existing, err := uc.repo.GetMonthlyGoalByID(ctx, cmd.ID, cmd.UserID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrMonthlyGoalNotFound
	}

	var achievedAt *time.Time
	if status == entities.MonthlyGoalStatusAchieved {
		now := time.Now().In(uc.timezone)
		achievedAt = &now
	}

	if err := uc.repo.UpdateMonthlyGoalStatus(ctx, cmd.ID, cmd.UserID, status, achievedAt); err != nil {
		return nil, err
	}

	existing.Status = status
	existing.AchievedAt = achievedAt
	return &UpdateMonthlyGoalStatusResponse{MonthlyGoal: existing}, nil
}
