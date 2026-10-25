package usecases

import (
	"context"
	"errors"
	"strings"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/repositories"
)

type CreateHabitUseCase struct {
	habitRepo repositories.HabitRepository
}

func NewCreateHabitUseCase(repo repositories.HabitRepository) *CreateHabitUseCase {
	return &CreateHabitUseCase{habitRepo: repo}
}

type CreateHabitCommand struct {
	UserID    string  `json:"user_id"`
	Name      string  `json:"name"`
	Frequency string  `json:"frequency"` // "daily" | "weekly"
	DomainID  *string `json:"domain_id,omitempty"`
}

func (uc *CreateHabitUseCase) Execute(ctx context.Context, cmd CreateHabitCommand) (*entities.Habit, error) {
	name := strings.TrimSpace(cmd.Name)
	if name == "" {
		return nil, errors.New("la claridad es esencial: el nombre del hábito no puede estar vacío")
	}
	if len(name) > 255 {
		return nil, errors.New("el nombre del hábito no puede superar los 255 caracteres")
	}

	freq := strings.ToLower(strings.TrimSpace(cmd.Frequency))
	if freq == "" {
		freq = "daily"
	}
	if freq != "daily" && freq != "weekly" {
		return nil, errors.New("la frecuencia debe ser 'daily' o 'weekly'")
	}

	var domainID *string
	if cmd.DomainID != nil && strings.TrimSpace(*cmd.DomainID) != "" {
		trimmed := strings.TrimSpace(*cmd.DomainID)
		domainID = &trimmed
	}

	habit := &entities.Habit{
		UserID:    cmd.UserID,
		Name:      name,
		Frequency: freq,
		DomainID:  domainID,
	}

	if err := uc.habitRepo.Create(ctx, habit); err != nil {
		return nil, err
	}

	return habit, nil
}
