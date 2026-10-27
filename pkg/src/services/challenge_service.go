package services

import (
	"context"
	"log/slog"
	"time"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	usecases "github.com/Juancho43/SOPRO-backend/pkg/internal/use_cases"
)

type ChallengeService struct {
	createChallenge *usecases.CreateChallengeUseCase
	manageChallenge *usecases.ManageChallengesUseCase
	manageMilestone *usecases.ManageMilestonesUseCase
	logger          *slog.Logger
}

func NewChallengeService(
	createChallenge *usecases.CreateChallengeUseCase,
	manageChallenge *usecases.ManageChallengesUseCase,
	manageMilestone *usecases.ManageMilestonesUseCase,
	logger *slog.Logger,
) *ChallengeService {
	return &ChallengeService{
		createChallenge: createChallenge,
		manageChallenge: manageChallenge,
		manageMilestone: manageMilestone,
		logger:          logger,
	}
}

// Retos
func (s *ChallengeService) CreateChallenge(ctx context.Context, cmd usecases.CreateChallengeCommand) (*entities.Challenge, error) {
	s.logger.Info("Declarando nuevo reto estratégico",
		slog.String("user_id", cmd.UserID),
		slog.String("title", cmd.Title),
	)
	challenge, err := s.createChallenge.Execute(ctx, cmd)
	if err != nil {
		s.logger.Error("Fallo al crear reto", slog.String("error", err.Error()))
		return nil, err
	}
	s.logger.Info("Reto estratégico establecido con éxito", slog.String("id", challenge.ID))
	return challenge, nil
}

func (s *ChallengeService) ListChallenges(ctx context.Context, userID string, status *string) ([]*entities.Challenge, error) {
	return s.manageChallenge.List(ctx, userID, status)
}

func (s *ChallengeService) GetChallengeDetail(ctx context.Context, id string, userID string, now time.Time) (*entities.ChallengeDetail, error) {
	return s.manageChallenge.GetDetail(ctx, id, userID, now)
}

func (s *ChallengeService) UpdateChallenge(ctx context.Context, cmd usecases.UpdateChallengeCommand) (*entities.Challenge, error) {
	s.logger.Info("Actualizando reto estratégico",
		slog.String("user_id", cmd.UserID),
		slog.String("id", cmd.ID),
	)
	return s.manageChallenge.Update(ctx, cmd)
}

func (s *ChallengeService) UpdateChallengeStatus(ctx context.Context, id string, userID string, status entities.ChallengeStatus) error {
	s.logger.Info("Actualizando estado de reto",
		slog.String("id", id),
		slog.String("user_id", userID),
		slog.String("status", string(status)),
	)
	return s.manageChallenge.UpdateStatus(ctx, id, userID, status)
}

// Hitos
func (s *ChallengeService) CreateMilestone(ctx context.Context, cmd usecases.CreateMilestoneCommand) (*entities.Milestone, error) {
	s.logger.Info("Añadiendo hito intermedio a reto",
		slog.String("user_id", cmd.UserID),
		slog.String("challenge_id", cmd.ChallengeID),
		slog.String("title", cmd.Title),
	)
	return s.manageMilestone.Create(ctx, cmd)
}

func (s *ChallengeService) UpdateMilestone(ctx context.Context, cmd usecases.UpdateMilestoneCommand) (*entities.Milestone, error) {
	s.logger.Info("Actualizando hito intermedio",
		slog.String("id", cmd.ID),
		slog.String("challenge_id", cmd.ChallengeID),
	)
	return s.manageMilestone.Update(ctx, cmd)
}

func (s *ChallengeService) DeleteMilestone(ctx context.Context, milestoneID string, challengeID string, userID string) error {
	s.logger.Info("Eliminando hito",
		slog.String("id", milestoneID),
		slog.String("challenge_id", challengeID),
	)
	return s.manageMilestone.Delete(ctx, milestoneID, challengeID, userID)
}

func (s *ChallengeService) GetUpcomingDeadlines(ctx context.Context, userID string, daysAhead int, now time.Time) ([]*entities.UpcomingDeadline, error) {
	return s.manageMilestone.ListUpcoming(ctx, userID, daysAhead, now)
}
