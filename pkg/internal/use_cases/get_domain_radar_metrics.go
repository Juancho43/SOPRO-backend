package usecases

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/repositories"
)

type GetDomainRadarMetricsUseCase struct {
	domainRepo repositories.DomainRepository
}

func NewGetDomainRadarMetricsUseCase(repo repositories.DomainRepository) *GetDomainRadarMetricsUseCase {
	return &GetDomainRadarMetricsUseCase{domainRepo: repo}
}

func (uc *GetDomainRadarMetricsUseCase) Execute(ctx context.Context, userID string, now time.Time) ([]*entities.DomainEffortPoint, error) {
	if userID == "" {
		return nil, errors.New("el userID es obligatorio")
	}

	// Ventana de los últimos 7 días
	fromDate := now.AddDate(0, 0, -6) // incluye hoy + 6 días previos = 7 días
	points, err := uc.domainRepo.GetWeeklyEffortByDomain(ctx, userID, fromDate)
	if err != nil {
		return nil, err
	}

	for _, p := range points {
		// Normalizado relativo a 7 días de la semana, máximo 100%
		rawScore := (float64(p.CompletedCount) / 7.0) * 100.0
		if rawScore > 100.0 {
			rawScore = 100.0
		}
		p.EffortScore = math.Round(rawScore*10) / 10
	}

	return points, nil
}
