package usecases_test

import (
	"context"
	"testing"
	"time"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	usecases "github.com/Juancho43/SOPRO-backend/pkg/internal/use_cases"
	"github.com/Juancho43/SOPRO-backend/pkg/src/bus"
)

type mockHabitRepo struct {
	habit         *entities.Habit
	logs          map[string]bool
	currentStreak int
	maxStreak     int
	lastUpdated   *time.Time
	logSaved      *entities.HabitLog
	streakSaved   bool
}

func (m *mockHabitRepo) Create(ctx context.Context, habit *entities.Habit) error {
	m.habit = habit
	return nil
}

func (m *mockHabitRepo) GetByID(ctx context.Context, id string) (*entities.Habit, error) {
	if m.habit != nil && m.habit.ID == id {
		return m.habit, nil
	}
	return nil, nil
}

func (m *mockHabitRepo) ListByUserID(ctx context.Context, userID string, todayDate string) ([]*entities.HabitWithStreak, error) {
	return nil, nil
}

func (m *mockHabitRepo) Update(ctx context.Context, habit *entities.Habit) error {
	return nil
}

func (m *mockHabitRepo) Archive(ctx context.Context, id string, userID string) error {
	return nil
}

func (m *mockHabitRepo) LogCompletion(ctx context.Context, log *entities.HabitLog) error {
	m.logSaved = log
	if m.logs == nil {
		m.logs = make(map[string]bool)
	}
	m.logs[log.LogDate] = true
	return nil
}

func (m *mockHabitRepo) HasLogForDate(ctx context.Context, habitID string, date string) (bool, error) {
	if m.logs == nil {
		return false, nil
	}
	return m.logs[date], nil
}

func (m *mockHabitRepo) GetLogs(ctx context.Context, habitID string, fromDate string, toDate string) ([]string, error) {
	return nil, nil
}

func (m *mockHabitRepo) GetStreakByHabitID(ctx context.Context, habitID string) (int, int, *time.Time, error) {
	return m.currentStreak, m.maxStreak, m.lastUpdated, nil
}

func (m *mockHabitRepo) SaveHabitStreak(
	ctx context.Context,
	habitID string,
	userID string,
	habitName string,
	frequency string,
	currentStreak int,
	maxStreak int,
	now time.Time,
) error {
	m.currentStreak = currentStreak
	m.maxStreak = maxStreak
	m.lastUpdated = &now
	m.streakSaved = true
	return nil
}

func TestCheckInHabit_StreakCalculationAndEvents(t *testing.T) {
	loc, _ := time.LoadLocation("America/Buenos_Aires")
	today := time.Date(2026, 10, 17, 10, 0, 0, 0, loc)
	yesterday := today.AddDate(0, 0, -1)

	t.Run("Primer check-in inicializa la racha en 1", func(t *testing.T) {
		repo := &mockHabitRepo{
			habit: &entities.Habit{
				ID:        "h1",
				UserID:    "u1",
				Name:      "Programación",
				Frequency: "daily",
			},
		}

		eventBus := bus.GetInstance()
		var publishedEvent *entities.HabitCompletedEvent
		eventBus.Subscribe("habit.completed", func(p interface{}) {
			if e, ok := p.(entities.HabitCompletedEvent); ok {
				publishedEvent = &e
			}
		})

		uc, err := usecases.NewCheckInHabitUseCase(repo, eventBus, "America/Buenos_Aires")
		if err != nil {
			t.Fatalf("Error al instanciar caso de uso: %v", err)
		}

		res, err := uc.Execute(context.Background(), "h1", "u1", today)
		if err != nil {
			t.Fatalf("No se esperaba error, se obtuvo: %v", err)
		}

		if res.CurrentStreak != 1 || res.MaxStreak != 1 {
			t.Errorf("Se esperaba racha 1/1, se obtuvo %d/%d", res.CurrentStreak, res.MaxStreak)
		}

		// Esperar brevemente para verificar evento asíncrono
		time.Sleep(50 * time.Millisecond)
		if publishedEvent == nil || publishedEvent.StreakDays != 1 {
			t.Errorf("No se publicó el evento habit.completed esperado")
		}
	})

	t.Run("Rechaza un segundo check-in en el mismo día calendario", func(t *testing.T) {
		repo := &mockHabitRepo{
			habit: &entities.Habit{
				ID:        "h1",
				UserID:    "u1",
				Name:      "Programación",
				Frequency: "daily",
			},
			logs: map[string]bool{
				"2026-10-17": true,
			},
		}

		uc, _ := usecases.NewCheckInHabitUseCase(repo, nil, "America/Buenos_Aires")
		_, err := uc.Execute(context.Background(), "h1", "u1", today)
		if err != usecases.ErrAlreadyCheckedIn {
			t.Errorf("Se esperaba ErrAlreadyCheckedIn, se obtuvo %v", err)
		}
	})

	t.Run("Check-in consecutivo al día siguiente incrementa la racha", func(t *testing.T) {
		repo := &mockHabitRepo{
			habit: &entities.Habit{
				ID:        "h1",
				UserID:    "u1",
				Name:      "Programación",
				Frequency: "daily",
			},
			currentStreak: 5,
			maxStreak:     5,
			lastUpdated:   &yesterday,
		}

		uc, _ := usecases.NewCheckInHabitUseCase(repo, nil, "America/Buenos_Aires")
		res, err := uc.Execute(context.Background(), "h1", "u1", today)
		if err != nil {
			t.Fatalf("Error inesperado: %v", err)
		}

		if res.CurrentStreak != 6 || res.MaxStreak != 6 {
			t.Errorf("Se esperaba racha 6/6, se obtuvo %d/%d", res.CurrentStreak, res.MaxStreak)
		}
	})
}
