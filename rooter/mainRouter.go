package rooter

import (
	"APIGenerator/model"
	"APIGenerator/util"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	_, err := model.Init()
	link := util.GetUrl()
	fmt.Println(link)
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
	NewUserInfoChan := make(chan *util.RegisterPostFrom, 10)
	go util.RegisterServer(NewUserInfoChan)
	//util.RegisterServer(NewUserInfoChan)
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{"link": link, "registerweb": "/register-1"})
	})
	/*
		Login Handler
	*/
	{
		r.POST("/login", model.LoginHandler)
		r.GET("/login", func(c *gin.Context) {
			c.HTML(http.StatusOK, "login.html", gin.H{"link": link, "registerweb": "/register-1"})
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
			c.HTML(http.StatusOK, "register-1.html", gin.H{
				"link": link,
			})
		})
		r.POST("/register-1", model.CheckEmail)
		r.GET("/register-2", func(c *gin.Context) {
			c.HTML(http.StatusOK, "register-2.html", gin.H{
				"link": link,
			})
		})
		r.POST("/register-2", func(c *gin.Context) {
			userInfo := model.CheckUsername(c)
			fmt.Println(userInfo)
			NewUserInfoChan <- userInfo
		})
	}

	// store Handler
	{
		r.GET("/download/:uuid/:filename", model.DownloadByAPI)
		r.GET("/home", model.UserCookieCheck, model.HomeHandler)
		r.POST("/upload", model.FileUploadHandler)
		r.GET("/check/:uuid/:name", model.CheckHandler)
		r.GET("/json/:uuid/:name", model.GetJson)
	}
	r.GET("/404", func(c *gin.Context) {
		c.HTML(http.StatusOK, "error.html", gin.H{"errorCode": "404", "info": "您访问的页面不存在"})
	})
	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"errorCode": "404", "info": "没有这个页面"})
	})
	return r
}
