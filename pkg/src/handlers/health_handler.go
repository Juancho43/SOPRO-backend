package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealtHandler struct {
}

func NewHealHandler() *HealtHandler {
	return &HealtHandler{}
}

func (h *HealtHandler) SetRoutes(router *gin.Engine) {
	router.GET("api/health", h.Health)
}

// Health godoc
// @Summary Comprueba la vitalidad del núcleo
// @Description Verifica que el motor del sistema esté en línea y operando a su máxima capacidad.
// @Tags Sistema
// @Produce json
// @Success 200 {object} map[string]string "Retorna un mensaje de victoria garantizada."
// @Router /api/health [get]
func (h *HealtHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "EXCELENT COMO LOS CAMPEONES QUE SOMOS",
	})
}
