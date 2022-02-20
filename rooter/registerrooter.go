package rooter

import (
	"APIGenerator/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func initRegisterRooter(r *gin.Engine) {
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
