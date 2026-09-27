package injectors

import (
	"github.com/Juancho43/SOPRO-backend/pkg/src/handlers"
	"github.com/gin-gonic/gin"
)

func HealthModule(c *gin.Engine) {
	handler := handlers.NewHealHandler()
	handler.SetRoutes(c)
}
