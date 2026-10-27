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

type CreateChallengeUseCase struct {
	repo repositories.ChallengeRepository
}

func NewCreateChallengeUseCase(repo repositories.ChallengeRepository) *CreateChallengeUseCase {
	return &CreateChallengeUseCase{repo: repo}
}

type CreateChallengeCommand struct {
	UserID      string `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	StartDate   string `json:"start_date"` // YYYY-MM-DD
	EndDate     string `json:"end_date"`   // YYYY-MM-DD
	Reward      string `json:"reward"`
}

func (uc *CreateChallengeUseCase) Execute(ctx context.Context, cmd CreateChallengeCommand) (*entities.Challenge, error) {
	title := strings.TrimSpace(cmd.Title)
	if title == "" {
		return nil, errors.New("la visión requiere un nombre: el título del reto no puede estar vacío")
	}
	if len(title) > 120 {
		return nil, errors.New("el título del reto no puede exceder 120 caracteres")
	}

	reward := strings.TrimSpace(cmd.Reward)
	if reward == "" {
		return nil, errors.New("la conquista merece un premio: la recompensa final es obligatoria")
	}

	startDateStr := strings.TrimSpace(cmd.StartDate)
	if startDateStr == "" {
		startDateStr = time.Now().Format("2006-01-02")
	}

	start, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return nil, fmt.Errorf("formato de fecha de inicio inválido (se espera YYYY-MM-DD): %w", err)
	}

	endDateStr := strings.TrimSpace(cmd.EndDate)
	if endDateStr == "" {
		return nil, errors.New("un objetivo sin fecha límite es solo un sueño: la fecha de fin es obligatoria")
	}

	end, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return nil, fmt.Errorf("formato de fecha de fin inválido (se espera YYYY-MM-DD): %w", err)
	}

	if !end.After(start) {
		return nil, errors.New("la fecha de fin debe ser posterior a la fecha de inicio")
	}

	challenge := &entities.Challenge{
		UserID:      cmd.UserID,
		Title:       title,
		Description: strings.TrimSpace(cmd.Description),
		StartDate:   startDateStr,
		EndDate:     endDateStr,
		Reward:      reward,
		Status:      entities.ChallengeStatusActive,
	}

	if err := uc.repo.CreateChallenge(ctx, challenge); err != nil {
		return nil, err
	}

	return challenge, nil
}
