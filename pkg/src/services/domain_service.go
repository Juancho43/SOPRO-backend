package services

import (
	"context"
	"log/slog"
	"time"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	usecases "github.com/Juancho43/SOPRO-backend/pkg/internal/use_cases"
)

type DomainService struct {
	createDomain *usecases.CreateDomainUseCase
	listDomains  *usecases.ListDomainsUseCase
	updateDomain *usecases.UpdateDomainUseCase
	deleteDomain *usecases.DeleteDomainUseCase
	radarMetrics *usecases.GetDomainRadarMetricsUseCase
	logger       *slog.Logger
}

func NewDomainService(
	createDomain *usecases.CreateDomainUseCase,
	listDomains *usecases.ListDomainsUseCase,
	updateDomain *usecases.UpdateDomainUseCase,
	deleteDomain *usecases.DeleteDomainUseCase,
	radarMetrics *usecases.GetDomainRadarMetricsUseCase,
	logger *slog.Logger,
) *DomainService {
	return &DomainService{
		createDomain: createDomain,
		listDomains:  listDomains,
		updateDomain: updateDomain,
		deleteDomain: deleteDomain,
		radarMetrics: radarMetrics,
		logger:       logger,
	}
}

func (s *DomainService) CreateDomain(ctx context.Context, cmd usecases.CreateDomainCommand) (*entities.Domain, error) {
	s.logger.Info("Creando nuevo Dominio de Excelencia",
		slog.String("user_id", cmd.UserID),
		slog.String("name", cmd.Name),
	)
	domain, err := s.createDomain.Execute(ctx, cmd)
	if err != nil {
		s.logger.Error("Fallo al crear dominio", slog.String("error", err.Error()))
		return nil, err
	}
	return domain, nil
}

func (s *DomainService) ListDomains(ctx context.Context, userID string) ([]*entities.DomainWithHabitCount, error) {
	return s.listDomains.Execute(ctx, userID)
}

func (s *DomainService) UpdateDomain(ctx context.Context, cmd usecases.UpdateDomainCommand) error {
	s.logger.Info("Actualizando Dominio de Excelencia",
		slog.String("id", cmd.ID),
		slog.String("user_id", cmd.UserID),
	)
	return s.updateDomain.Execute(ctx, cmd)
}

func (s *DomainService) DeleteDomain(ctx context.Context, id string, userID string) error {
	s.logger.Info("Eliminando Dominio de Excelencia",
		slog.String("id", id),
		slog.String("user_id", userID),
	)
	return s.deleteDomain.Execute(ctx, id, userID)
}

func (s *DomainService) GetRadarMetrics(ctx context.Context, userID string, now time.Time) ([]*entities.DomainEffortPoint, error) {
	s.logger.Info("Consultando métricas de radar de dominios", slog.String("user_id", userID))
	return s.radarMetrics.Execute(ctx, userID, now)
}
