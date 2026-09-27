package middleware

import (
	"net/http"
	"strings"

	"github.com/Juancho43/SOPRO-backend/app/internal"
	"github.com/gin-gonic/gin"
)

func FirebaseAuthMiddleware(verifier *internal.FirebaseVerifier) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header requerido"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, "Bearer ")
		if len(parts) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Formato de token inválido"})
			c.Abort()
			return
		}

		idToken := parts[1]
		// Utiliza tu FirebaseVerifier para decodificar el JWT
		user, err := verifier.Verify(idToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		// Guarda la entidad User en el contexto para que el handler la recoja
		c.Set("user", user)
		c.Next()
	}
}
