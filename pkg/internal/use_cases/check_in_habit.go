package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/repositories"
	"github.com/Juancho43/SOPRO-backend/pkg/src/bus"
)

var (
	ErrHabitNotFound       = errors.New("hábito no encontrado o no pertenece al usuario")
	ErrAlreadyCheckedIn    = errors.New("ya has registrado el cumplimiento de este hábito hoy")
)

type CheckInHabitUseCase struct {
	habitRepo repositories.HabitRepository
	eventBus  bus.EventBus
	location  *time.Location
}

func NewCheckInHabitUseCase(repo repositories.HabitRepository, eventBus bus.EventBus, timezoneStr string) (*CheckInHabitUseCase, error) {
	loc, err := time.LoadLocation(timezoneStr)
	if err != nil {
		loc = time.UTC
	}
	return &CheckInHabitUseCase{
		habitRepo: repo,
		eventBus:  eventBus,
		location:  loc,
	}, nil
}

func (uc *CheckInHabitUseCase) Execute(ctx context.Context, habitID string, userID string, now time.Time) (*entities.HabitCheckInResult, error) {
	if habitID == "" || userID == "" {
		return nil, errors.New("habitID y userID son obligatorios")
	}

	habit, err := uc.habitRepo.GetByID(ctx, habitID)
	if err != nil {
		return nil, err
	}
	if habit == nil || habit.UserID != userID || habit.Archived {
		return nil, ErrHabitNotFound
	}

	localNow := now.In(uc.location)
	todayStr := localNow.Format("2006-01-02")

	// 1. Validar que no exista log para la fecha local
	alreadyChecked, err := uc.habitRepo.HasLogForDate(ctx, habitID, todayStr)
	if err != nil {
		return nil, err
	}
	if alreadyChecked {
		return nil, ErrAlreadyCheckedIn
	}

	// 2. Registrar en habit_logs
	log := &entities.HabitLog{
		HabitID:     habitID,
		UserID:      userID,
		LogDate:     todayStr,
		CompletedAt: localNow,
	}
	if err := uc.habitRepo.LogCompletion(ctx, log); err != nil {
		return nil, err
	}

	// 3. Recalcular racha en habit_streaks
	currentStreak, maxStreak, lastUpdated, err := uc.habitRepo.GetStreakByHabitID(ctx, habitID)
	if err != nil {
		return nil, err
	}

	if lastUpdated == nil {
		currentStreak = 1
		maxStreak = 1
	} else {
		lastLocal := lastUpdated.In(uc.location)
		lastDate := time.Date(lastLocal.Year(), lastLocal.Month(), lastLocal.Day(), 0, 0, 0, 0, uc.location)
		todayDate := time.Date(localNow.Year(), localNow.Month(), localNow.Day(), 0, 0, 0, 0, uc.location)

		daysDiff := int(todayDate.Sub(lastDate).Hours() / 24)
		if daysDiff == 1 {
			currentStreak++
			if currentStreak > maxStreak {
				maxStreak = currentStreak
			}
		} else if daysDiff > 1 {
			currentStreak = 1
		}
	}

	if err := uc.habitRepo.SaveHabitStreak(ctx, habitID, userID, habit.Name, habit.Frequency, currentStreak, maxStreak, localNow); err != nil {
		return nil, err
	}

	// 4. Emitir asíncronamente evento habit.completed en el EventBus
	if uc.eventBus != nil {
		event := entities.HabitCompletedEvent{
			UserID:     userID,
			HabitID:    habitID,
			HabitName:  habit.Name,
			StreakDays: currentStreak,
			Date:       todayStr,
		}
		uc.eventBus.Publish("habit.completed", event)
	}

	return &entities.HabitCheckInResult{
		HabitID:       habitID,
		LogDate:       todayStr,
		CurrentStreak: currentStreak,
		MaxStreak:     maxStreak,
		Message:       "Hábito cumplido. Tu constancia forja tu destino.",
	}, nil
}
