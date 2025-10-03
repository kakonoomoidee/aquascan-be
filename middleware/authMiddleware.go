package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"server_aquascan/services"
	"server_aquascan/utils"
)

// AuthMiddleware mem-parse JWT, memvalidasi, dan menyimpan claim penting ke context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		parts := strings.Split(authHeader, " ")
		if authHeader == "" || len(parts) != 2 || parts[0] != "Bearer" {
			utils.RespondError(c, http.StatusUnauthorized, "Authorization header dibutuhkan atau format token tidak valid", nil)
			c.Abort()
			return
		}

		tokenString := parts[1]
		token, err := services.ParseJWT(tokenString)
		if err != nil || !token.Valid {
			utils.RespondError(c, http.StatusUnauthorized, "Token tidak valid atau expired", err.Error())
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utils.RespondError(c, http.StatusUnauthorized, "Token claims tidak valid", nil)
			c.Abort()
			return
		}

		// helper convert ke string (jwt sering encode numeric id sebagai float64)
		toStr := func(v interface{}) string {
			switch t := v.(type) {
			case string:
				return t
			case float64:
				return fmt.Sprintf("%.0f", t)
			default:
				return fmt.Sprintf("%v", v)
			}
		}

		if uid, exists := claims["user_id"]; exists {
			c.Set("user_id", toStr(uid))
		}
		if email, exists := claims["email"]; exists {
			c.Set("email", toStr(email))
		}
		if role, exists := claims["role"]; exists {
			c.Set("role", toStr(role))
		}

		c.Next()
	}
}
