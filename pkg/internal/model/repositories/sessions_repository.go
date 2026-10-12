package repositories

import (
	"context"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
)

// FocusmateAPI es el contrato innegociable. Tu caso de uso solo conoce esta interfaz,
// no le importa si los datos vienen de HTTP, un mock o una base de datos.
type FocusmateAPI interface {
	FetchFutureSessions(ctx context.Context) ([]entities.FocusmateSession, error)
}
