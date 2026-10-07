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
	HealthModule(app.Router)
	AuthModule(app, middleware)
	RitualModule(app, middleware)
	StreakModule(app, middleware)
}
