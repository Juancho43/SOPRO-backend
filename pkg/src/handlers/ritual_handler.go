package handlers

import (
	"net/http"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
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
// @Tags Rituals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body usecases.DailyRitualCommand true "Tu enfoque y gratitud del día"
// @Success 201 {object} map[string]string "¡Ritual completado! Tu enfoque está asegurado para hoy."
// @Failure 401 {object} map[string]string "Sin disciplina no hay resultados (No autorizado)."
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
// @Tags Rituals
// @Produce json
// @Security BearerAuth
// @Success 200 {object} bool "Ritual del día obtenido con éxito."
// @Router /api/rituals/today [get]
func (h *RitualHandler) CheckTodayRitual(c *gin.Context) {
	userObj, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Usuario no encontrado en el contexto"})
		return
	}

	user, _ := userObj.(*entities.User)
	ritual, err := h.service.ExecuteGetTodayRitual(user.UID)
	if err != nil || ritual == false {
		c.JSON(http.StatusOK, false)
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
