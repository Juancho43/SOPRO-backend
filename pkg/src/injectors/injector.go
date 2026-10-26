package injectors

import (
	_ "github.com/Juancho43/SOPRO-backend/docs"
	"github.com/Juancho43/SOPRO-backend/pkg/core"
	"github.com/Juancho43/SOPRO-backend/pkg/external"
	"github.com/Juancho43/SOPRO-backend/pkg/internal"
	"github.com/Juancho43/SOPRO-backend/pkg/src/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InjectAll(app *core.Core) {
	app.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	external.GetFirebaseApp()
	client, _ := external.GetFirebaseAuthClient()
	firebase := internal.NewFirebaseVerifier(client, app.Logger)
	middleware := middleware.FirebaseAuthMiddleware(firebase)

	// Épica 1: El Ritual del Ganador
	HealthModule(app.Router)
	AuthModule(app, middleware)
	RitualModule(app, middleware)
	streakService := StreakModule(app, middleware)
	SessionsModule(app, middleware)

	// Épica 2: El Arsenal de Hábitos y Trofeos
	DomainModule(app, middleware)
	HabitModule(app, middleware)
	trophyService := TrophyModule(app, middleware, streakService)

	// Épica 3: Conquista de Retos y Sprints
	ChallengeModule(app, middleware)

	// Event Bus Subscriptions
	SetupEventSubscriptions(streakService, trophyService)
}
