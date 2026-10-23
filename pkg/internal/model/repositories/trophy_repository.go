package repositories

import (
	"context"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
)

// TrophyRepository define el contrato para trofeos, condiciones multihábito y desbloqueos.
type TrophyRepository interface {
	Create(ctx context.Context, trophy *entities.Trophy, conditions []entities.TrophyCondition) error
	GetByID(ctx context.Context, id string) (*entities.Trophy, error)
	ListByUserID(ctx context.Context, userID string) ([]*entities.Trophy, error)
	GetUserShowcase(ctx context.Context, userID string) ([]*entities.ShowcaseTrophy, error)
	UnlockTrophy(ctx context.Context, userID string, trophyID string) (bool, error)
	IsTrophyUnlocked(ctx context.Context, userID string, trophyID string) (bool, error)
	GetConditionsByTrophyID(ctx context.Context, trophyID string) ([]entities.TrophyCondition, error)
	FindSynergyTrophiesByHabit(ctx context.Context, habitID string) ([]*entities.Trophy, error)
	FindActiveTrophiesByTarget(ctx context.Context, userID string, targetType entities.TrophyTargetType, habitID *string) ([]*entities.Trophy, error)
}
