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
		rituals.GET("/by-date", h.GetRitualByDate)
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

// GetRitualByDate godoc
// @Summary Obtiene un ritual por fecha específica
// @Description Devuelve el ritual registrado en una fecha enviada como query param (ej. ?date=YYYY-MM-DD).
// @Tags Rituals
// @Produce json
// @Security BearerAuth
// @Param date query string true "Fecha del ritual en formato YYYY-MM-DD"
// @Success 200 {object} entities.Ritual "Ritual obtenido con éxito."
// @Failure 400 {object} map[string]string "El parámetro date es requerido."
// @Failure 401 {object} map[string]string "No autorizado."
// @Failure 404 {object} map[string]string "Ritual no encontrado."
// @Router /api/rituals/by-date [get]
func (h *RitualHandler) GetRitualByDate(c *gin.Context) {
	userObj, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Usuario no encontrado en el contexto"})
		return
	}
	user, _ := userObj.(*entities.User)

	dateParam := c.Query("date")

	if dateParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "La claridad es poder: El parámetro 'date' es requerido (ej. ?date=2026-10-07)"})
		return
	}

	ritual, err := h.service.ExecuteGetRitualByDate(user.UID, dateParam)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No se encontraron registros para esta fecha. El éxito requiere acción continua."})
		return
	}

	c.JSON(http.StatusOK, ritual)
}
