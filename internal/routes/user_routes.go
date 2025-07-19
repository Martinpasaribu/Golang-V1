package routes

import (
	userController "github.com/Martinpasaribu/Golang-V1/internal/controllers/user"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.RouterGroup, userCtrl *userController.UserController) {
	userGroup := router.Group("/users")
	{
		userGroup.POST("", userCtrl.Register)
	}
}