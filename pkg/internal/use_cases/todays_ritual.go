package usecases

import (
	"time"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/repositories"
)

type TodaysRitual struct {
	repo     repositories.RitualRepository
	location *time.Location
}

func NewTodaysRitual(repo repositories.RitualRepository, timezoneStr string) (*TodaysRitual, error) {
	loc, err := time.LoadLocation(timezoneStr)
	if err != nil {
		return nil, err
	}
	return &TodaysRitual{
		repo:     repo,
		location: loc,
	}, nil
}

func (u *TodaysRitual) Execute(userID string) (bool, error) {
	now := time.Now().In(u.location)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, u.location)

	return u.repo.ExistsForDate(userID, today)
}
