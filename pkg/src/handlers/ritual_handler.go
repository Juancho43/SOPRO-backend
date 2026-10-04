package handlers

import (
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
// @Success 200 {object} map[string]string "¡Ritual completado! Tu enfoque está asegurado para hoy."
// @Failure 400 {object} map[string]string "La excelencia requiere precisión (faltan datos o no son 10 metas)."
// @Failure 401 {object} map[string]string "Sin disciplina no hay resultados (No autorizado)."
// @Failure 409 {object} map[string]string "Ya has forjado tu disciplina hoy."
// @Failure 500 {object} map[string]string "Error interno del servidor."
// @Router /api/rituals [post]
func (h *RitualHandler) CreateRitual(c *gin.Context) {

}

// CheckTodayRitual	godoc
// @Tags Ritual
// @Router /api/rituals/today [get]
func (h *RitualHandler) CheckTodayRitual(c *gin.Context) {

}

// GetRitualsHistory godoc
// @Tags Ritual
// @Router /api/rituals [get]
func (h *RitualHandler) GetRitualsHistory(c *gin.Context)
