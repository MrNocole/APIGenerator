package rooter

import (
	"APIGenerator/model"
	"APIGenerator/util"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

// redis 连接池
var redisPool = model.InitRedis()

// 当前链接
var link = util.GetUrl()

// NewUserInfoChan 新用户注册channel
var NewUserInfoChan = make(chan *util.RegisterPostFrom, 10)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	_, err := model.InitSQL()
	fmt.Println(link)
	if err != nil {
		fmt.Println("数据库初始化失败!")
	} else {
		fmt.Println("数据库初始化成功!")
	}
	err = model.InitStore()
	if err != nil {
		fmt.Println("仓库初始化失败!")
	} else {
		fmt.Println("仓库初始化成功!")
	}

	r.Use(model.SessionDefault("regular"))
	r.LoadHTMLGlob("view/*")
	// 初始化 store 相关的路由
	initStoreRooter(r)
	// 初始化 登录注册相关的路由
	initLoginRooter(r)
	initRegisterRooter(r)
	// 注册服务启动
	go util.RegisterServer(NewUserInfoChan)
	//util.RegisterServer(NewUserInfoChan)

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{"link": link, "registerweb": "/register-1"})
	})

	r.GET("/404", func(c *gin.Context) {
		c.HTML(http.StatusOK, "error.html", gin.H{"errorCode": "404", "info": "您访问的页面不存在"})
	})
	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"errorCode": "404", "info": "没有这个页面"})
	})
	return r
}
