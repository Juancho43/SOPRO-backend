package handlers

import (
	"net/http"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	"github.com/Juancho43/SOPRO-backend/pkg/src/services"
	"github.com/gin-gonic/gin"
)

type StreakHandler struct {
	service *services.StreakService
}

func NewStreakHandler(service *services.StreakService) *StreakHandler {
	return &StreakHandler{service: service}
}

func (h *StreakHandler) SetRoutes(router *gin.Engine, authMiddleware gin.HandlerFunc) {
	streaks := router.Group("api/streaks", authMiddleware)
	{
		// Modificamos la ruta para exigir el objetivo exacto como Path Parameter[cite: 22]
		streaks.GET("/habit/:habitName", h.GetStreakByHabit)
	}
}

// GetStreakByHabit godoc
// @Summary Obtener la racha de un hábito específico
// @Description Recupera la inercia actual y máxima de un hábito para medir tu nivel de disciplina.
// @Tags Streaks
// @Accept json
// @Produce json
// @Param habitName path string true "Nombre del hábito (ej. Programacion)"
// @Success 200 {object} entities.HabitStreak "El estado actual de tu disciplina"
// @Router /api/streaks/habit/{habitName} [get]
// @Security BearerAuth
func (h *StreakHandler) GetStreakByHabit(c *gin.Context) {
	userObj, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Usuario no encontrado en el contexto"})
		return
	}

	user, _ := userObj.(*entities.User)
	// 2. Extraer el objetivo desde la variable de la URL
	habitName := c.Param("habitName")

	// 3. Delegar la ejecución al servicio
	streak, err := h.service.GetStreak(user.UID, habitName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 4. Entregar los resultados
	c.JSON(http.StatusOK, streak)
}
