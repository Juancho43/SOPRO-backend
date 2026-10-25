package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/repositories"
)

type GetHabitHistoryUseCase struct {
	habitRepo repositories.HabitRepository
}

func NewGetHabitHistoryUseCase(repo repositories.HabitRepository) *GetHabitHistoryUseCase {
	return &GetHabitHistoryUseCase{habitRepo: repo}
}

func (uc *GetHabitHistoryUseCase) Execute(ctx context.Context, habitID string, userID string, now time.Time) ([]string, error) {
	if habitID == "" || userID == "" {
		return nil, errors.New("habitID y userID son obligatorios")
	}

	habit, err := uc.habitRepo.GetByID(ctx, habitID)
	if err != nil {
		return nil, err
	}
	if habit == nil || habit.UserID != userID {
		return nil, ErrHabitNotFound
	}

	toDate := now.Format("2006-01-02")
	fromDate := now.AddDate(-1, 0, 0).Format("2006-01-02") // 365 días atrás

	logs, err := uc.habitRepo.GetLogs(ctx, habitID, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	if logs == nil {
		logs = []string{}
	}

	return logs, nil
}
