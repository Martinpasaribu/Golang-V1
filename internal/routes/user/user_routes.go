package routes

import (
	"project-name/internal/controllers"
	"project-name/internal/services"
	"project-name/internal/repositories"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Initialize dependencies
	userRepo := repositories.NewUserRepository()
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	// Routes
	api := r.Group("/api/v1")
	{
		api.POST("/users", userController.Register)
	}

	return r
}