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

type ManageChallengesUseCase struct {
	repo     repositories.ChallengeRepository
	timezone *time.Location
}

func NewManageChallengesUseCase(repo repositories.ChallengeRepository, tz string) (*ManageChallengesUseCase, error) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		loc = time.Local
	}
	return &ManageChallengesUseCase{
		repo:     repo,
		timezone: loc,
	}, nil
}

type UpdateChallengeCommand struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Reward      string `json:"reward"`
}

func (uc *ManageChallengesUseCase) Update(ctx context.Context, cmd UpdateChallengeCommand) (*entities.Challenge, error) {
	title := strings.TrimSpace(cmd.Title)
	if title == "" {
		return nil, errors.New("el título del reto no puede estar vacío")
	}
	if len(title) > 120 {
		return nil, errors.New("el título del reto no puede exceder 120 caracteres")
	}

	reward := strings.TrimSpace(cmd.Reward)
	if reward == "" {
		return nil, errors.New("la recompensa final no puede estar vacía")
	}

	existing, err := uc.repo.GetChallengeByID(ctx, cmd.ID, cmd.UserID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("reto no encontrado")
	}

	existing.Title = title
	existing.Description = strings.TrimSpace(cmd.Description)
	existing.Reward = reward

	if err := uc.repo.UpdateChallenge(ctx, existing); err != nil {
		return nil, err
	}

	return existing, nil
}

func (uc *ManageChallengesUseCase) UpdateStatus(ctx context.Context, id string, userID string, status entities.ChallengeStatus) error {
	switch status {
	case entities.ChallengeStatusActive, entities.ChallengeStatusCompleted,
		entities.ChallengeStatusFailed, entities.ChallengeStatusArchived, entities.ChallengeStatusCancelled:
	default:
		return errors.New("estado de reto no válido")
	}

	return uc.repo.UpdateChallengeStatus(ctx, id, userID, status)
}

func (uc *ManageChallengesUseCase) List(ctx context.Context, userID string, statusFilter *string) ([]*entities.Challenge, error) {
	var status *entities.ChallengeStatus
	if statusFilter != nil && strings.TrimSpace(*statusFilter) != "" {
		s := entities.ChallengeStatus(strings.ToUpper(strings.TrimSpace(*statusFilter)))
		status = &s
	}
	return uc.repo.ListChallengesByUserID(ctx, userID, status)
}

func (uc *ManageChallengesUseCase) GetDetail(ctx context.Context, id string, userID string, now time.Time) (*entities.ChallengeDetail, error) {
	challenge, err := uc.repo.GetChallengeByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if challenge == nil {
		return nil, errors.New("reto no encontrado")
	}

	milestones, err := uc.repo.ListMilestonesByChallengeID(ctx, id)
	if err != nil {
		return nil, err
	}

	total := len(milestones)
	done := 0
	for _, m := range milestones {
		if m.IsCompleted {
			done++
		}
	}

	progress := 0.0
	if total > 0 {
		progress = math.Round((float64(done)/float64(total))*10000) / 100
	}

	// Cálculo de días restantes y progreso temporal
	userNow := now.In(uc.timezone)
	todayDate := time.Date(userNow.Year(), userNow.Month(), userNow.Day(), 0, 0, 0, 0, uc.timezone)

	start, err := time.ParseInLocation("2006-01-02", challenge.StartDate, uc.timezone)
	if err != nil {
		start = todayDate
	}
	end, err := time.ParseInLocation("2006-01-02", challenge.EndDate, uc.timezone)
	if err != nil {
		end = todayDate
	}

	totalDays := int(end.Sub(start).Hours() / 24)
	daysRemaining := int(end.Sub(todayDate).Hours() / 24)
	if daysRemaining < 0 {
		daysRemaining = 0
	}
	daysElapsed := int(todayDate.Sub(start).Hours() / 24)
	if daysElapsed < 0 {
		daysElapsed = 0
	}

	detail := &entities.ChallengeDetail{
		Challenge:       *challenge,
		Milestones:      milestones,
		TotalMilestones: total,
		DoneMilestones:  done,
		ProgressPercent: progress,
		DaysRemaining:   daysRemaining,
		TotalDays:       totalDays,
		DaysElapsed:     daysElapsed,
	}

	return detail, nil
}
