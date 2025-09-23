package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "error", "message": "Role tidak ditemukan"})
			return
		}

		userRole, ok := roleVal.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "error", "message": "Format role tidak valid"})
			return
		}

		for _, r := range roles {
			if r == "*" {
				c.Next()
				return
			}
			if userRole == r {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "error", "message": "Akses ditolak, role tidak sesuai"})
	}
}
