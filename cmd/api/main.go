package routes

import (
	"github.com/gin-gonic/gin" // Ganti dengan framework yang Anda pakai
	"github.com/Martinpasaribu/Golang-V1/internal/controllers/user"
	"github.com/Martinpasaribu/Golang-V1/internal/services"
	"github.com/Martinpasaribu/Golang-V1/internal/repositories"
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