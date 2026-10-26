package injectors

import (
	"github.com/Juancho43/SOPRO-backend/pkg/core"
	"github.com/Juancho43/SOPRO-backend/pkg/external/postgress"
	usecases "github.com/Juancho43/SOPRO-backend/pkg/internal/use_cases"
	"github.com/Juancho43/SOPRO-backend/pkg/src/bus"
	"github.com/Juancho43/SOPRO-backend/pkg/src/handlers"
	"github.com/Juancho43/SOPRO-backend/pkg/src/services"
	"github.com/gin-gonic/gin"
)

func HabitModule(app *core.Core, authMiddleware gin.HandlerFunc) *services.HabitService {
	repo := postgress.NewPostgresHabitRepository(app.Database)
	eventBus := bus.GetInstance()

	createHabit := usecases.NewCreateHabitUseCase(repo)
	listHabits, _ := usecases.NewListUserHabitsUseCase(repo, app.Config.Timezone)
	checkInHabit, _ := usecases.NewCheckInHabitUseCase(repo, eventBus, app.Config.Timezone)
	habitHistory := usecases.NewGetHabitHistoryUseCase(repo)

	service := services.NewHabitService(
		createHabit,
		listHabits,
		checkInHabit,
		habitHistory,
		app.Logger,
	)

	handler := handlers.NewHabitHandler(service)
	handler.SetRoutes(app.Router, authMiddleware)

	return service
}
