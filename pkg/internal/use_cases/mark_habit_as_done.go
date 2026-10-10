package usecases

import (
	"errors"
	"time"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/repositories"
)

type MarkHabitAsDone struct {
	streakRepo repositories.HabitStreakRepository
	location   *time.Location
}

// Inyectamos el timezone string exacto (ej. "America/Buenos_Aires") para mantener la consistencia
func NewMarkHabitAsDone(repo repositories.HabitStreakRepository, timezoneStr string) (*MarkHabitAsDone, error) {
	loc, err := time.LoadLocation(timezoneStr)
	if err != nil {
		return nil, errors.New("error al cargar la zona horaria: " + err.Error())
	}

	return &MarkHabitAsDone{
		streakRepo: repo,
		location:   loc,
	}, nil
}

func (uc *MarkHabitAsDone) Execute(userID string, habitName string, now time.Time) error {
	if habitName == "" {
		return errors.New("la claridad es esencial: el nombre del hábito no puede estar vacío")
	}

	// 1. Alinear el tiempo de ejecución a la zona horaria de tu comando central
	localNow := now.In(uc.location)

	streak, err := uc.streakRepo.GetByUserIDAndHabit(userID, habitName)
	if err != nil {
		return err
	}

	if streak == nil {
		streak = &entities.HabitStreak{
			UserID:        userID,
			HabitName:     habitName,
			Frequency:     entities.Daily,
			CurrentStreak: 1,
			MaxStreak:     1,
			CreatedAt:     localNow,
			UpdatedAt:     localNow,
		}
	} else {
		// 2. Alinear la última actualización a la misma zona horaria antes de evaluar
		lastUpdate := streak.UpdatedAt.In(uc.location)

		// 3. Evaluar el paso del tiempo usando estrictamente la ubicación definida
		lastDate := time.Date(lastUpdate.Year(), lastUpdate.Month(), lastUpdate.Day(), 0, 0, 0, 0, uc.location)
		todayDate := time.Date(localNow.Year(), localNow.Month(), localNow.Day(), 0, 0, 0, 0, uc.location)

		daysDifference := int(todayDate.Sub(lastDate).Hours() / 24)

		if daysDifference == 0 {
			return errors.New("ya has forjado tu disciplina hoy; ahora enfócate en ejecutar")
		} else if daysDifference == 1 {
			streak.CurrentStreak++
			if streak.CurrentStreak > streak.MaxStreak {
				streak.MaxStreak = streak.CurrentStreak
			}
		} else {
			streak.CurrentStreak = 1
		}

		streak.UpdatedAt = localNow
	}

	if err := uc.streakRepo.Save(streak); err != nil {
		return err
	}

	return nil
}
