package userController

import (
	"net/http"

	"github.com/Martinpasaribu/Golang-V1/internal/models"
	userService "github.com/Martinpasaribu/Golang-V1/internal/services/user"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	service userService.UserService
}

func NewUserController(service userService.UserService) *UserController {
	return &UserController{service: service}
}

func (c *UserController) Register(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := c.service.RegisterUser(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Hilangkan password di response
	result.Password = ""
	ctx.JSON(http.StatusCreated, result)
}

func (c *UserController) Login(ctx *gin.Context) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	
	if err := ctx.ShouldBindJSON(&loginData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.service.LoginUser(loginData.Email, loginData.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Generate JWT token (akan kita implementasikan nanti)
	token, err := generateToken(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	// Hilangkan password di response
	user.Password = ""
	ctx.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": token,
	})
}

// Helper function untuk generate JWT
func generateToken(user *models.User) (string, error) {
	// Implementasi nyata menggunakan JWT
	return "dummy-token", nil
}