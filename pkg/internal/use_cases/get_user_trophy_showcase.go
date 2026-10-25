package usecases

import (
	"context"
	"errors"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/repositories"
)

type GetUserTrophyShowcaseUseCase struct {
	trophyRepo repositories.TrophyRepository
}

func NewGetUserTrophyShowcaseUseCase(repo repositories.TrophyRepository) *GetUserTrophyShowcaseUseCase {
	return &GetUserTrophyShowcaseUseCase{trophyRepo: repo}
}

func (uc *GetUserTrophyShowcaseUseCase) Execute(ctx context.Context, userID string) ([]*entities.ShowcaseTrophy, error) {
	if userID == "" {
		return nil, errors.New("el userID es obligatorio")
	}
	return uc.trophyRepo.GetUserShowcase(ctx, userID)
}

type ListTrophiesUseCase struct {
	trophyRepo repositories.TrophyRepository
}

func NewListTrophiesUseCase(repo repositories.TrophyRepository) *ListTrophiesUseCase {
	return &ListTrophiesUseCase{trophyRepo: repo}
}

func (uc *ListTrophiesUseCase) Execute(ctx context.Context, userID string) ([]*entities.Trophy, error) {
	if userID == "" {
		return nil, errors.New("el userID es obligatorio")
	}
	return uc.trophyRepo.ListByUserID(ctx, userID)
}
