package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"server_aquascan/services"
	"server_aquascan/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil token dari header Authorization
		authHeader := c.GetHeader("Authorization")
		parts := strings.Split(authHeader, " ")
		if authHeader == "" || len(parts) != 2 || parts[0] != "Bearer" {
			utils.RespondError(c, http.StatusUnauthorized, "Authorization header dibutuhkan atau format token tidak valid", nil)
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Parse dan validasi token
		token, err := services.ParseJWT(tokenString)
		if err != nil || !token.Valid {
			utils.RespondError(c, http.StatusUnauthorized, "Token tidak valid atau expired", err.Error())
			c.Abort()
			return
		}

		// Set claims ke context
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if uid, exists := claims["user_id"]; exists {
				c.Set("user_id", uid)
			}
			if email, exists := claims["email"]; exists {
				c.Set("email", email)
			}
			if role, exists := claims["role"]; exists {
				c.Set("role", role)
			}
		}

		c.Next()
	}
}
