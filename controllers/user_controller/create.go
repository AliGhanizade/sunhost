package user_controller

import (
	"net/http"
	"sunhost/controllers/log_controller"
	"sunhost/model"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (uc *UserController) Create(ctx *gin.Context) {
	var user model.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingUser model.User
	if err := existingUser.GetByUsername(user.Username); err == nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	if err := existingUser.GetByEmail(user.Email); err == nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}
	user.Password = string(hash)

	if err := user.Create(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	ctx.Set("action", "User Created")
	ctx.Set("username", user.Username)
	logC := log_controller.LogController{}
    logC.Create(ctx)

	
	ctx.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}
