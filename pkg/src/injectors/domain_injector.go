package injectors

import (
	"github.com/Juancho43/SOPRO-backend/pkg/core"
	"github.com/Juancho43/SOPRO-backend/pkg/external/postgress"
	usecases "github.com/Juancho43/SOPRO-backend/pkg/internal/use_cases"
	"github.com/Juancho43/SOPRO-backend/pkg/src/handlers"
	"github.com/Juancho43/SOPRO-backend/pkg/src/services"
	"github.com/gin-gonic/gin"
)

func DomainModule(app *core.Core, authMiddleware gin.HandlerFunc) *services.DomainService {
	repo := postgress.NewPostgresDomainRepository(app.Database)

	createDomain := usecases.NewCreateDomainUseCase(repo)
	listDomains := usecases.NewListDomainsUseCase(repo)
	updateDomain := usecases.NewUpdateDomainUseCase(repo)
	deleteDomain := usecases.NewDeleteDomainUseCase(repo)
	radarMetrics := usecases.NewGetDomainRadarMetricsUseCase(repo)

	service := services.NewDomainService(
		createDomain,
		listDomains,
		updateDomain,
		deleteDomain,
		radarMetrics,
		app.Logger,
	)

	handler := handlers.NewDomainHandler(service)
	handler.SetRoutes(app.Router, authMiddleware)

	return service
}
