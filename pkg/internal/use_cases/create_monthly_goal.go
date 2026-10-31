package usecases

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/repositories"
)

// ErrMonthlyGoalNotFound is returned when a monthly tactical goal does not exist or does not belong to the user.
var ErrMonthlyGoalNotFound = errors.New("meta mensual no encontrada")

// CreateMonthlyGoalCommand represents the input payload to create a new monthly tactical goal (HU 4.1).
type CreateMonthlyGoalCommand struct {
	UserID                     string  `json:"user_id"`
	Year                       int     `json:"year"`
	Month                      int     `json:"month"`
	Title                      string  `json:"title"`
	Description                string  `json:"description"`
	AnnualGoalReferenceID      *string `json:"annual_goal_reference_id,omitempty"`
	AnnualGoalReferenceIDCamel *string `json:"annualGoalReferenceId,omitempty" swaggerignore:"true"`
}

// CreateMonthlyGoalResponse represents the output returned after creating a monthly tactical goal.
type CreateMonthlyGoalResponse struct {
	*entities.MonthlyGoal
}

// CreateMonthlyGoalUseCase handles the creation of a new monthly tactical goal with an initial PENDING status.
type CreateMonthlyGoalUseCase struct {
	repo     repositories.FinanceRepository
	timezone *time.Location
}

// NewCreateMonthlyGoalUseCase constructs a new CreateMonthlyGoalUseCase with the given repository and timezone.
func NewCreateMonthlyGoalUseCase(repo repositories.FinanceRepository, tz string) (*CreateMonthlyGoalUseCase, error) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		loc = time.Local
	}
	return &CreateMonthlyGoalUseCase{
		repo:     repo,
		timezone: loc,
	}, nil
}

// Execute validates the command and persists a new monthly tactical goal with initial status PENDING.
func (uc *CreateMonthlyGoalUseCase) Execute(ctx context.Context, cmd CreateMonthlyGoalCommand) (*CreateMonthlyGoalResponse, error) {
	title := strings.TrimSpace(cmd.Title)
	if title == "" {
		return nil, errors.New("la ejecución táctica requiere claridad: el título de la meta mensual no puede estar vacío")
	}
	if len(title) > 255 {
		return nil, errors.New("el título de la meta mensual no puede exceder 255 caracteres")
	}

	year, month, err := normalizeFinancePeriod(cmd.Year, cmd.Month, uc.timezone)
	if err != nil {
		return nil, err
	}

	refID := normalizeOptionalString(cmd.AnnualGoalReferenceID, cmd.AnnualGoalReferenceIDCamel)

	goal := &entities.MonthlyGoal{
		UserID:                cmd.UserID,
		Year:                  year,
		Month:                 month,
		Title:                 title,
		Description:           strings.TrimSpace(cmd.Description),
		AnnualGoalReferenceID: refID,
		Status:                entities.MonthlyGoalStatusPending,
	}

	if err := uc.repo.CreateMonthlyGoal(ctx, goal); err != nil {
		return nil, err
	}

	return &CreateMonthlyGoalResponse{MonthlyGoal: goal}, nil
}

// normalizeFinancePeriod resolves default year/month when omitted (0) and validates their ranges.
func normalizeFinancePeriod(year int, month int, loc *time.Location) (int, int, error) {
	if loc == nil {
		loc = time.Local
	}
	now := time.Now().In(loc)
	if year == 0 {
		year = now.Year()
	}
	if month == 0 {
		month = int(now.Month())
	}
	if year < 2000 || year > 2100 {
		return 0, 0, errors.New("el año debe estar comprendido entre 2000 y 2100")
	}
	if month < 1 || month > 12 {
		return 0, 0, errors.New("el mes debe estar comprendido entre 1 y 12")
	}
	return year, month, nil
}

// normalizeOptionalString trims and resolves an optional string pointer with a fallback pointer.
func normalizeOptionalString(primary *string, fallback *string) *string {
	val := primary
	if val == nil {
		val = fallback
	}
	if val == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*val)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}
