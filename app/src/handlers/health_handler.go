package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealtHandler struct {
}

func NewHealHandler() *HealtHandler {
	return &HealtHandler{}
}

func (h *HealtHandler) SetRoutes(router *gin.Engine) {
	router.GET("api/health", h.Health)
}

func (h *HealtHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "EXCELENT COMO LOS CAMPEONES QUE SOMOS",
	})
}
