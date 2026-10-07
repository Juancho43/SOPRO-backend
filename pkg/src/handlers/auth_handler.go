package handlers

import (
	"net/http"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	"github.com/Juancho43/SOPRO-backend/pkg/src/services"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) SetRoutes(router *gin.Engine, authMiddleware gin.HandlerFunc) {
	router.POST("api/login/google", authMiddleware, h.GoogleLogin)
}

// GoogleLogin godoc
// @Summary Autenticación maestra con Google
// @Description Gestiona el inicio de sesión o creación de cuenta. Requiere que el token de Firebase haya sido validado por el middleware.
// @Tags Auht
// @Produce json
// @Success 200 {object} map[string]interface{} "Login exitoso y datos del usuario."
// @Failure 500 {object} map[string]string "Error interno o usuario no encontrado en el contexto."
// @Router /api/login/google [post]
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
