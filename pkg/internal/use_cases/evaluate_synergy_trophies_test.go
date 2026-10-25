package usecases_test

import (
	"context"
	"testing"
	"time"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	usecases "github.com/Juancho43/SOPRO-backend/pkg/internal/use_cases"
	"github.com/Juancho43/SOPRO-backend/pkg/src/bus"
)

type mockHabitStreakLookupRepo struct {
	mockHabitRepo
	streaks map[string]int
}

func (m *mockHabitStreakLookupRepo) GetStreakByHabitID(ctx context.Context, habitID string) (int, int, *time.Time, error) {
	streak := m.streaks[habitID]
	return streak, streak, nil, nil
}

func TestEvaluateSynergyTrophies_StrictANDLogic(t *testing.T) {
	trophyID := "synergy-titan"
	userID := "user-1"
	habitA := "habit-code"
	habitB := "habit-gym"

	trophy := &entities.Trophy{
		ID:         trophyID,
		UserID:     userID,
		Name:       "Titán Integral",
		TargetType: entities.TargetSynergy,
	}

	trophyRepo := &mockTrophyRepo{
		trophies:    []*entities.Trophy{trophy},
		unlockedMap: make(map[string]bool),
		conditionsMap: map[string][]entities.TrophyCondition{
			trophyID: {
				{HabitID: habitA, RequiredStreakDays: 21},
				{HabitID: habitB, RequiredStreakDays: 21},
			},
		},
	}

	habitRepo := &mockHabitStreakLookupRepo{
		streaks: map[string]int{
			habitA: 21,
			habitB: 20, // Aún no cumple
		},
	}

	eventBus := bus.GetInstance()
	var unlockedEvent *entities.TrophyUnlockedEvent
	eventBus.Subscribe("trophy.unlocked", func(p interface{}) {
		if e, ok := p.(entities.TrophyUnlockedEvent); ok {
			unlockedEvent = &e
		}
	})

	uc := usecases.NewEvaluateSynergyTrophiesUseCase(trophyRepo, habitRepo, eventBus)

	t.Run("Cuando Hábito A tiene 21 y Hábito B tiene 20, NO se desbloquea", func(t *testing.T) {
		err := uc.Execute(context.Background(), userID, habitA)
		if err != nil {
			t.Fatalf("Error inesperado: %v", err)
		}

		if trophyRepo.unlockedMap[trophyID] {
			t.Errorf("El trofeo de sinergia no debía desbloquearse porque Hábito B no alcanza 21 días")
		}
	})

	t.Run("En el instante en que Hábito B alcanza 21, se desbloquea inmediatamente", func(t *testing.T) {
		habitRepo.streaks[habitB] = 21 // Ambos cumplen

		err := uc.Execute(context.Background(), userID, habitB)
		if err != nil {
			t.Fatalf("Error inesperado: %v", err)
		}

		if !trophyRepo.unlockedMap[trophyID] {
			t.Errorf("El trofeo de sinergia debió desbloquearse al cumplirse todas las condiciones AND")
		}

		time.Sleep(50 * time.Millisecond)
		if unlockedEvent == nil || !unlockedEvent.IsSynergy || unlockedEvent.TrophyID != trophyID {
			t.Errorf("Se esperaba evento trophy.unlocked con IsSynergy=true")
		}
	})
}
