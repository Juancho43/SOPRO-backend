package handlers

import (
	"net/http"
	"time"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	usecases "github.com/Juancho43/SOPRO-backend/pkg/internal/use_cases"
	"github.com/Juancho43/SOPRO-backend/pkg/src/services"
	"github.com/gin-gonic/gin"
)

type DomainHandler struct {
	service *services.DomainService
}

func NewDomainHandler(service *services.DomainService) *DomainHandler {
	return &DomainHandler{service: service}
}

func (h *DomainHandler) SetRoutes(router *gin.Engine, authMiddleware gin.HandlerFunc) {
	domains := router.Group("api/domains", authMiddleware)
	{
		domains.POST("", h.CreateDomain)
		domains.GET("", h.ListDomains)
		domains.PUT("/:id", h.UpdateDomain)
		domains.DELETE("/:id", h.DeleteDomain)
		domains.GET("/radar", h.GetRadarMetrics)
	}
}

func getAuthUserID(c *gin.Context) string {
	if userObj, exists := c.Get("user"); exists {
		if user, ok := userObj.(*entities.User); ok && user != nil {
			return user.UID
		}
	}
	return c.GetString("user_id")
}

// CreateDomain godoc
// @Summary Crear un dominio de excelencia
// @Description Permite categorizar áreas de desarrollo personal con paleta de color personalizada
// @Tags Domains
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body usecases.CreateDomainCommand true "Datos del dominio"
// @Success 201 {object} entities.Domain
// @Failure 400 {object} map[string]string
// @Router /api/domains [post]
func (h *DomainHandler) CreateDomain(c *gin.Context) {
	userID := getAuthUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autorizado"})
		return
	}

	var cmd usecases.CreateDomainCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}
	cmd.UserID = userID

	domain, err := h.service.CreateDomain(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, domain)
}

// ListDomains godoc
// @Summary Listar dominios de excelencia del usuario
// @Description Recupera los dominios configurados con el recuento de hábitos asociados
// @Tags Domains
// @Produce json
// @Security BearerAuth
// @Success 200 {array} entities.DomainWithHabitCount
// @Router /api/domains [get]
func (h *DomainHandler) ListDomains(c *gin.Context) {
	userID := getAuthUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autorizado"})
		return
	}

	domains, err := h.service.ListDomains(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if domains == nil {
		domains = []*entities.DomainWithHabitCount{}
	}

	c.JSON(http.StatusOK, domains)
}

// UpdateDomain godoc
// @Summary Actualizar un dominio de excelencia
// @Tags Domains
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID del dominio"
// @Param request body usecases.UpdateDomainCommand true "Datos a actualizar"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/domains/{id} [put]
func (h *DomainHandler) UpdateDomain(c *gin.Context) {
	userID := getAuthUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autorizado"})
		return
	}

	id := c.Param("id")
	var cmd usecases.UpdateDomainCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}
	cmd.ID = id
	cmd.UserID = userID

	if err := h.service.UpdateDomain(c.Request.Context(), cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Dominio actualizado con éxito"})
}

// DeleteDomain godoc
// @Summary Eliminar un dominio de excelencia
// @Tags Domains
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID del dominio"
// @Success 200 {object} map[string]string
// @Router /api/domains/{id} [delete]
func (h *DomainHandler) DeleteDomain(c *gin.Context) {
	userID := getAuthUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autorizado"})
		return
	}

	id := c.Param("id")
	if err := h.service.DeleteDomain(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Dominio eliminado con éxito"})
}

// GetRadarMetrics godoc
// @Summary Métricas de radar semanal por dominio
// @Description Retorna el esfuerzo agregado de los últimos 7 días por dominio para graficar en radar
// @Tags Domains
// @Produce json
// @Security BearerAuth
// @Success 200 {array} entities.DomainEffortPoint
// @Router /api/domains/radar [get]
func (h *DomainHandler) GetRadarMetrics(c *gin.Context) {
	userID := getAuthUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autorizado"})
		return
	}

	metrics, err := h.service.GetRadarMetrics(c.Request.Context(), userID, time.Now())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if metrics == nil {
		metrics = []*entities.DomainEffortPoint{}
	}

	c.JSON(http.StatusOK, metrics)
}
