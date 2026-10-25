package usecases

import (
	"context"
	"errors"
	"strings"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/repositories"
)

type CreateTrophyUseCase struct {
	trophyRepo repositories.TrophyRepository
	habitRepo  repositories.HabitRepository
}

func NewCreateTrophyUseCase(trophyRepo repositories.TrophyRepository, habitRepo repositories.HabitRepository) *CreateTrophyUseCase {
	return &CreateTrophyUseCase{
		trophyRepo: trophyRepo,
		habitRepo:  habitRepo,
	}
}

type TrophyConditionInput struct {
	HabitID            string `json:"habit_id"`
	RequiredStreakDays int    `json:"required_streak_days"`
}

type CreateTrophyCommand struct {
	UserID             string                 `json:"user_id"`
	Name               string                 `json:"name"`
	Description        string                 `json:"description"`
	Icon               string                 `json:"icon"`
	Color              string                 `json:"color"`
	TargetType         string                 `json:"target_type"` // "master_streak" | "habit" | "synergy"
	HabitID            *string                `json:"habit_id,omitempty"`
	RequiredStreakDays int                    `json:"required_streak_days"`
	Conditions         []TrophyConditionInput `json:"conditions,omitempty"`
}

func (uc *CreateTrophyUseCase) Execute(ctx context.Context, cmd CreateTrophyCommand) (*entities.Trophy, error) {
	name := strings.TrimSpace(cmd.Name)
	if name == "" {
		return nil, errors.New("el nombre del trofeo es obligatorio")
	}

	icon := strings.TrimSpace(cmd.Icon)
	if icon == "" {
		icon = "military_tech"
	}

	color := strings.TrimSpace(cmd.Color)
	if color == "" {
		color = "#FFD700"
	}

	targetType := entities.TrophyTargetType(strings.ToLower(strings.TrimSpace(cmd.TargetType)))
	if targetType != entities.TargetMasterStreak && targetType != entities.TargetHabit && targetType != entities.TargetSynergy {
		return nil, errors.New("target_type inválido: debe ser 'master_streak', 'habit' o 'synergy'")
	}

	var conditions []entities.TrophyCondition

	switch targetType {
	case entities.TargetMasterStreak:
		if cmd.RequiredStreakDays <= 0 {
			return nil, errors.New("required_streak_days debe ser mayor a 0 para master_streak")
		}
	case entities.TargetHabit:
		if cmd.HabitID == nil || strings.TrimSpace(*cmd.HabitID) == "" {
			return nil, errors.New("habit_id es obligatorio para trofeos de tipo habit")
		}
		if cmd.RequiredStreakDays <= 0 {
			return nil, errors.New("required_streak_days debe ser mayor a 0 para trofeos de hábito")
		}
		// Validar que el hábito pertenezca al usuario
		habit, err := uc.habitRepo.GetByID(ctx, *cmd.HabitID)
		if err != nil {
			return nil, err
		}
		if habit == nil || habit.UserID != cmd.UserID {
			return nil, errors.New("el hábito asociado no existe o no pertenece al usuario")
		}
	case entities.TargetSynergy:
		if len(cmd.Conditions) < 2 {
			return nil, errors.New("los trofeos de sinergia requieren al menos 2 condiciones de hábitos")
		}
		for _, cond := range cmd.Conditions {
			if strings.TrimSpace(cond.HabitID) == "" || cond.RequiredStreakDays <= 0 {
				return nil, errors.New("cada condición de sinergia debe incluir un habit_id válido y required_streak_days > 0")
			}
			habit, err := uc.habitRepo.GetByID(ctx, cond.HabitID)
			if err != nil {
				return nil, err
			}
			if habit == nil || habit.UserID != cmd.UserID {
				return nil, errors.New("uno de los hábitos en las condiciones no existe o no pertenece al usuario")
			}
			conditions = append(conditions, entities.TrophyCondition{
				HabitID:            cond.HabitID,
				RequiredStreakDays: cond.RequiredStreakDays,
			})
		}
	}

	trophy := &entities.Trophy{
		UserID:             cmd.UserID,
		Name:               name,
		Description:        strings.TrimSpace(cmd.Description),
		Icon:               icon,
		Color:              color,
		TargetType:         targetType,
		HabitID:            cmd.HabitID,
		RequiredStreakDays: cmd.RequiredStreakDays,
	}

	if err := uc.trophyRepo.Create(ctx, trophy, conditions); err != nil {
		return nil, err
	}

	return trophy, nil
}
