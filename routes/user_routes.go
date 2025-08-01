package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mferdian/golang_boiller_plate/controller"
	"github.com/mferdian/golang_boiller_plate/middleware"
	"github.com/mferdian/golang_boiller_plate/service"
)

func UserRoutes(
	r *gin.Engine,
	userController controller.IUserController,
	jwtService service.InterfaceJWTService,
) {
	user := r.Group("/api/users")
	user.Use(middleware.Authentication(jwtService))

	// --- User Routes ---
	user.PATCH("/update-profile/:id", userController.UpdateUser)
	user.GET("/get-detail-user/:id", userController.GetUserByID)
	user.DELETE("/delete-profile/:id", userController.DeleteUser)
	
}
