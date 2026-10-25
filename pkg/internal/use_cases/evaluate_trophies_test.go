package usecases_test

import (
	"context"
	"testing"
	"time"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	usecases "github.com/Juancho43/SOPRO-backend/pkg/internal/use_cases"
	"github.com/Juancho43/SOPRO-backend/pkg/src/bus"
)

type mockTrophyRepo struct {
	trophies        []*entities.Trophy
	unlockedMap     map[string]bool
	conditionsMap   map[string][]entities.TrophyCondition
	unlockedTrophies []string
}

func (m *mockTrophyRepo) Create(ctx context.Context, trophy *entities.Trophy, conditions []entities.TrophyCondition) error {
	m.trophies = append(m.trophies, trophy)
	if len(conditions) > 0 {
		if m.conditionsMap == nil {
			m.conditionsMap = make(map[string][]entities.TrophyCondition)
		}
		m.conditionsMap[trophy.ID] = conditions
	}
	return nil
}

func (m *mockTrophyRepo) GetByID(ctx context.Context, id string) (*entities.Trophy, error) {
	for _, t := range m.trophies {
		if t.ID == id {
			return t, nil
		}
	}
	return nil, nil
}

func (m *mockTrophyRepo) ListByUserID(ctx context.Context, userID string) ([]*entities.Trophy, error) {
	return m.trophies, nil
}

func (m *mockTrophyRepo) GetUserShowcase(ctx context.Context, userID string) ([]*entities.ShowcaseTrophy, error) {
	return nil, nil
}

func (m *mockTrophyRepo) UnlockTrophy(ctx context.Context, userID string, trophyID string) (bool, error) {
	if m.unlockedMap == nil {
		m.unlockedMap = make(map[string]bool)
	}
	if m.unlockedMap[trophyID] {
		return false, nil // Ya estaba desbloqueado (idempotencia)
	}
	m.unlockedMap[trophyID] = true
	m.unlockedTrophies = append(m.unlockedTrophies, trophyID)
	return true, nil
}

func (m *mockTrophyRepo) IsTrophyUnlocked(ctx context.Context, userID string, trophyID string) (bool, error) {
	if m.unlockedMap == nil {
		return false, nil
	}
	return m.unlockedMap[trophyID], nil
}

func (m *mockTrophyRepo) GetConditionsByTrophyID(ctx context.Context, trophyID string) ([]entities.TrophyCondition, error) {
	if m.conditionsMap == nil {
		return nil, nil
	}
	return m.conditionsMap[trophyID], nil
}

func (m *mockTrophyRepo) FindSynergyTrophiesByHabit(ctx context.Context, habitID string) ([]*entities.Trophy, error) {
	var result []*entities.Trophy
	for _, t := range m.trophies {
		if t.TargetType == entities.TargetSynergy {
			for _, cond := range m.conditionsMap[t.ID] {
				if cond.HabitID == habitID {
					result = append(result, t)
					break
				}
			}
		}
	}
	return result, nil
}

func (m *mockTrophyRepo) FindActiveTrophiesByTarget(
	ctx context.Context,
	userID string,
	targetType entities.TrophyTargetType,
	habitID *string,
) ([]*entities.Trophy, error) {
	var result []*entities.Trophy
	for _, t := range m.trophies {
		if t.TargetType == targetType && !m.unlockedMap[t.ID] {
			if habitID == nil || (t.HabitID != nil && *t.HabitID == *habitID) {
				result = append(result, t)
			}
		}
	}
	return result, nil
}

func TestEvaluateTrophies_SimpleTrophies(t *testing.T) {
	habitID := "habit-1"
	trophy := &entities.Trophy{
		ID:                 "trophy-7d",
		UserID:             "user-1",
		Name:               "Guerrero de la Constancia",
		TargetType:         entities.TargetHabit,
		HabitID:            &habitID,
		RequiredStreakDays: 7,
	}

	repo := &mockTrophyRepo{
		trophies:    []*entities.Trophy{trophy},
		unlockedMap: make(map[string]bool),
	}

	eventBus := bus.GetInstance()
	var unlockedEvent *entities.TrophyUnlockedEvent
	eventBus.Subscribe("trophy.unlocked", func(p interface{}) {
		if e, ok := p.(entities.TrophyUnlockedEvent); ok {
			unlockedEvent = &e
		}
	})

	uc := usecases.NewEvaluateTrophiesUseCase(repo, eventBus)

	t.Run("Racha de 6 días no desbloquea el trofeo de 7", func(t *testing.T) {
		err := uc.Execute(context.Background(), "user-1", entities.TargetHabit, &habitID, 6)
		if err != nil {
			t.Fatalf("Error inesperado: %v", err)
		}
		if repo.unlockedMap["trophy-7d"] {
			t.Errorf("El trofeo no debió desbloquearse con racha 6")
		}
	})

	t.Run("Racha de 7 días desbloquea el trofeo y emite evento", func(t *testing.T) {
		err := uc.Execute(context.Background(), "user-1", entities.TargetHabit, &habitID, 7)
		if err != nil {
			t.Fatalf("Error inesperado: %v", err)
		}
		if !repo.unlockedMap["trophy-7d"] {
			t.Errorf("El trofeo debió desbloquearse al alcanzar 7 días")
		}

		time.Sleep(50 * time.Millisecond)
		if unlockedEvent == nil || unlockedEvent.TrophyID != "trophy-7d" {
			t.Errorf("No se emitió el evento trophy.unlocked esperado")
		}
	})

	t.Run("Racha de 8 días es idempotente (no duplica desbloqueo)", func(t *testing.T) {
		unlockedEvent = nil
		err := uc.Execute(context.Background(), "user-1", entities.TargetHabit, &habitID, 8)
		if err != nil {
			t.Fatalf("Error inesperado: %v", err)
		}
		time.Sleep(50 * time.Millisecond)
		if unlockedEvent != nil {
			t.Errorf("No debió emitirse evento duplicado")
		}
	})
}
