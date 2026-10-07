package injectors

import (
	"github.com/Juancho43/SOPRO-backend/pkg/core"
	"github.com/Juancho43/SOPRO-backend/pkg/external/postgress"
	usecases "github.com/Juancho43/SOPRO-backend/pkg/internal/use_cases"
	"github.com/Juancho43/SOPRO-backend/pkg/src/handlers"
	"github.com/Juancho43/SOPRO-backend/pkg/src/services"
	"github.com/gin-gonic/gin"
)

func StreakModule(app *core.Core, authMiddleware gin.HandlerFunc) {
	repo := postgress.NewPostgresHabitStreakRepository(app.Database)
	markAsDone := usecases.NewMarkHabitAsDone(repo)
	getHabit := usecases.NewGetHabitStreak(repo)
	service := services.NewStreakService(markAsDone, getHabit, app.Logger)
	handler := handlers.NewStreakHandler(service)
	handler.SetRoutes(app.Router, authMiddleware)
	SetupEventSubscriptions(service)
}
