package repositories

import (
	"context"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
)

// ChallengeRepository define el contrato para persistir y consultar retos e hitos (Épica 3).
type ChallengeRepository interface {
	// Retos
	CreateChallenge(ctx context.Context, challenge *entities.Challenge) error
	GetChallengeByID(ctx context.Context, id string, userID string) (*entities.Challenge, error)
	ListChallengesByUserID(ctx context.Context, userID string, status *entities.ChallengeStatus) ([]*entities.Challenge, error)
	UpdateChallenge(ctx context.Context, challenge *entities.Challenge) error
	UpdateChallengeStatus(ctx context.Context, id string, userID string, status entities.ChallengeStatus) error

	// Hitos / Deadlines Intermedias
	CreateMilestone(ctx context.Context, milestone *entities.Milestone) error
	GetMilestoneByID(ctx context.Context, id string) (*entities.Milestone, error)
	ListMilestonesByChallengeID(ctx context.Context, challengeID string) ([]*entities.Milestone, error)
	UpdateMilestone(ctx context.Context, milestone *entities.Milestone) error
	DeleteMilestone(ctx context.Context, id string, challengeID string) error
	ListUpcomingMilestones(ctx context.Context, userID string, maxDaysAhead int, today string) ([]*entities.UpcomingDeadline, error)
}
