package services

import usecases "github.com/Juancho43/SOPRO-backend/pkg/internal/use_cases"

type RitualService struct {
	dailyRitual *usecases.DailyRitual
}

func NewRitualService(dailyRitual *usecases.DailyRitual) *RitualService {
	return &RitualService{dailyRitual: dailyRitual}
}

func (s *RitualService) ExecuteDailyRitual(command usecases.DailyRitualCommand) {

}

func (s *RitualService) ExecuteGetTodayRitual() {

}
