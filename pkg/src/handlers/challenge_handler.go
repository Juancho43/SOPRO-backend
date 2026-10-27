package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	usecases "github.com/Juancho43/SOPRO-backend/pkg/internal/use_cases"
	"github.com/Juancho43/SOPRO-backend/pkg/src/services"
	"github.com/gin-gonic/gin"
)

type ChallengeHandler struct {
	service *services.ChallengeService
}

func NewChallengeHandler(service *services.ChallengeService) *ChallengeHandler {
	return &ChallengeHandler{service: service}
}

func (h *ChallengeHandler) SetRoutes(router *gin.Engine, authMiddleware gin.HandlerFunc) {
	challenges := router.Group("api/v1/challenges", authMiddleware)
	{
		// Endpoints de Retos (HU 3.1)
		challenges.POST("", h.CreateChallenge)
		challenges.GET("", h.ListChallenges)
		challenges.GET("/upcoming-deadlines", h.GetUpcomingDeadlines)
		challenges.GET("/:id", h.GetChallengeDetail)
		challenges.PUT("/:id", h.UpdateChallenge)
		challenges.PATCH("/:id/status", h.UpdateChallengeStatus)

		// Endpoints de Hitos (HU 3.2)
		challenges.POST("/:id/milestones", h.CreateMilestone)
		challenges.PATCH("/:id/milestones/:milestoneId", h.UpdateMilestone)
		challenges.DELETE("/:id/milestones/:milestoneId", h.DeleteMilestone)
	}
}

// CreateChallenge godoc
// @Summary Crear un reto estratégico
// @Description Permite definir un reto a plazo fijo con fechas, objetivo y recompensa final.
// @Tags Challenges
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body usecases.CreateChallengeCommand true "Datos del reto"
// @Success 201 {object} entities.Challenge
// @Failure 400 {object} map[string]string
// @Router /api/v1/challenges [post]
func (h *ChallengeHandler) CreateChallenge(c *gin.Context) {
	userID := getAuthUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autorizado"})
		return
	}

	var cmd usecases.CreateChallengeCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}
	cmd.UserID = userID

	challenge, err := h.service.CreateChallenge(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, challenge)
}

// ListChallenges godoc
// @Summary Listar retos estratégicos del usuario
// @Description Retorna los retos del usuario filtrables opcionalmente por status (ej: ACTIVE, COMPLETED).
// @Tags Challenges
// @Produce json
// @Security BearerAuth
// @Param status query string false "Filtrar por estado (ACTIVE, COMPLETED, FAILED, ARCHIVED, CANCELLED)"
// @Success 200 {array} entities.Challenge
// @Router /api/v1/challenges [get]
func (h *ChallengeHandler) ListChallenges(c *gin.Context) {
	userID := getAuthUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autorizado"})
		return
	}

	var statusFilter *string
	if q := c.Query("status"); q != "" {
		statusFilter = &q
	}

	challenges, err := h.service.ListChallenges(c.Request.Context(), userID, statusFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, challenges)
}

// GetChallengeDetail godoc
// @Summary Obtener detalle de un reto con sus hitos y progreso
// @Description Retorna el reto, su lista de hitos, porcentaje de progreso y días restantes.
// @Tags Challenges
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID del Reto"
// @Success 200 {object} entities.ChallengeDetail
// @Failure 404 {object} map[string]string
// @Router /api/v1/challenges/{id} [get]
func (h *ChallengeHandler) GetChallengeDetail(c *gin.Context) {
	userID := getAuthUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autorizado"})
		return
	}

	challengeID := c.Param("id")
	detail, err := h.service.GetChallengeDetail(c.Request.Context(), challengeID, userID, time.Now())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, detail)
}

// UpdateChallenge godoc
// @Summary Actualizar título, descripción o recompensa de un reto
// @Tags Challenges
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID del Reto"
// @Param request body usecases.UpdateChallengeCommand true "Datos a modificar"
// @Success 200 {object} entities.Challenge
// @Failure 400 {object} map[string]string
// @Router /api/v1/challenges/{id} [put]
func (h *ChallengeHandler) UpdateChallenge(c *gin.Context) {
	userID := getAuthUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autorizado"})
		return
	}

	var cmd usecases.UpdateChallengeCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}
	cmd.ID = c.Param("id")
	cmd.UserID = userID

	updated, err := h.service.UpdateChallenge(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}

type updateStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// UpdateChallengeStatus godoc
// @Summary Cambiar estado del reto (COMPLETED, CANCELLED, ARCHIVED, etc.)
// @Tags Challenges
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID del Reto"
// @Param request body updateStatusRequest true "Nuevo estado"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/v1/challenges/{id}/status [patch]
func (h *ChallengeHandler) UpdateChallengeStatus(c *gin.Context) {
	userID := getAuthUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autorizado"})
		return
	}

	var req updateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El campo 'status' es obligatorio"})
		return
	}

	challengeID := c.Param("id")
	err := h.service.UpdateChallengeStatus(c.Request.Context(), challengeID, userID, entities.ChallengeStatus(req.Status))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Estado de reto actualizado con éxito"})
}

// CreateMilestone godoc
// @Summary Añadir un hito o deadline intermedia a un reto
// @Tags Challenges
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID del Reto"
// @Param request body usecases.CreateMilestoneCommand true "Datos del hito"
// @Success 201 {object} entities.Milestone
// @Failure 400 {object} map[string]string
// @Router /api/v1/challenges/{id}/milestones [post]
func (h *ChallengeHandler) CreateMilestone(c *gin.Context) {
	userID := getAuthUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autorizado"})
		return
	}

	var cmd usecases.CreateMilestoneCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}
	cmd.UserID = userID
	cmd.ChallengeID = c.Param("id")

	milestone, err := h.service.CreateMilestone(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, milestone)
}

// UpdateMilestone godoc
// @Summary Modificar datos o marcar un hito como completado
// @Tags Challenges
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID del Reto"
// @Param milestoneId path string true "ID del Hito"
// @Param request body usecases.UpdateMilestoneCommand true "Campos a modificar"
// @Success 200 {object} entities.Milestone
// @Failure 400 {object} map[string]string
// @Router /api/v1/challenges/{id}/milestones/{milestoneId} [patch]
func (h *ChallengeHandler) UpdateMilestone(c *gin.Context) {
	userID := getAuthUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autorizado"})
		return
	}

	var cmd usecases.UpdateMilestoneCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}
	cmd.ChallengeID = c.Param("id")
	cmd.ID = c.Param("milestoneId")
	cmd.UserID = userID

	milestone, err := h.service.UpdateMilestone(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, milestone)
}

// DeleteMilestone godoc
// @Summary Eliminar un hito de un reto
// @Tags Challenges
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID del Reto"
// @Param milestoneId path string true "ID del Hito"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/v1/challenges/{id}/milestones/{milestoneId} [delete]
func (h *ChallengeHandler) DeleteMilestone(c *gin.Context) {
	userID := getAuthUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autorizado"})
		return
	}

	challengeID := c.Param("id")
	milestoneID := c.Param("milestoneId")

	err := h.service.DeleteMilestone(c.Request.Context(), milestoneID, challengeID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Hito eliminado correctamente"})
}

// GetUpcomingDeadlines godoc
// @Summary Próximas deadlines intermedias para alertas del Dashboard
// @Description Devuelve hitos pendientes de retos activos con vencimiento próximo (por defecto en los siguientes 3 días)
// @Tags Challenges
// @Produce json
// @Security BearerAuth
// @Param days query int false "Días de anticipación (default: 3)"
// @Success 200 {array} entities.UpcomingDeadline
// @Router /api/v1/challenges/upcoming-deadlines [get]
func (h *ChallengeHandler) GetUpcomingDeadlines(c *gin.Context) {
	userID := getAuthUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autorizado"})
		return
	}

	days := 3
	if q := c.Query("days"); q != "" {
		if val, err := strconv.Atoi(q); err == nil && val > 0 {
			days = val
		}
	}

	deadlines, err := h.service.GetUpcomingDeadlines(c.Request.Context(), userID, days, time.Now())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, deadlines)
}
