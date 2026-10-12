package injectors

import (
	"github.com/Juancho43/SOPRO-backend/pkg/core"
	"github.com/Juancho43/SOPRO-backend/pkg/internal"
	usecases "github.com/Juancho43/SOPRO-backend/pkg/internal/use_cases"
	"github.com/Juancho43/SOPRO-backend/pkg/src/handlers"
	"github.com/Juancho43/SOPRO-backend/pkg/src/services"
	"github.com/gin-gonic/gin"
)

func SessionsModule(app *core.Core, authMiddleware gin.HandlerFunc) {
	repository, _ := internal.NewHTTPFocusmateAPI(app.Config.FocusmateAPIKEY, app.Config.Timezone)
	futureSessions := usecases.NewGetAllFutureSessions(repository)
	service := services.NewFocusmateService(futureSessions, app.Logger)
	handler := handlers.NewFocusmateHandler(service)
	handler.SetRoutes(app.Router, authMiddleware)
}
