package usecases

import (
	"errors"
	"strings"
	"time"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/repositories"
)

type DailyRitual struct {
	repo     repositories.RitualRepository
	location *time.Location
}

func NewDailyRitual(repo repositories.RitualRepository, timezoneStr string) (*DailyRitual, error) {
	loc, err := time.LoadLocation(timezoneStr)
	if err != nil {
		return nil, errors.New("error al cargar la zona horaria: " + err.Error())
	}

	return &DailyRitual{
		repo:     repo,
		location: loc,
	}, nil
}

type DailyRitualCommand struct {
	UserID      string   `json:"user_id" example:"f47ac10b-58cc-4372-a567-0e02b2c3d479"`
	GratefulFor string   `json:"GratefulFor" example:"Estoy agradecido por la oportunidad de crear sistemas increíbles hoy"`
	Goals       []string `json:"Goals" example:"Finalizar la integración de pagos,Leer 20 páginas de un libro,Hacer 45 minutos de ejercicio"`
}

func (u *DailyRitual) CreateDailyRitual(command DailyRitualCommand) error {
	if len(command.Goals) != 10 {
		return errors.New("la excelencia requiere precisión: debes definir exactamente 10 metas")
	}
	for _, goal := range command.Goals {
		if strings.TrimSpace(goal) == "" {
			return errors.New("la claridad es el punto de partida del éxito: ninguna de tus 10 metas puede estar vacía")
		}
	}
	if command.GratefulFor == "" {
		return errors.New("la gratitud es el punto de partida: el agradecimiento no puede estar vacío")
	}
	now := time.Now().In(u.location)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, u.location)

	exists, err := u.repo.ExistsForDate(command.UserID, today)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("ya has forjado tu disciplina hoy; ahora enfócate en ejecutar")
	}

	ritual := &entities.Ritual{
		UserID:         command.UserID,
		Date:           today,
		Agradecimiento: command.GratefulFor,
		Metas:          command.Goals,
		CreatedAt:      now,
	}

	if err := u.repo.Save(ritual); err != nil {
		return err
	}

	return nil
}
