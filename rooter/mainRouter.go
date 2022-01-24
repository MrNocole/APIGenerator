package rooter

import (
	"APIGenerator/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("view/*")
	r.GET("/", helloHandler)

	/*
		Login Handler
	*/
	{
		r.POST("/login", model.LoginHandler)
		r.GET("/login", func(c *gin.Context) {
			c.HTML(http.StatusOK, "login.html", nil)
		})
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
