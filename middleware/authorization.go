package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mferdian/golang_boiller_plate/logging" 
)

func AuthorizeRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, ok := c.Get("role")
		if !ok {
			logging.Log.Warn("Authorization failed: role not found in context")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "unauthorized: role not found",
			})
			return
		}

		roleStr := userRole.(string)
		for _, allowed := range allowedRoles {
			if strings.EqualFold(roleStr, allowed) {
				logging.Log.Infof("Authorized access: role=%s allowed", roleStr)
				c.Next()
				return
			}
		}

		logging.Log.Warnf("Forbidden access attempt: role=%s not in allowedRoles=%v", roleStr, allowedRoles)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  false,
			"message": "forbidden: you don't have access to this resource",
		})
	}
}
