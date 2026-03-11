package routes

import (
	"html/template"
	"log"
	"net/http"
	"sunhost/controllers/log_controller"
	"sunhost/controllers/system_stat_cotroller"
	"sunhost/controllers/user_controller"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	tmpl := template.New("")
	tmpl, err := tmpl.ParseGlob("./public/*.html")
	if err != nil {
		log.Panic(err)
	}
	tmpl, err = tmpl.ParseGlob("./public/panel/*.html")
	if err != nil {
		log.Panic(err)
	}
	r.SetHTMLTemplate(tmpl)
	userC := user_controller.UserController{}
	logC := log_controller.LogController{}
	sysC := system_stat_cotroller.SystemController{}
	api := r.Group("/api")
	{
		users := api.Group("/users")
		{
			users.POST("/register", userC.Create)
			users.POST("/login", userC.Select)
			users.GET("/logs", logC.Select)
		}
		api.GET("/system/stream", sysC.LiveStream)

		api.GET("/system/summary", sysC.GetSystemSummary)
	}
	r.Static("/host/assets/", "./public/assets/")
	hostGroup := r.Group("/host")
	{
		hostGroup.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", gin.H{})
		})
		hostGroup.GET("/products", func(c *gin.Context) {
			c.HTML(http.StatusOK, "products.html", gin.H{})
		})
		hostGroup.GET("/blog", func(c *gin.Context) {
			c.HTML(http.StatusOK, "blog.html", gin.H{})
		})
		hostGroup.GET("/about", func(c *gin.Context) {
			c.HTML(http.StatusOK, "about.html", gin.H{})
		})
		hostGroup.GET("/contact", func(c *gin.Context) {
			c.HTML(http.StatusOK, "contact.html", gin.H{})
		})
		hostGroup.GET("/login", func(c *gin.Context) {
			c.HTML(http.StatusOK, "login.html", gin.H{})
		})
	}

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/host/")
	})

	r.Static("/panel/assets", "./public/assets")

	panelGroup := r.Group("/panel")
	{
		panelGroup.GET("/", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "panel.html", gin.H{})
		})
		panelGroup.GET("/users", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "users.html", gin.H{})
		})
		panelGroup.GET("/servers", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "servers.html", gin.H{})
		})
		panelGroup.GET("/databases", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "databases.html", gin.H{})
		})
		panelGroup.GET("/system-info", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "system_info.html", gin.H{})
		})
	}

	return r
}
