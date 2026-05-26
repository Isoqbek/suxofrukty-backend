package middleware

import (
	"net/http"
	"strings"

	"github.com/Isoqbek/suxofrukty-backend/pkg/jwtutil"
	"github.com/gin-gonic/gin"
)

func AdminOnly(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		claims, err := jwtutil.Parse(strings.TrimPrefix(header, "Bearer "), secret)
		if err != nil || claims.Role != "admin" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.Set("admin_id", claims.AdminID)
		c.Next()
	}
}
