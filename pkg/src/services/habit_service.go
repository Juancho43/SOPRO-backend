package services

import (
	"context"
	"log/slog"
	"time"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	usecases "github.com/Juancho43/SOPRO-backend/pkg/internal/use_cases"
)

type HabitService struct {
	createHabit   *usecases.CreateHabitUseCase
	listHabits    *usecases.ListUserHabitsUseCase
	checkInHabit  *usecases.CheckInHabitUseCase
	habitHistory  *usecases.GetHabitHistoryUseCase
	logger        *slog.Logger
}

func NewHabitService(
	createHabit *usecases.CreateHabitUseCase,
	listHabits *usecases.ListUserHabitsUseCase,
	checkInHabit *usecases.CheckInHabitUseCase,
	habitHistory *usecases.GetHabitHistoryUseCase,
	logger *slog.Logger,
) *HabitService {
	return &HabitService{
		createHabit:  createHabit,
		listHabits:   listHabits,
		checkInHabit: checkInHabit,
		habitHistory: habitHistory,
		logger:       logger,
	}
}

func (s *HabitService) CreateHabit(ctx context.Context, cmd usecases.CreateHabitCommand) (*entities.Habit, error) {
	s.logger.Info("Añadiendo hábito al arsenal",
		slog.String("user_id", cmd.UserID),
		slog.String("name", cmd.Name),
		slog.String("frequency", cmd.Frequency),
	)
	habit, err := s.createHabit.Execute(ctx, cmd)
	if err != nil {
		s.logger.Error("Fallo al crear hábito", slog.String("error", err.Error()))
		return nil, err
	}
	return habit, nil
}

func (s *HabitService) ListUserHabits(ctx context.Context, userID string, now time.Time) ([]*entities.HabitWithStreak, error) {
	return s.listHabits.Execute(ctx, userID, now)
}

func (s *HabitService) CheckInHabit(ctx context.Context, habitID string, userID string, now time.Time) (*entities.HabitCheckInResult, error) {
	s.logger.Info("Registrando check-in de disciplina",
		slog.String("habit_id", habitID),
		slog.String("user_id", userID),
	)
	result, err := s.checkInHabit.Execute(ctx, habitID, userID, now)
	if err != nil {
		s.logger.Error("Obstáculo al registrar check-in",
			slog.String("habit_id", habitID),
			slog.String("error", err.Error()),
		)
		return nil, err
	}
	s.logger.Info("¡Hábito marcado con éxito!",
		slog.String("habit_id", habitID),
		slog.Int("streak", result.CurrentStreak),
	)
	return result, nil
}

func (s *HabitService) GetHabitHistory(ctx context.Context, habitID string, userID string, now time.Time) ([]string, error) {
	return s.habitHistory.Execute(ctx, habitID, userID, now)
}
