package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/Martinpasaribu/Golang-V1/internal/config"
	blogController "github.com/Martinpasaribu/Golang-V1/internal/controllers/blog"
	userController "github.com/Martinpasaribu/Golang-V1/internal/controllers/user"
	blogRepository "github.com/Martinpasaribu/Golang-V1/internal/repositories/blog"
	userRepository "github.com/Martinpasaribu/Golang-V1/internal/repositories/user"
	blogService "github.com/Martinpasaribu/Golang-V1/internal/services/blog"
	userService "github.com/Martinpasaribu/Golang-V1/internal/services/user"
	"github.com/Martinpasaribu/Golang-V1/internal/routes" // Import package routes
	"github.com/gin-gonic/gin"
)

func main() {


	// 1. Initialize MongoDB
	if err := config.InitMongoDB(); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	defer config.CloseMongoDB()

	config.InitRedis()  
	config.GetImageKitKeys()

	// 2. Setup Router
	r := SetupRouter()

	// 3. Start Server
	log.Println("🚀 Server running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// ✅ Tambahkan konfigurasi CORS di sini
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000","http://localhost:3001"}, // URL Next.js
		AllowMethods:     []string{"POST", "GET", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	// ====================== USER ======================
	userRepo := userRepository.NewUserRepository(config.GetDB())
	userSvc := userService.NewUserService(userRepo)
	userCtrl := userController.NewUserController(userSvc)

	// ====================== BLOG ======================
	blogRepo := blogRepository.NewBlogRepository(config.GetDB())
	blogSvc := blogService.NewBlogService(blogRepo)
	blogCtrl := blogController.NewBlogController(blogSvc)

	// Setup routes
	routes.SetupRoutes(r, routes.RouteConfig{
		UserController: userCtrl,
		BlogController: blogCtrl,
	})

	return r
}
