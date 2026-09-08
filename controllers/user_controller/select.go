package user_controller

import (
	"net/http"
	"sunhost/controllers/log_controller"
	"sunhost/model"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (uc *UserController) Select(ctx *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	logC := log_controller.LogController{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user model.User

	if err := user.GetByUsername(req.Username); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		ctx.Set("action", "Failed login")
		ctx.Set("username", req.Username)
		logC.Create(ctx)

		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}
	ctx.Set("action", "User logged in")
	ctx.Set("username", user.Username)

	logC.Create(ctx)
	ctx.JSON(http.StatusOK, gin.H{"user": user})
}
