package repositories

import (
	"time"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
)

type RitualRepository interface {
	Save(ritual *entities.Ritual) error
	ExistsForDate(userID string, date time.Time) (bool, error)
	GetTodays(userID string, date time.Time) (*entities.Ritual, error)
}
