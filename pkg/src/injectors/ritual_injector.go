package injectors

import (
	"github.com/Juancho43/SOPRO-backend/pkg/core"
	"github.com/Juancho43/SOPRO-backend/pkg/external/postgress"
	usecases "github.com/Juancho43/SOPRO-backend/pkg/internal/use_cases"
	"github.com/Juancho43/SOPRO-backend/pkg/src/handlers"
	"github.com/Juancho43/SOPRO-backend/pkg/src/services"
	"github.com/gin-gonic/gin"
)

func RitualModule(core *core.Core, authMiddleware gin.HandlerFunc) {
	repo := postgress.NewPostgresRitualRepository(core.Database)
	useCase, _ := usecases.NewDailyRitual(repo, core.Config.Timezone)
	todaysUsecase, _ := usecases.NewTodaysRitual(repo, core.Config.Timezone)
	service := services.NewRitualService(useCase, todaysUsecase, core.Logger)
	handler := handlers.NewRitualHandler(service)
	handler.SetRoutes(core.Router, authMiddleware)
}
