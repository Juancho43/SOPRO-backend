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

func TrophyModule(app *core.Core, authMiddleware gin.HandlerFunc, streakService *services.StreakService) *services.TrophyService {
	trophyRepo := postgress.NewPostgresTrophyRepository(app.Database)
	habitRepo := postgress.NewPostgresHabitRepository(app.Database)
	eventBus := bus.GetInstance()

	createTrophy := usecases.NewCreateTrophyUseCase(trophyRepo, habitRepo)
	showcase := usecases.NewGetUserTrophyShowcaseUseCase(trophyRepo)
	listTrophies := usecases.NewListTrophiesUseCase(trophyRepo)
	evaluateTrophies := usecases.NewEvaluateTrophiesUseCase(trophyRepo, eventBus)
	evaluateSynergy := usecases.NewEvaluateSynergyTrophiesUseCase(trophyRepo, habitRepo, eventBus)
	synergyStatus := usecases.NewGetSynergyStatusUseCase(trophyRepo, habitRepo)

	service := services.NewTrophyService(
		createTrophy,
		showcase,
		listTrophies,
		evaluateTrophies,
		evaluateSynergy,
		synergyStatus,
		streakService,
		app.Logger,
	)

	handler := handlers.NewTrophyHandler(service)
	handler.SetRoutes(app.Router, authMiddleware)

	return service
}
