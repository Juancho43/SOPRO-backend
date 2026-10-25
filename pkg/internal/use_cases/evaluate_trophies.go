package usecases

import (
	"context"
	"time"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/repositories"
	"github.com/Juancho43/SOPRO-backend/pkg/src/bus"
)

type EvaluateTrophiesUseCase struct {
	trophyRepo repositories.TrophyRepository
	eventBus   bus.EventBus
}

func NewEvaluateTrophiesUseCase(trophyRepo repositories.TrophyRepository, eventBus bus.EventBus) *EvaluateTrophiesUseCase {
	return &EvaluateTrophiesUseCase{
		trophyRepo: trophyRepo,
		eventBus:   eventBus,
	}
}

func (uc *EvaluateTrophiesUseCase) Execute(
	ctx context.Context,
	userID string,
	targetType entities.TrophyTargetType,
	habitID *string,
	currentStreak int,
) error {
	candidates, err := uc.trophyRepo.FindActiveTrophiesByTarget(ctx, userID, targetType, habitID)
	if err != nil {
		return err
	}

	for _, trophy := range candidates {
		if currentStreak >= trophy.RequiredStreakDays {
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
					IsSynergy:  false,
				}
				uc.eventBus.Publish("trophy.unlocked", event)
			}
		}
	}

	return nil
}
