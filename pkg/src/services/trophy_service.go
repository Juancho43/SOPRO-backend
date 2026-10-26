package services

import (
	"context"
	"log/slog"
	"time"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	usecases "github.com/Juancho43/SOPRO-backend/pkg/internal/use_cases"
)

type TrophyService struct {
	createTrophy     *usecases.CreateTrophyUseCase
	showcase         *usecases.GetUserTrophyShowcaseUseCase
	listTrophies     *usecases.ListTrophiesUseCase
	evaluateTrophies *usecases.EvaluateTrophiesUseCase
	evaluateSynergy  *usecases.EvaluateSynergyTrophiesUseCase
	synergyStatus    *usecases.GetSynergyStatusUseCase
	streakService    *StreakService
	logger           *slog.Logger
}

func NewTrophyService(
	createTrophy *usecases.CreateTrophyUseCase,
	showcase *usecases.GetUserTrophyShowcaseUseCase,
	listTrophies *usecases.ListTrophiesUseCase,
	evaluateTrophies *usecases.EvaluateTrophiesUseCase,
	evaluateSynergy *usecases.EvaluateSynergyTrophiesUseCase,
	synergyStatus *usecases.GetSynergyStatusUseCase,
	streakService *StreakService,
	logger *slog.Logger,
) *TrophyService {
	return &TrophyService{
		createTrophy:     createTrophy,
		showcase:         showcase,
		listTrophies:     listTrophies,
		evaluateTrophies: evaluateTrophies,
		evaluateSynergy:  evaluateSynergy,
		synergyStatus:    synergyStatus,
		streakService:    streakService,
		logger:           logger,
	}
}

func (s *TrophyService) CreateTrophy(ctx context.Context, cmd usecases.CreateTrophyCommand) (*entities.Trophy, error) {
	s.logger.Info("Diseñando nuevo trofeo de victoria",
		slog.String("user_id", cmd.UserID),
		slog.String("name", cmd.Name),
		slog.String("target_type", cmd.TargetType),
	)
	trophy, err := s.createTrophy.Execute(ctx, cmd)
	if err != nil {
		s.logger.Error("Fallo al crear trofeo", slog.String("error", err.Error()))
		return nil, err
	}
	return trophy, nil
}

func (s *TrophyService) GetUserShowcase(ctx context.Context, userID string) ([]*entities.ShowcaseTrophy, error) {
	return s.showcase.Execute(ctx, userID)
}

func (s *TrophyService) ListTrophies(ctx context.Context, userID string) ([]*entities.Trophy, error) {
	return s.listTrophies.Execute(ctx, userID)
}

func (s *TrophyService) GetSynergyStatus(ctx context.Context, trophyID string, userID string) (*entities.SynergyStatusResponse, error) {
	return s.synergyStatus.Execute(ctx, trophyID, userID)
}

// OnHabitCompleted es el oyente asíncrono que evalúa desbloqueos ante cada check-in de hábito.
func (s *TrophyService) OnHabitCompleted(payload interface{}) {
	event, ok := payload.(entities.HabitCompletedEvent)
	if !ok {
		s.logger.Error("Payload inválido recibido en evento habit.completed")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Evaluar trofeos simples asociados directamente a este hábito
	if err := s.evaluateTrophies.Execute(ctx, event.UserID, entities.TargetHabit, &event.HabitID, event.StreakDays); err != nil {
		s.logger.Error("Error al evaluar trofeos de hábito", slog.String("error", err.Error()))
	}

	// 2. Evaluar trofeos de sinergia multihábito
	if err := s.evaluateSynergy.Execute(ctx, event.UserID, event.HabitID); err != nil {
		s.logger.Error("Error al evaluar trofeos de sinergia", slog.String("error", err.Error()))
	}
}

// OnRitualCompleted es el oyente asíncrono para evaluar trofeos de racha maestra tras el ritual diario.
func (s *TrophyService) OnRitualCompleted(payload interface{}) {
	cmd, ok := payload.(usecases.DailyRitualCommand)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Consultamos la racha maestra calculada
	if s.streakService != nil {
		streak, err := s.streakService.GetStreak(cmd.UserID, "Ritual Diario")
		if err == nil && streak != nil {
			if err := s.evaluateTrophies.Execute(ctx, cmd.UserID, entities.TargetMasterStreak, nil, streak.CurrentStreak); err != nil {
				s.logger.Error("Error al evaluar trofeos de racha maestra", slog.String("error", err.Error()))
			}
		}
	}
}
