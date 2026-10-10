package services

import (
	"log/slog"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	usecases "github.com/Juancho43/SOPRO-backend/pkg/internal/use_cases"
	"github.com/Juancho43/SOPRO-backend/pkg/src/bus"
)

type RitualService struct {
	dailyRitual     *usecases.DailyRitual
	todaysRitual    *usecases.TodaysRitual
	getRitualByDate *usecases.GetRitualByDate
	logger          *slog.Logger
	eventBus        bus.EventBus
}

func NewRitualService(
	dailyRitual *usecases.DailyRitual,
	todaysRitual *usecases.TodaysRitual,
	getRitualByDate *usecases.GetRitualByDate,
	logger *slog.Logger,
	bus bus.EventBus,
) *RitualService {
	return &RitualService{
		dailyRitual:     dailyRitual,
		todaysRitual:    todaysRitual,
		getRitualByDate: getRitualByDate,
		logger:          logger,
		eventBus:        bus,
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

	s.logger.Info("Ritual de hoy obtenido con éxito", slog.String("user_id", userID))

	return ritual, nil
}

func (s *RitualService) ExecuteGetRitualByDate(userID string, date string) (*entities.Ritual, error) {
	s.logger.Info("Consultando ritual por fecha", slog.String("user_id", userID), slog.String("date", date))

	ritual, err := s.getRitualByDate.Execute(userID, date)
	if err != nil {
		s.logger.Error("Error al obtener el ritual por fecha",
			slog.String("error", err.Error()),
			slog.String("user_id", userID),
			slog.String("date", date),
		)
		return nil, err
	}

	return ritual, nil
}
