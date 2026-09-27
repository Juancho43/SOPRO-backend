package injectors

import (
	"github.com/Juancho43/SOPRO-backend/pkg/core"
	"github.com/Juancho43/SOPRO-backend/pkg/external"
	"github.com/Juancho43/SOPRO-backend/pkg/external/postgress"
	"github.com/Juancho43/SOPRO-backend/pkg/internal"
	usecases "github.com/Juancho43/SOPRO-backend/pkg/internal/use_cases"
	"github.com/Juancho43/SOPRO-backend/pkg/src/handlers"
	"github.com/Juancho43/SOPRO-backend/pkg/src/middleware"
	"github.com/Juancho43/SOPRO-backend/pkg/src/services"
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
