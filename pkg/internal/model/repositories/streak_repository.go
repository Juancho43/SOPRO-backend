package repositories

import "github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"

// HabitStreakRepository establece las leyes inmutables para acceder a tus victorias.
type HabitStreakRepository interface {
	// Save crea un nuevo hábito o actualiza la racha (Current/Max) de uno existente
	// tras una victoria diaria.
	Save(streak *entities.HabitStreak) error

	// GetByUserIDAndHabit recupera una racha específica para evaluarla.
	// Esencial para que el motor de gamificación sepa si debe sumar +1 hoy o reiniciar a 0.
	GetByUserIDAndHabit(userID string, habitName string) (*entities.HabitStreak, error)

	// GetAllByUserID carga tu arsenal completo de hábitos.
	// Se utiliza para renderizar los mapas de calor y el progreso en tu Dashboard.
	GetAllByUserID(userID string) ([]*entities.HabitStreak, error)
}
