package handlers

import (
	"net/http"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	usecases "github.com/Juancho43/SOPRO-backend/pkg/internal/use_cases"
	"github.com/Juancho43/SOPRO-backend/pkg/src/services"
	"github.com/gin-gonic/gin"
)

type TrophyHandler struct {
	service *services.TrophyService
}

func NewTrophyHandler(service *services.TrophyService) *TrophyHandler {
	return &TrophyHandler{service: service}
}

func (h *TrophyHandler) SetRoutes(router *gin.Engine, authMiddleware gin.HandlerFunc) {
	trophies := router.Group("api/trophies", authMiddleware)
	{
		trophies.POST("", h.CreateTrophy)
		trophies.GET("", h.ListTrophies)
		trophies.GET("/:id/synergy-status", h.GetSynergyStatus)
	}

	userShowcase := router.Group("api/users/me/trophies", authMiddleware)
	{
		userShowcase.GET("", h.GetUserShowcase)
	}
}

// CreateTrophy godoc
// @Summary Diseñar un nuevo trofeo o logro
// @Description Permite definir trofeos simples vinculados a hábitos o racha maestra, o trofeos de sinergia multihábito (AND)
// @Tags Trophies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body usecases.CreateTrophyCommand true "Definición del trofeo"
// @Success 201 {object} entities.Trophy
// @Failure 400 {object} map[string]string
// @Router /api/trophies [post]
func (h *TrophyHandler) CreateTrophy(c *gin.Context) {
	userID := getAuthUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autorizado"})
		return
	}

	var cmd usecases.CreateTrophyCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}
	cmd.UserID = userID

	trophy, err := h.service.CreateTrophy(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, trophy)
}

// ListTrophies godoc
// @Summary Listar trofeos creados
// @Tags Trophies
// @Produce json
// @Security BearerAuth
// @Success 200 {array} entities.Trophy
// @Router /api/trophies [get]
func (h *TrophyHandler) ListTrophies(c *gin.Context) {
	userID := getAuthUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autorizado"})
		return
	}

	trophies, err := h.service.ListTrophies(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if trophies == nil {
		trophies = []*entities.Trophy{}
	}

	c.JSON(http.StatusOK, trophies)
}

// GetUserShowcase godoc
// @Summary Vitrina de trofeos del usuario autenticado
// @Description Retorna todos los trofeos indicando si han sido desbloqueados y en qué fecha
// @Tags Trophies
// @Produce json
// @Security BearerAuth
// @Success 200 {array} entities.ShowcaseTrophy
// @Router /api/users/me/trophies [get]
func (h *TrophyHandler) GetUserShowcase(c *gin.Context) {
	userID := getAuthUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autorizado"})
		return
	}

	showcase, err := h.service.GetUserShowcase(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if showcase == nil {
		showcase = []*entities.ShowcaseTrophy{}
	}

	c.JSON(http.StatusOK, showcase)
}

// GetSynergyStatus godoc
// @Summary Consultar estado de un trofeo de sinergia
// @Description Retorna el progreso detallado de cada hábito que compone el trofeo compuesto
// @Tags Trophies
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID del trofeo"
// @Success 200 {object} entities.SynergyStatusResponse
// @Failure 400 {object} map[string]string
// @Router /api/trophies/{id}/synergy-status [get]
func (h *TrophyHandler) GetSynergyStatus(c *gin.Context) {
	userID := getAuthUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autorizado"})
		return
	}

	id := c.Param("id")
	status, err := h.service.GetSynergyStatus(c.Request.Context(), id, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}
