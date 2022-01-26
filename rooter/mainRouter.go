package rooter

import (
	"APIGenerator/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(model.SessionDefault("regular"))
	r.LoadHTMLGlob("view/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{"registerweb": "47.93.212.155:8000/register-1"})
	})
	/*
		Login Handler
	*/
	{
		r.POST("/login", model.LoginHandler)
		r.GET("/login", func(c *gin.Context) {
			c.HTML(http.StatusOK, "login.html", gin.H{"registerweb": "47.93.212.155:8000/register-1"})
		})
		r.GET("/captcha", func(c *gin.Context) {
			model.Captcha(c, 4)
		})
	}
	/*
		Register handler
	*/
	{
		r.GET("/register-1", model.Session("regular", 180, "emailverify"), func(c *gin.Context) {
			c.HTML(http.StatusOK, "register-1.html", nil)
		})
		r.POST("/register-1", model.CheckEmail)
	}

	{
		r.GET("/home", model.UserCookieCheck, model.HomeHandler)
	}

	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"errorCode": "404"})
	})
	return r
}

func helloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello World!",
	})
}
