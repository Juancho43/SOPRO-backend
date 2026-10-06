package usecases

import (
	"errors"
	"time"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/repositories"
)

// MarkHabitAsDone orquesta la victoria diaria sobre un hábito específico.
type MarkHabitAsDone struct {
	streakRepo repositories.HabitStreakRepository
}

func NewMarkHabitAsDone(repo repositories.HabitStreakRepository) *MarkHabitAsDone {
	return &MarkHabitAsDone{
		streakRepo: repo,
	}
}

// Execute registra la ejecución del hábito, incrementa la inercia y notifica al sistema.
func (uc *MarkHabitAsDone) Execute(userID string, habitName string, now time.Time) error {
	// 1. Claridad de propósito.
	if habitName == "" {
		return errors.New("la claridad es esencial: el nombre del hábito no puede estar vacío")
	}

	// 2. Evaluar el estado actual de la disciplina.
	streak, err := uc.streakRepo.GetByUserIDAndHabit(userID, habitName)
	if err != nil {
		return err
	}

	if streak == nil {
		// El inicio de un nuevo estándar. Se establece la marca en 1.
		streak = &entities.HabitStreak{
			UserID:        userID,
			HabitName:     habitName,
			Frequency:     "diario",
			CurrentStreak: 1,
			MaxStreak:     1,
			CreatedAt:     now,
			UpdatedAt:     now,
		}
	} else {
		// Mantener el fuego vivo. Se acumula la inercia.
		streak.CurrentStreak++

		// Romper los propios límites: si la racha actual supera la histórica, se actualiza el récord.
		if streak.CurrentStreak > streak.MaxStreak {
			streak.MaxStreak = streak.CurrentStreak
		}
		streak.UpdatedAt = now
	}

	// 3. Consolidar el avance en la base de datos.
	if err := uc.streakRepo.Save(streak); err != nil {
		return err
	}

	return nil
}
