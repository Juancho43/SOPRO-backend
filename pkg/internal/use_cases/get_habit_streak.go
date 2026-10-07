package usecases

import (
	"errors"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/repositories"
)

type GetHabitStreak struct {
	streakRepo repositories.HabitStreakRepository
}

func NewGetHabitStreak(repo repositories.HabitStreakRepository) *GetHabitStreak {
	return &GetHabitStreak{
		streakRepo: repo,
	}
}

// Execute recupera la inercia actual de tu hábito.
func (uc *GetHabitStreak) Execute(userID string, habitName string) (*entities.HabitStreak, error) {
	if habitName == "" {
		return nil, errors.New("la claridad es poder: debes especificar el nombre del hábito que deseas evaluar")
	}

	streak, err := uc.streakRepo.GetByUserIDAndHabit(userID, habitName)
	if err != nil {
		return nil, err
	}

	if streak == nil {
		return &entities.HabitStreak{
			UserID:        userID,
			HabitName:     habitName,
			Frequency:     entities.Daily,
			CurrentStreak: 0,
			MaxStreak:     0,
		}, nil
	}

	return streak, nil
}
