package injectors

import (
	"github.com/Juancho43/SOPRO-backend/pkg/src/bus"
	"github.com/Juancho43/SOPRO-backend/pkg/src/services"
)

// SetupEventSubscriptions forja la alianza entre tus módulos independientes.
func SetupEventSubscriptions(streakService *services.StreakService, trophyService *services.TrophyService) {
	eventBus := bus.GetInstance()

	if streakService != nil {
		eventBus.Subscribe("ritual.completed", streakService.OnRitualCompleted)
	}

	if trophyService != nil {
		eventBus.Subscribe("ritual.completed", trophyService.OnRitualCompleted)
		eventBus.Subscribe("habit.completed", trophyService.OnHabitCompleted)
	}
}
