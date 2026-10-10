package usecases

import (
	"time"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/repositories"
)

type GetRitualByDate struct {
	repo repositories.RitualRepository
}

func NewGetRitualByDate(repo repositories.RitualRepository) *GetRitualByDate {
	return &GetRitualByDate{repo: repo}
}

func (u *GetRitualByDate) Execute(userID string, date string) (*entities.Ritual, error) {
	ritualDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}

	return u.repo.GetTodays(userID, ritualDate)
}
