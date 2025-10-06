package helpers

import (
	"context"

	"github.com/gin-gonic/gin"
)

// GetUserID mengambil user ID dari context (bisa *gin.Context atau context.Context)
func GetUserID(ctx context.Context) string {
	// jika context berasal dari gin.Context
	if ginCtx, ok := ctx.(*gin.Context); ok {
		if val, exists := ginCtx.Get("id"); exists {
			if id, ok := val.(string); ok {
				return id
			}
		}
	}

	// fallback untuk context biasa (kalau suatu saat kamu set manual)
	if val := ctx.Value("id"); val != nil {
		if id, ok := val.(string); ok {
			return id
		}
	}

	return ""
}

// GetUserRole mengambil role user dari context (bisa *gin.Context atau context.Context)
func GetUserRole(ctx context.Context) string {
	// jika context berasal dari gin.Context
	if ginCtx, ok := ctx.(*gin.Context); ok {
		if val, exists := ginCtx.Get("role"); exists {
			if role, ok := val.(string); ok {
				return role
			}
		}
	}

	// fallback untuk context biasa (kalau diset manual)
	if val := ctx.Value("role"); val != nil {
		if role, ok := val.(string); ok {
			return role
		}
	}

	return ""
}
