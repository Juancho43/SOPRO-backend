package handlers

import (
	"net/http"

	"github.com/Juancho43/SOPRO-backend/pkg/src/services"
	"github.com/gin-gonic/gin"
)

type FocusmateHandler struct {
	service *services.FocusmateService
}

func NewFocusmateHandler(service *services.FocusmateService) *FocusmateHandler {
	return &FocusmateHandler{service: service}
}

func (h *FocusmateHandler) SetRoutes(router *gin.Engine, authMiddleware gin.HandlerFunc) {
	focusmate := router.Group("api/focusmate", authMiddleware)
	{
		focusmate.GET("/sessions/future", h.GetAllFutureSessions)
	}
}

// GetAllFutureSessions godoc
// @Summary Obtiene todas las sesiones futuras de Focusmate
// @Description Retorna la lista de bloques de tiempo de alta concentración programados.
// @Tags Focusmate
// @Produce json
// @Security BearerAuth
// @Success 200 {array} entities.FocusmateSession "Sesiones obtenidas con éxito."
// @Failure 500 {object} map[string]string "Error interno al contactar la API."
// @Router /api/focusmate/sessions/future [get]
func (h *FocusmateHandler) GetAllFutureSessions(c *gin.Context) {
	// Reutilizamos el contexto de Gin para respetar tiempos de espera o cancelaciones
	sessions, err := h.service.ExecuteGetFutureSessions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al sincronizar con el entorno de alta concentración."})
		return
	}

	c.JSON(http.StatusOK, sessions)
}
