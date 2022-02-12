package rooter

import (
	"APIGenerator/model"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	_, err := model.Init()
	if err != nil {
		fmt.Println("数据库初始化失败!")
	} else {
		fmt.Println("数据库初始化成功!")
	}
	err = model.InitDocumentation()
	if err != nil {
		fmt.Println("仓库初始化失败!")
	} else {
		fmt.Println("仓库初始化成功!")
	}
	r.Use(model.SessionDefault("regular"))
	r.LoadHTMLGlob("view/*")
	NewUserInfoChan := make(chan *model.NewUserInfoInMysql, 10)
	//util.RegisterServer(NewUserInfoChan)
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{"registerweb": "localhost:8000/register-1"})
	})
	/*
		Login Handler
	*/
	{
		r.POST("/login", model.LoginHandler)
		r.GET("/login", func(c *gin.Context) {
			c.HTML(http.StatusOK, "login.html", gin.H{"registerweb": "localhost:8000/register-1"})
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
		r.GET("/register-2", func(c *gin.Context) {
			c.HTML(http.StatusOK, "register-2.html", nil)
		})
		r.POST("/register-2", func(c *gin.Context) {
			userInfo := model.CheckUsername(c)
			NewUserInfoChan <- &userInfo
		})
	}

	{
		r.GET("/home", model.UserCookieCheck, model.HomeHandler)
		r.POST("/upload", model.FileUploadHandler)
		r.GET("/upload", func(c *gin.Context) {
			c.HTML(http.StatusOK, "upload.html", nil)
		})
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
