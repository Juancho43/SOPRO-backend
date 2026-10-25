package usecases

import (
	"context"
	"time"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/repositories"
	"github.com/Juancho43/SOPRO-backend/pkg/src/bus"
)

type EvaluateSynergyTrophiesUseCase struct {
	trophyRepo repositories.TrophyRepository
	habitRepo  repositories.HabitRepository
	eventBus   bus.EventBus
}

func NewEvaluateSynergyTrophiesUseCase(
	trophyRepo repositories.TrophyRepository,
	habitRepo repositories.HabitRepository,
	eventBus bus.EventBus,
) *EvaluateSynergyTrophiesUseCase {
	return &EvaluateSynergyTrophiesUseCase{
		trophyRepo: trophyRepo,
		habitRepo:  habitRepo,
		eventBus:   eventBus,
	}
}

func (uc *EvaluateSynergyTrophiesUseCase) Execute(ctx context.Context, userID string, habitID string) error {
	synergyTrophies, err := uc.trophyRepo.FindSynergyTrophiesByHabit(ctx, habitID)
	if err != nil {
		return err
	}

	for _, trophy := range synergyTrophies {
		// Validar que pertenezca al usuario del evento
		if trophy.UserID != userID {
			continue
		}

		// Si ya está desbloqueado, no hacemos nada (idempotencia)
		alreadyUnlocked, err := uc.trophyRepo.IsTrophyUnlocked(ctx, userID, trophy.ID)
		if err != nil || alreadyUnlocked {
			continue
		}

		conditions, err := uc.trophyRepo.GetConditionsByTrophyID(ctx, trophy.ID)
		if err != nil || len(conditions) == 0 {
			continue
		}

		allMet := true
		for _, cond := range conditions {
			currentStreak, _, _, err := uc.habitRepo.GetStreakByHabitID(ctx, cond.HabitID)
			if err != nil || currentStreak < cond.RequiredStreakDays {
				allMet = false
				break
			}
		}

		if allMet {
			unlocked, err := uc.trophyRepo.UnlockTrophy(ctx, userID, trophy.ID)
			if err != nil {
				return err
			}

			if unlocked && uc.eventBus != nil {
				event := entities.TrophyUnlockedEvent{
					UserID:     userID,
					TrophyID:   trophy.ID,
					TrophyName: trophy.Name,
					Icon:       trophy.Icon,
					Color:      trophy.Color,
					UnlockedAt: time.Now(),
					IsSynergy:  true,
				}
				uc.eventBus.Publish("trophy.unlocked", event)
			}
		}
	}

	return nil
}
