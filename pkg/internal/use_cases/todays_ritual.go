package usecases

import "github.com/Juancho43/SOPRO-backend/pkg/internal/model/repositories"

type TodaysRitual struct {
	repo repositories.RitualRepository
}

func NewTodaysRitual(repo repositories.RitualRepository) *TodaysRitual {
	return &TodaysRitual{repo: repo}
}

func (u *TodaysRitual) Execute() {

}
