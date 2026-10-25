package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/repositories"
)

type ListUserHabitsUseCase struct {
	habitRepo repositories.HabitRepository
	location  *time.Location
}

func NewListUserHabitsUseCase(repo repositories.HabitRepository, timezoneStr string) (*ListUserHabitsUseCase, error) {
	loc, err := time.LoadLocation(timezoneStr)
	if err != nil {
		loc = time.UTC
	}
	return &ListUserHabitsUseCase{
		habitRepo: repo,
		location:  loc,
	}, nil
}

func (uc *ListUserHabitsUseCase) Execute(ctx context.Context, userID string, now time.Time) ([]*entities.HabitWithStreak, error) {
	if userID == "" {
		return nil, errors.New("el userID es obligatorio")
	}

	todayDate := now.In(uc.location).Format("2006-01-02")
	return uc.habitRepo.ListByUserID(ctx, userID, todayDate)
}
