package services

import (
	"context"
	"log/slog"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	usecases "github.com/Juancho43/SOPRO-backend/pkg/internal/use_cases"
)

type FocusmateService struct {
	getAllFutureSessions *usecases.GetAllFutureSessions
	logger               *slog.Logger
}

func NewFocusmateService(
	getAllFutureSessions *usecases.GetAllFutureSessions,
	logger *slog.Logger,
) *FocusmateService {
	return &FocusmateService{
		getAllFutureSessions: getAllFutureSessions,
		logger:               logger,
	}
}

func (s *FocusmateService) ExecuteGetFutureSessions(ctx context.Context) ([]entities.FocusmateSession, error) {
	s.logger.Info("Solicitando sesiones futuras a Focusmate")

	sessions, err := s.getAllFutureSessions.Execute(ctx)
	if err != nil {
		s.logger.Error("Fallo al obtener sesiones", slog.String("error", err.Error()))
		return nil, err
	}

	return sessions, nil
}
