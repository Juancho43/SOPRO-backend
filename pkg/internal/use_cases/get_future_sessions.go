package usecases

import (
	"context"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/repositories"
)

type GetAllFutureSessions struct {
	api repositories.FocusmateAPI
}

func NewGetAllFutureSessions(api repositories.FocusmateAPI) *GetAllFutureSessions {
	return &GetAllFutureSessions{
		api: api,
	}
}

// Execute recupera tus próximos bloques de concentración.
func (uc *GetAllFutureSessions) Execute(ctx context.Context) ([]entities.FocusmateSession, error) {
	// Aquí podrías agregar lógica de negocio adicional si fuera necesario
	// (ej. filtrar sesiones canceladas, mapear zonas horarias, etc.)
	sessions, err := uc.api.FetchFutureSessions(ctx)
	if err != nil {
		return nil, err
	}

	return sessions, nil
}
