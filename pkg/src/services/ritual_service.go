package services

import (
	"log/slog"

	usecases "github.com/Juancho43/SOPRO-backend/pkg/internal/use_cases"
	"github.com/Juancho43/SOPRO-backend/pkg/src/bus"
)

type RitualService struct {
	dailyRitual  *usecases.DailyRitual
	todaysRitual *usecases.TodaysRitual
	logger       *slog.Logger
	eventBus     bus.EventBus
}

func NewRitualService(
	dailyRitual *usecases.DailyRitual,
	todaysRitual *usecases.TodaysRitual,
	logger *slog.Logger,
	bus bus.EventBus,
) *RitualService {
	return &RitualService{
		dailyRitual:  dailyRitual,
		todaysRitual: todaysRitual,
		logger:       logger,
		eventBus:     bus,
	}
}

func (s *RitualService) ExecuteDailyRitual(command usecases.DailyRitualCommand) error {
	s.logger.Info("Iniciando creación de ritual diario", slog.String("user_id", command.UserID))

	err := s.dailyRitual.CreateDailyRitual(command)
	if err != nil {
		s.logger.Error("Fallo al forjar la disciplina",
			slog.String("error", err.Error()),
			slog.String("user_id", command.UserID),
		)
		return err
	}

	s.logger.Info("Ritual diario creado con éxito", slog.String("user_id", command.UserID))
	s.eventBus.Publish("ritual.completed", command)
	return nil
}

func (s *RitualService) ExecuteGetTodayRitual(userID string) (bool, error) {
	s.logger.Info("Consultando ritual de hoy", slog.String("user_id", userID))

	ritual, err := s.todaysRitual.Execute(userID)
	if err != nil {
		s.logger.Error("Error al obtener el ritual de hoy",
			slog.String("error", err.Error()),
			slog.String("user_id", userID),
		)
		return false, err
	}

	return ritual, nil
}
