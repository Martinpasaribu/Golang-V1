package routes

import (
	"log"

	blogController "github.com/Martinpasaribu/Golang-V1/internal/controllers/blog"
	userController "github.com/Martinpasaribu/Golang-V1/internal/controllers/user"
	"github.com/gin-gonic/gin"
)

type RouteConfig struct {
	UserController *userController.UserController
	BlogController *blogController.BlogController
	// Tambahkan controller lain di sini
}

func SetupRoutes(router *gin.Engine, config RouteConfig) {
	// Base route
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "API is running",
			"routes": []string{
				"POST /api/v1/users",
				"POST /api/v1/blogs",
				"GET  /api/v1/blogs/:id",
			},
		})
	})
	
	// API v1 routes
	apiV1 := router.Group("/api/v1")
	{
		if config.UserController != nil {
			RegisterUserRoutes(apiV1, config.UserController)
		} else {
			log.Println("⚠️ UserController not initialized. Skipping user routes.")
		}
		
		if config.BlogController != nil {
			RegisterBlogRoutes(apiV1, config.BlogController)
		} else {
			log.Println("⚠️ BlogController not initialized. Skipping blog routes.")
		}
	}
}