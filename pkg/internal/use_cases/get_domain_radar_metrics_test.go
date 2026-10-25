package usecases_test

import (
	"context"
	"testing"
	"time"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	usecases "github.com/Juancho43/SOPRO-backend/pkg/internal/use_cases"
)

type mockDomainRepo struct {
	points []*entities.DomainEffortPoint
}

func (m *mockDomainRepo) Create(ctx context.Context, domain *entities.Domain) error { return nil }
func (m *mockDomainRepo) GetByID(ctx context.Context, id string) (*entities.Domain, error) {
	return nil, nil
}
func (m *mockDomainRepo) ListByUserID(ctx context.Context, userID string) ([]*entities.DomainWithHabitCount, error) {
	return nil, nil
}
func (m *mockDomainRepo) Update(ctx context.Context, domain *entities.Domain) error { return nil }
func (m *mockDomainRepo) Delete(ctx context.Context, id string, userID string) error { return nil }
func (m *mockDomainRepo) GetWeeklyEffortByDomain(ctx context.Context, userID string, fromDate time.Time) ([]*entities.DomainEffortPoint, error) {
	return m.points, nil
}

func TestGetDomainRadarMetrics_Calculation(t *testing.T) {
	now := time.Now()
	repo := &mockDomainRepo{
		points: []*entities.DomainEffortPoint{
			{DomainID: "d1", DomainName: "Arquitectura Go / TDD", CompletedCount: 6},
			{DomainID: "d2", DomainName: "Acondicionamiento Físico", CompletedCount: 5},
			{DomainID: "d3", DomainName: "Idiomas", CompletedCount: 2},
			{DomainID: "d4", DomainName: "Finanzas", CompletedCount: 0},
		},
	}

	uc := usecases.NewGetDomainRadarMetricsUseCase(repo)
	res, err := uc.Execute(context.Background(), "user-1", now)
	if err != nil {
		t.Fatalf("Error inesperado: %v", err)
	}

	expectedScores := map[string]float64{
		"d1": 85.7,
		"d2": 71.4,
		"d3": 28.6,
		"d4": 0.0,
	}

	for _, p := range res {
		expected, exists := expectedScores[p.DomainID]
		if !exists {
			t.Errorf("Dominio inesperado: %s", p.DomainID)
		}
		if p.EffortScore != expected {
			t.Errorf("Para dominio %s se esperaba score %f, se obtuvo %f", p.DomainID, expected, p.EffortScore)
		}
	}
}
