package handlers

import (
	"net/http"

	"github.com/Juancho43/SOPRO-backend/app/internal/model/entities"
	"github.com/Juancho43/SOPRO-backend/app/src/services"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// Asegúrate de inyectar el middleware en las rutas en tu main o router
func (h *AuthHandler) SetRoutes(router *gin.Engine, authMiddleware gin.HandlerFunc) {
	router.POST("api/login/google", authMiddleware, h.GoogleLogin)
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	userObj, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Usuario no encontrado en el contexto"})
		return
	}

	user, ok := userObj.(*entities.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error de parseo de usuario"})
		return
	}

	err := h.service.ExecuteGoogleLogin(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login exitoso", "user": user})
}
