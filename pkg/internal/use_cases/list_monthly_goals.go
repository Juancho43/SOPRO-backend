package usecases

import (
	"context"
	"time"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/repositories"
)

// ListMonthlyGoalsQuery represents the query parameters for retrieving monthly tactical goals (HU 4.1).
type ListMonthlyGoalsQuery struct {
	UserID string `json:"user_id"`
	Year   int    `json:"year"`
	Month  int    `json:"month"`
}

// ListMonthlyGoalsResponse wraps the list of monthly tactical goals for the requested period.
type ListMonthlyGoalsResponse struct {
	Year  int                     `json:"year"`
	Month int                     `json:"month"`
	Goals []*entities.MonthlyGoal `json:"goals"`
}

// ListMonthlyGoalsUseCase retrieves all monthly tactical goals for a specific user and period.
type ListMonthlyGoalsUseCase struct {
	repo     repositories.FinanceRepository
	timezone *time.Location
}

// NewListMonthlyGoalsUseCase constructs a new ListMonthlyGoalsUseCase instance.
func NewListMonthlyGoalsUseCase(repo repositories.FinanceRepository, tz string) (*ListMonthlyGoalsUseCase, error) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		loc = time.Local
	}
	return &ListMonthlyGoalsUseCase{
		repo:     repo,
		timezone: loc,
	}, nil
}

// Execute validates the query period and returns the list of monthly goals for the user.
func (uc *ListMonthlyGoalsUseCase) Execute(ctx context.Context, query ListMonthlyGoalsQuery) (*ListMonthlyGoalsResponse, error) {
	year, month, err := normalizeFinancePeriod(query.Year, query.Month, uc.timezone)
	if err != nil {
		return nil, err
	}

	goals, err := uc.repo.ListMonthlyGoalsByPeriod(ctx, query.UserID, year, month)
	if err != nil {
		return nil, err
	}
	if goals == nil {
		goals = []*entities.MonthlyGoal{}
	}

	return &ListMonthlyGoalsResponse{
		Year:  year,
		Month: month,
		Goals: goals,
	}, nil
}
