package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mferdian/golang_boiller_plate/controller"
)

func PublicRoutes(r *gin.Engine, userController controller.IUserController) {
	public := r.Group("/api")
	public.POST("/register", userController.Register)
	public.POST("/login", userController.Login)
}
