package injectors

import (
	"github.com/Juancho43/SOPRO-backend/pkg/core"
	"github.com/Juancho43/SOPRO-backend/pkg/external/postgress"
	usecases "github.com/Juancho43/SOPRO-backend/pkg/internal/use_cases"
	"github.com/Juancho43/SOPRO-backend/pkg/src/handlers"
	"github.com/Juancho43/SOPRO-backend/pkg/src/services"
	"github.com/gin-gonic/gin"
)

func ChallengeModule(app *core.Core, authMiddleware gin.HandlerFunc) *services.ChallengeService {
	repo := postgress.NewPostgresChallengeRepository(app.Database)

	createChallenge := usecases.NewCreateChallengeUseCase(repo)
	manageChallenges, _ := usecases.NewManageChallengesUseCase(repo, app.Config.Timezone)
	manageMilestones, _ := usecases.NewManageMilestonesUseCase(repo, app.Config.Timezone)

	service := services.NewChallengeService(
		createChallenge,
		manageChallenges,
		manageMilestones,
		app.Logger,
	)

	handler := handlers.NewChallengeHandler(service)
	handler.SetRoutes(app.Router, authMiddleware)

	return service
}
