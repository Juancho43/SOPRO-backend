package handlers

import (
	"net/http"

	usecases "github.com/Juancho43/SOPRO-backend/pkg/internal/use_cases"
	"github.com/Juancho43/SOPRO-backend/pkg/src/services"
	"github.com/gin-gonic/gin"
)

type RitualHandler struct {
	service *services.RitualService
}

func NewRitualHandler(service *services.RitualService) *RitualHandler {
	return &RitualHandler{service: service}
}

func (h *RitualHandler) SetRoutes(router *gin.Engine, authMiddleware gin.HandlerFunc) {
	rituals := router.Group("api/rituals", authMiddleware)
	{
		rituals.POST("", h.CreateRitual)
		rituals.GET("/today", h.CheckTodayRitual)
		rituals.GET("", h.GetRitualsHistory)
	}
}

// CreateRitual	godoc
//
// @Summary Crea tu ritual diario
// @Description Registra el agradecimiento y las 10 metas innegociables del día. Regla inquebrantable: Solo se permite uno por día calendario.
// @Tags Ritual
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body usecases.DailyRitualCommand true "Tu enfoque y gratitud del día"
// @Success 201 {object} map[string]string "¡Ritual completado! Tu enfoque está asegurado para hoy."
// @Failure 400 {object} map[string]string "La excelencia requiere precisión (faltan datos o no son 10 metas)."
// @Failure 401 {object} map[string]string "Sin disciplina no hay resultados (No autorizado)."
// @Failure 409 {object} map[string]string "Ya has forjado tu disciplina hoy."
// @Failure 500 {object} map[string]string "Error interno del servidor."
// @Router /api/rituals [post]
func (h *RitualHandler) CreateRitual(c *gin.Context) {
	var cmd usecases.DailyRitualCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "La excelencia requiere precisión en los datos enviados."})
		return
	}

	// NOTA: Idealmente extraes el UserID del token JWT en tu AuthMiddleware
	// cmd.UserID = c.GetString("user_id")

	if err := h.service.ExecuteDailyRitual(cmd); err != nil {
		// Evaluamos si es error de validación/duplicidad o del servidor
		if err.Error() == "ya has forjado tu disciplina hoy; ahora enfócate en ejecutar" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "¡Ritual completado! Tu enfoque está asegurado para hoy."})
}

// CheckTodayRitual	godoc
// @Summary Verifica y obtiene el ritual de hoy
// @Description Retorna el agradecimiento y las metas establecidas para el día actual si ya fue creado.
// @Tags Ritual
// @Produce json
// @Security BearerAuth
// @Success 200 {object} bool "Ritual del día obtenido con éxito."
// @Failure 401 {object} map[string]string "No autorizado."
// @Failure 404 {object} map[string]string "Aún no has forjado tu disciplina hoy."
// @Failure 500 {object} map[string]string "Error interno del servidor."
// @Router /api/rituals/today [get]
func (h *RitualHandler) CheckTodayRitual(c *gin.Context) {
	userID := c.GetString("user_id") // Extraído de tu middleware JWT
	if userID == "" {
		// Mock temporal por si aún no tienes el middleware conectado
		userID = "default-user-id"
	}

	ritual, err := h.service.ExecuteGetTodayRitual(userID)
	if err != nil || ritual == false {
		c.JSON(http.StatusNotFound, gin.H{"error": "Aún no has forjado tu disciplina hoy."})
		return
	}

	c.JSON(http.StatusOK, ritual)
}

// GetRitualsHistory godoc
// @Summary Obtiene el historial de rituales
// @Description Devuelve la lista histórica de todos los rituales diarios forjados por el usuario.
// @Tags Ritual
// @Produce json
// @Security BearerAuth
// @Success 200 {array} entities.Ritual "Historial obtenido con éxito."
// @Failure 401 {object} map[string]string "No autorizado."
// @Router /api/rituals [get]
func (h *RitualHandler) GetRitualsHistory(c *gin.Context) {
	// Aquí conectarías el servicio para obtener el historial completo
	c.JSON(http.StatusOK, gin.H{"message": "El éxito deja huellas: Historial en construcción."})
}
