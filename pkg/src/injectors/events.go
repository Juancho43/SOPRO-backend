package injectors

import (
	"github.com/Juancho43/SOPRO-backend/pkg/src/bus"
	"github.com/Juancho43/SOPRO-backend/pkg/src/services"
)

// SetupEventSubscriptions forja la alianza entre tus módulos independientes.
func SetupEventSubscriptions(streakService *services.StreakService) {
	eventBus := bus.GetInstance()

	eventBus.Subscribe("ritual.completed", streakService.OnRitualCompleted)

}
