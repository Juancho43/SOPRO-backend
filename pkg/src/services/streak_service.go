package services

import (
	"log/slog"
	"time"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	usecases "github.com/Juancho43/SOPRO-backend/pkg/internal/use_cases"
)

type StreakService struct {
	markHabitDone  *usecases.MarkHabitAsDone
	getHabitStreak *usecases.GetHabitStreak
	logger         *slog.Logger
}

// NewStreakService consolida tu centro de mando con todas sus dependencias.
func NewStreakService(
	markHabitDone *usecases.MarkHabitAsDone,
	getHabitStreak *usecases.GetHabitStreak,
	logger *slog.Logger,
) *StreakService {
	return &StreakService{
		markHabitDone:  markHabitDone,
		getHabitStreak: getHabitStreak,
		logger:         logger,
	}
}

// OnRitualCompleted reacciona al evento del sistema y forja la disciplina diaria.
func (s *StreakService) OnRitualCompleted(payload interface{}) {
	command, ok := payload.(usecases.DailyRitualCommand)
	if !ok {
		s.logger.Error("Falla de alineación: el payload recibido no es un DailyRitualCommand válido")
		return
	}

	err := s.markHabitDone.Execute(command.UserID, "Ritual Diario", time.Now())
	if err != nil {
		s.logger.Error("Obstáculo al incrementar la racha",
			slog.String("error", err.Error()),
			slog.String("user_id", command.UserID),
		)
		return
	}

	s.logger.Info("¡Inercia imparable! Racha maestra actualizada", slog.String("user_id", command.UserID))
}

// GetStreak recupera el estado de tu disciplina para que puedas medir tu progreso.
func (s *StreakService) GetStreak(userID string, habitName string) (*entities.HabitStreak, error) {
	s.logger.Info("Evaluando nivel de disciplina",
		slog.String("user_id", userID),
		slog.String("habit_name", habitName),
	)

	// Delegamos la lógica estricta al especialista
	streak, err := s.getHabitStreak.Execute(userID, habitName)
	if err != nil {
		s.logger.Error("Fallo al consultar la racha",
			slog.String("error", err.Error()),
			slog.String("user_id", userID),
			slog.String("habit_name", habitName),
		)
		return nil, err
	}

	s.logger.Info("Disciplina evaluada con éxito",
		slog.String("user_id", userID),
		slog.String("habit_name", habitName),
		slog.Int("current_streak", streak.CurrentStreak),
		slog.Int("max_streak", streak.MaxStreak),
	)

	return streak, nil
}
