package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mferdian/golang_boiller_plate/command"
	"github.com/mferdian/golang_boiller_plate/config/database"
	"github.com/mferdian/golang_boiller_plate/controller"
	"github.com/mferdian/golang_boiller_plate/logging"
	"github.com/mferdian/golang_boiller_plate/middleware"
	"github.com/mferdian/golang_boiller_plate/repository"
	"github.com/mferdian/golang_boiller_plate/routes"
	"github.com/mferdian/golang_boiller_plate/service"
)

func main() {
	// ==== Set up logger ====
	logging.SetUpLogger()
	logging.Log.Info("Logger initialized")

	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// DB
	db := database.SetUpPostgreSQLConnection()
	defer database.ClosePostgreSQLConnection(db)

	// Seeder command
	if len(os.Args) > 1 {
		cmd.Command(db)
		return
	}

	var (
		jwtService = service.NewJWTService()

		userRepo       = repository.NewUserRepository(db)
		userService    = service.NewUserService(userRepo, jwtService)
		userController = controller.NewUserController(userService)
	)

	server := gin.Default()
	server.Use(middleware.CORSMiddleware())

	routes.PublicRoutes(server, userController)
	routes.AdminRoutes(server, userController, jwtService)
	routes.UserRoutes(server, userController, jwtService)

	server.Static("/assets", "./assets")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	var serve string
	if os.Getenv("APP_ENV") == "localhost" {
		serve = "127.0.0.1:" + port
	} else {
		serve = ":" + port
	}

	if err := server.Run(serve); err != nil {
		log.Fatalf("error running server: %v", err)
	}
}
