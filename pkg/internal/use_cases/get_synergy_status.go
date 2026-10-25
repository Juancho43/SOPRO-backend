package usecases

import (
	"context"
	"errors"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/repositories"
)

type GetSynergyStatusUseCase struct {
	trophyRepo repositories.TrophyRepository
	habitRepo  repositories.HabitRepository
}

func NewGetSynergyStatusUseCase(trophyRepo repositories.TrophyRepository, habitRepo repositories.HabitRepository) *GetSynergyStatusUseCase {
	return &GetSynergyStatusUseCase{
		trophyRepo: trophyRepo,
		habitRepo:  habitRepo,
	}
}

func (uc *GetSynergyStatusUseCase) Execute(ctx context.Context, trophyID string, userID string) (*entities.SynergyStatusResponse, error) {
	if trophyID == "" || userID == "" {
		return nil, errors.New("trophyID y userID son requeridos")
	}

	trophy, err := uc.trophyRepo.GetByID(ctx, trophyID)
	if err != nil {
		return nil, err
	}
	if trophy == nil || trophy.UserID != userID {
		return nil, errors.New("trofeo no encontrado")
	}

	if trophy.TargetType != entities.TargetSynergy {
		return nil, errors.New("el trofeo especificado no es de tipo sinergia")
	}

	unlocked, err := uc.trophyRepo.IsTrophyUnlocked(ctx, userID, trophyID)
	if err != nil {
		return nil, err
	}

	conditions, err := uc.trophyRepo.GetConditionsByTrophyID(ctx, trophyID)
	if err != nil {
		return nil, err
	}

	var progress []entities.SynergyHabitProgress
	for _, cond := range conditions {
		habit, err := uc.habitRepo.GetByID(ctx, cond.HabitID)
		habitName := "Hábito desconocido"
		if err == nil && habit != nil {
			habitName = habit.Name
		}

		currentStreak, _, _, _ := uc.habitRepo.GetStreakByHabitID(ctx, cond.HabitID)
		progress = append(progress, entities.SynergyHabitProgress{
			HabitID:            cond.HabitID,
			HabitName:          habitName,
			CurrentStreak:      currentStreak,
			RequiredStreakDays: cond.RequiredStreakDays,
			IsMet:              currentStreak >= cond.RequiredStreakDays,
		})
	}

	return &entities.SynergyStatusResponse{
		TrophyID:   trophy.ID,
		TrophyName: trophy.Name,
		Unlocked:   unlocked,
		Progress:   progress,
	}, nil
}
