package injectors

import (
	"github.com/Juancho43/SOPRO-backend/app/core"
	"github.com/Juancho43/SOPRO-backend/app/external"
	"github.com/Juancho43/SOPRO-backend/app/external/postgress"
	"github.com/Juancho43/SOPRO-backend/app/internal"
	usecases "github.com/Juancho43/SOPRO-backend/app/internal/use_cases"
	"github.com/Juancho43/SOPRO-backend/app/src/handlers"
	"github.com/Juancho43/SOPRO-backend/app/src/middleware"
	"github.com/Juancho43/SOPRO-backend/app/src/services"
)

func AuthModule(core *core.Core) {
	repo := postgress.NewPostgresUserRepository(core.Database)
	usecase := usecases.NewLoginUseCase(repo)
	service := services.NewAuthService(usecase)
	handler := handlers.NewAuthHandler(service)
	client, _ := external.GetFirebaseAuthClient()
	firebase := internal.NewFirebaseVerifier(client, core.Logger)
	middleware := middleware.FirebaseAuthMiddleware(firebase)
	handler.SetRoutes(core.Router, middleware)
}
