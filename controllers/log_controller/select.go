package log_controller

import (
	"net/http"
	"sunhost/model"

	"github.com/gin-gonic/gin"
)

func (lc *LogController) Select(ctx *gin.Context) {
    username := ctx.Query("username")
    if username == "" {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Username required"})
        return
    }

    var uLog model.UserLog
    uLogs, err := uLog.GetByUsername(username)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"logs": uLogs})
}
