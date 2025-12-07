package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mferdian/golang_boiller_plate/constants"
	"github.com/mferdian/golang_boiller_plate/controller"
	"github.com/mferdian/golang_boiller_plate/middleware"
	"github.com/mferdian/golang_boiller_plate/service"
)

func AdminRoutes(r *gin.Engine, userController controller.IUserController,
	jwtService service.InterfaceJWTService) {
	admin := r.Group("/api/users")
	admin.Use(middleware.Authentication(jwtService))
	admin.Use(middleware.AuthorizeRole(constants.ENUM_ROLE_ADMIN))

	// User management
	admin.POST("", userController.CreateUser)
	admin.GET("", userController.GetAllUser)
}
