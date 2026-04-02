package log_controller

import (
	"log"
	"strings"
	"sunhost/model"

	"github.com/gin-gonic/gin"
)

func (lc *LogController) Create(ctx *gin.Context) {
	action := ctx.GetString("action")
	username := ctx.GetString("username")
	if username == "" {
		return
	}
	if action == "" {
		return
	}
	var uLog model.UserLog
	count, err := uLog.CheckCountByUsername(username)
	if err != nil {
		return
	}

	uLog.Username = username
	uLog.IPAddress = ctx.ClientIP()
	uaString := ctx.Request.UserAgent()
	uLog.SystemInfo = parseUserAgent(uaString)
	uLog.Action = action
	err = uLog.Create()
	if err != nil {
		log.Println(err)
		return
	}
	if count >= 20 {
		err = uLog.DeleteOnceByUsername()
		if err != nil {
			log.Println(err)

			return
		}
	}

}

func parseUserAgent(ua string) string {
	os := "Unknown OS"
	if strings.Contains(ua, "iPhone") {
		os = "iOS"
	} else if strings.Contains(ua, "Android") {
		os = "Android"
	} else if strings.Contains(ua, "Windows") {
		os = "Windows"
	} else if strings.Contains(ua, "Macintosh") {
		os = "macOS"
	}

	browser := "Unknown Browser"
	if strings.Contains(ua, "Edg/") {
		browser = "Edge"
	} else if strings.Contains(ua, "Chrome") {
		browser = "Chrome"
	} else if strings.Contains(ua, "Safari") {
		browser = "Safari"
	} else if strings.Contains(ua, "Firefox") {
		browser = "Firefox"
	}

	return os + " - " + browser
}
