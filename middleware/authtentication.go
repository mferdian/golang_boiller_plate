package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mferdian/golang_boiller_plate/constants"
	"github.com/mferdian/golang_boiller_plate/logging"
	"github.com/mferdian/golang_boiller_plate/service"
	"github.com/mferdian/golang_boiller_plate/utils"
)

func Authentication(jwtService service.InterfaceJWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			logging.Log.Warn("Authorization header not found")
			res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_PROSES_REQUEST, constants.MESSAGE_FAILED_TOKEN_NOT_FOUND, nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			logging.Log.Warn("Authorization header format is invalid")
			res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_PROSES_REQUEST, constants.MESSAGE_FAILED_TOKEN_NOT_VALID, nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwtService.ValidateToken(tokenStr)
		if err != nil || !token.Valid {
			logging.Log.Warnf("Invalid token: %v", err)
			res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_PROSES_REQUEST, constants.MESSAGE_FAILED_TOKEN_NOT_VALID, nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		userID, err := jwtService.GetUserIDByToken(tokenStr)
		if err != nil {
			logging.Log.Warnf("Failed to extract user ID from token: %v", err)
			res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_PROSES_REQUEST, err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		role, err := jwtService.GetRoleByToken(tokenStr)
		if err != nil {
			logging.Log.Warnf("Failed to extract role from token: %v", err)
			res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_PROSES_REQUEST, err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		logging.Log.Infof("Authenticated request - UserID: %s, Role: %s", userID, role)

		ctx.Set("Authorization", tokenStr)
		ctx.Set("user_id", userID)
		ctx.Set("role", role)

		ctx.Next()
	}
}
