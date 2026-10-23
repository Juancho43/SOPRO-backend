package repositories

import (
	"context"
	"time"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
)

// DomainRepository define las operaciones de persistencia para los Dominios de Excelencia.
type DomainRepository interface {
	Create(ctx context.Context, domain *entities.Domain) error
	GetByID(ctx context.Context, id string) (*entities.Domain, error)
	ListByUserID(ctx context.Context, userID string) ([]*entities.DomainWithHabitCount, error)
	Update(ctx context.Context, domain *entities.Domain) error
	Delete(ctx context.Context, id string, userID string) error
	GetWeeklyEffortByDomain(ctx context.Context, userID string, fromDate time.Time) ([]*entities.DomainEffortPoint, error)
}
