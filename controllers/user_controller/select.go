package user_controller

import (
	"net/http"
	"sunhost/controllers/log_controller"
	"sunhost/model"

	"github.com/gin-gonic/gin"
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
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if req.Password != user.Password {
		logC := log_controller.LogController{}
        ctx.Set("action", "Failed login")
        ctx.Set("username", req.Username)
        logC.Create(ctx)

		return
	}
	ctx.Set("action", "User loggin")
	ctx.Set("username", user.Username)

	logC.Create(ctx)
	ctx.JSON(http.StatusOK, gin.H{"user": user})
}
