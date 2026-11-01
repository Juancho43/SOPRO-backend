package usecases

import (
	"context"
	"errors"
	"strings"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/repositories"
)

// UpdateMonthlyGoalCommand represents the input payload to update a monthly tactical goal's content or linked annual goal (HU 4.1).
type UpdateMonthlyGoalCommand struct {
	ID                         string  `json:"id"`
	UserID                     string  `json:"user_id"`
	Title                      string  `json:"title"`
	Description                string  `json:"description"`
	AnnualGoalReferenceID      *string `json:"annual_goal_reference_id,omitempty"`
	AnnualGoalReferenceIDCamel *string `json:"annualGoalReferenceId,omitempty" swaggerignore:"true"`
}

// UpdateMonthlyGoalResponse represents the updated monthly tactical goal returned by the use case.
type UpdateMonthlyGoalResponse struct {
	*entities.MonthlyGoal
}

// UpdateMonthlyGoalUseCase handles updating the content and annual goal reference of an existing monthly goal.
type UpdateMonthlyGoalUseCase struct {
	repo repositories.FinanceRepository
}

// NewUpdateMonthlyGoalUseCase constructs a new UpdateMonthlyGoalUseCase instance.
func NewUpdateMonthlyGoalUseCase(repo repositories.FinanceRepository) *UpdateMonthlyGoalUseCase {
	return &UpdateMonthlyGoalUseCase{repo: repo}
}

// Execute validates the command, updates the target monthly tactical goal, and returns the updated entity.
func (uc *UpdateMonthlyGoalUseCase) Execute(ctx context.Context, cmd UpdateMonthlyGoalCommand) (*UpdateMonthlyGoalResponse, error) {
	title := strings.TrimSpace(cmd.Title)
	if title == "" {
		return nil, errors.New("el título de la meta mensual no puede estar vacío")
	}
	if len(title) > 255 {
		return nil, errors.New("el título de la meta mensual no puede exceder 255 caracteres")
	}

	existing, err := uc.repo.GetMonthlyGoalByID(ctx, cmd.ID, cmd.UserID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrMonthlyGoalNotFound
	}

	existing.Title = title
	existing.Description = strings.TrimSpace(cmd.Description)
	if cmd.AnnualGoalReferenceID != nil || cmd.AnnualGoalReferenceIDCamel != nil {
		existing.AnnualGoalReferenceID = normalizeOptionalString(cmd.AnnualGoalReferenceID, cmd.AnnualGoalReferenceIDCamel)
	}

	if err := uc.repo.UpdateMonthlyGoal(ctx, existing); err != nil {
		return nil, err
	}

	return &UpdateMonthlyGoalResponse{MonthlyGoal: existing}, nil
}
