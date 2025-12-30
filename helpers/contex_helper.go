package helpers

import (
	"context"

	"github.com/gin-gonic/gin"
)

func GetUserID(ctx context.Context) string {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		if val, exists := ginCtx.Get("id"); exists {
			if id, ok := val.(string); ok {
				return id
			}
		}
	}
	if val := ctx.Value("id"); val != nil {
		if id, ok := val.(string); ok {
			return id
		}
	}

	return ""
}

func GetUserRole(ctx context.Context) string {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		if val, exists := ginCtx.Get("role"); exists {
			if role, ok := val.(string); ok {
				return role
			}
		}
	}
	if val := ctx.Value("role"); val != nil {
		if role, ok := val.(string); ok {
			return role
		}
	}

	return ""
}
