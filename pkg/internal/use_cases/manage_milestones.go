package usecases

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/repositories"
)

type ManageMilestonesUseCase struct {
	repo     repositories.ChallengeRepository
	timezone *time.Location
}

func NewManageMilestonesUseCase(repo repositories.ChallengeRepository, tz string) (*ManageMilestonesUseCase, error) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		loc = time.Local
	}
	return &ManageMilestonesUseCase{
		repo:     repo,
		timezone: loc,
	}, nil
}

type CreateMilestoneCommand struct {
	UserID      string `json:"user_id"`
	ChallengeID string `json:"challenge_id"`
	Title       string `json:"title"`
	DueDate     string `json:"due_date"` // YYYY-MM-DD
}

func (uc *ManageMilestonesUseCase) Create(ctx context.Context, cmd CreateMilestoneCommand) (*entities.Milestone, error) {
	title := strings.TrimSpace(cmd.Title)
	if title == "" {
		return nil, errors.New("el hito requiere un objetivo concreto: el título no puede estar vacío")
	}

	dueDateStr := strings.TrimSpace(cmd.DueDate)
	if dueDateStr == "" {
		return nil, errors.New("la fecha límite intermedia (due_date) es obligatoria")
	}

	due, err := time.Parse("2006-01-02", dueDateStr)
	if err != nil {
		return nil, fmt.Errorf("formato de fecha límite inválido (se espera YYYY-MM-DD): %w", err)
	}

	// Validar existencia y fechas del reto padre
	challenge, err := uc.repo.GetChallengeByID(ctx, cmd.ChallengeID, cmd.UserID)
	if err != nil {
		return nil, err
	}
	if challenge == nil {
		return nil, errors.New("el reto asociado no existe o no pertenece al usuario")
	}

	if challenge.Status != entities.ChallengeStatusActive {
		return nil, errors.New("no se pueden agregar hitos a un reto que no está activo")
	}

	start, _ := time.Parse("2006-01-02", challenge.StartDate)
	end, _ := time.Parse("2006-01-02", challenge.EndDate)

	if due.Before(start) || due.After(end) {
		return nil, fmt.Errorf("la fecha límite del hito (%s) debe estar entre el inicio (%s) y fin (%s) del reto",
			dueDateStr, challenge.StartDate, challenge.EndDate)
	}

	milestones, err := uc.repo.ListMilestonesByChallengeID(ctx, cmd.ChallengeID)
	if err != nil {
		return nil, err
	}

	milestone := &entities.Milestone{
		ChallengeID: cmd.ChallengeID,
		Title:       title,
		DueDate:     dueDateStr,
		IsCompleted: false,
		OrderIndex:  len(milestones) + 1,
	}

	if err := uc.repo.CreateMilestone(ctx, milestone); err != nil {
		return nil, err
	}

	return milestone, nil
}

type UpdateMilestoneCommand struct {
	ID          string `json:"id"`
	ChallengeID string `json:"challenge_id"`
	UserID      string `json:"user_id"`
	Title       *string `json:"title,omitempty"`
	DueDate     *string `json:"due_date,omitempty"`
	IsCompleted *bool   `json:"is_completed,omitempty"`
}

func (uc *ManageMilestonesUseCase) Update(ctx context.Context, cmd UpdateMilestoneCommand) (*entities.Milestone, error) {
	// Verificar pertenencia del reto
	challenge, err := uc.repo.GetChallengeByID(ctx, cmd.ChallengeID, cmd.UserID)
	if err != nil {
		return nil, err
	}
	if challenge == nil {
		return nil, errors.New("reto no encontrado o no autorizado")
	}

	milestone, err := uc.repo.GetMilestoneByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}
	if milestone == nil || milestone.ChallengeID != cmd.ChallengeID {
		return nil, errors.New("hito no encontrado")
	}

	if cmd.Title != nil {
		title := strings.TrimSpace(*cmd.Title)
		if title == "" {
			return nil, errors.New("el título del hito no puede estar vacío")
		}
		milestone.Title = title
	}

	if cmd.DueDate != nil {
		dueStr := strings.TrimSpace(*cmd.DueDate)
		due, err := time.Parse("2006-01-02", dueStr)
		if err != nil {
			return nil, errors.New("formato de due_date inválido")
		}

		start, _ := time.Parse("2006-01-02", challenge.StartDate)
		end, _ := time.Parse("2006-01-02", challenge.EndDate)
		if due.Before(start) || due.After(end) {
			return nil, fmt.Errorf("la fecha límite (%s) debe estar entre %s y %s", dueStr, challenge.StartDate, challenge.EndDate)
		}
		milestone.DueDate = dueStr
	}

	if cmd.IsCompleted != nil {
		milestone.IsCompleted = *cmd.IsCompleted
		if *cmd.IsCompleted {
			now := time.Now()
			milestone.CompletedAt = &now
		} else {
			milestone.CompletedAt = nil
		}
	}

	if err := uc.repo.UpdateMilestone(ctx, milestone); err != nil {
		return nil, err
	}

	return milestone, nil
}

func (uc *ManageMilestonesUseCase) Delete(ctx context.Context, milestoneID string, challengeID string, userID string) error {
	challenge, err := uc.repo.GetChallengeByID(ctx, challengeID, userID)
	if err != nil {
		return err
	}
	if challenge == nil {
		return errors.New("reto no encontrado o no autorizado")
	}

	return uc.repo.DeleteMilestone(ctx, milestoneID, challengeID)
}

func (uc *ManageMilestonesUseCase) ListUpcoming(ctx context.Context, userID string, maxDaysAhead int, now time.Time) ([]*entities.UpcomingDeadline, error) {
	if maxDaysAhead <= 0 {
		maxDaysAhead = 3
	}

	userNow := now.In(uc.timezone)
	todayStr := userNow.Format("2006-01-02")

	return uc.repo.ListUpcomingMilestones(ctx, userID, maxDaysAhead, todayStr)
}
