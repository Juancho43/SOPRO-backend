package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	usecases "github.com/Juancho43/SOPRO-backend/pkg/internal/use_cases"
	"github.com/Juancho43/SOPRO-backend/pkg/src/services"
	"github.com/gin-gonic/gin"
)

type HabitHandler struct {
	service *services.HabitService
}

func NewHabitHandler(service *services.HabitService) *HabitHandler {
	return &HabitHandler{service: service}
}

func (h *HabitHandler) SetRoutes(router *gin.Engine, authMiddleware gin.HandlerFunc) {
	habits := router.Group("api/habits", authMiddleware)
	{
		habits.POST("", h.CreateHabit)
		habits.GET("", h.ListHabits)
		habits.POST("/:id/check-in", h.CheckInHabit)
		habits.GET("/:id/history", h.GetHabitHistory)
	}
}

// CreateHabit godoc
// @Summary Crear un nuevo hábito en el arsenal
// @Description Permite registrar un hábito individual con frecuencia diaria o semanal y dominio opcional
// @Tags Habits
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body usecases.CreateHabitCommand true "Datos del hábito"
// @Success 201 {object} entities.Habit
// @Failure 400 {object} map[string]string
// @Router /api/habits [post]
func (h *HabitHandler) CreateHabit(c *gin.Context) {
	userID := getAuthUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autorizado"})
		return
	}

	var cmd usecases.CreateHabitCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}
	cmd.UserID = userID

	habit, err := h.service.CreateHabit(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, habit)
}

// ListHabits godoc
// @Summary Listar hábitos del arsenal con rachas
// @Description Retorna los hábitos activos del usuario, sus rachas y si ya fueron completados hoy
// @Tags Habits
// @Produce json
// @Security BearerAuth
// @Success 200 {array} entities.HabitWithStreak
// @Router /api/habits [get]
func (h *HabitHandler) ListHabits(c *gin.Context) {
	userID := getAuthUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autorizado"})
		return
	}

	habits, err := h.service.ListUserHabits(c.Request.Context(), userID, time.Now())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if habits == nil {
		habits = []*entities.HabitWithStreak{}
	}

	c.JSON(http.StatusOK, habits)
}

// CheckInHabit godoc
// @Summary Registrar cumplimiento diario de un hábito
// @Description Marca el hábito como completado hoy, recalcula racha y emite evento de gamificación
// @Tags Habits
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID del hábito"
// @Success 200 {object} entities.HabitCheckInResult
// @Failure 404 {object} map[string]string "Hábito no encontrado"
// @Failure 409 {object} map[string]string "Ya registrado hoy"
// @Router /api/habits/{id}/check-in [post]
func (h *HabitHandler) CheckInHabit(c *gin.Context) {
	userID := getAuthUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autorizado"})
		return
	}

	habitID := c.Param("id")
	result, err := h.service.CheckInHabit(c.Request.Context(), habitID, userID, time.Now())
	if err != nil {
		if errors.Is(err, usecases.ErrAlreadyCheckedIn) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, usecases.ErrHabitNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetHabitHistory godoc
// @Summary Historial de cumplimiento para mapa de calor
// @Description Retorna las fechas (YYYY-MM-DD) de cumplimiento de los últimos 365 días
// @Tags Habits
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID del hábito"
// @Success 200 {array} string
// @Router /api/habits/{id}/history [get]
func (h *HabitHandler) GetHabitHistory(c *gin.Context) {
	userID := getAuthUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autorizado"})
		return
	}

	habitID := c.Param("id")
	history, err := h.service.GetHabitHistory(c.Request.Context(), habitID, userID, time.Now())
	if err != nil {
		if errors.Is(err, usecases.ErrHabitNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if history == nil {
		history = []string{}
	}

	c.JSON(http.StatusOK, history)
}
