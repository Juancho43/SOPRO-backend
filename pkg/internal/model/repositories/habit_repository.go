package repositories

import (
	"context"
	"time"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
)

// HabitRepository define el contrato para almacenar y consultar hábitos, logs y sus rachas.
type HabitRepository interface {
	Create(ctx context.Context, habit *entities.Habit) error
	GetByID(ctx context.Context, id string) (*entities.Habit, error)
	ListByUserID(ctx context.Context, userID string, todayDate string) ([]*entities.HabitWithStreak, error)
	Update(ctx context.Context, habit *entities.Habit) error
	Archive(ctx context.Context, id string, userID string) error
	LogCompletion(ctx context.Context, log *entities.HabitLog) error
	HasLogForDate(ctx context.Context, habitID string, date string) (bool, error)
	GetLogs(ctx context.Context, habitID string, fromDate string, toDate string) ([]string, error)
	GetStreakByHabitID(ctx context.Context, habitID string) (currentStreak int, maxStreak int, lastUpdated *time.Time, err error)
	SaveHabitStreak(ctx context.Context, habitID string, userID string, habitName string, frequency string, currentStreak int, maxStreak int, now time.Time) error
}
