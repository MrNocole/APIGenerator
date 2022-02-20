package rooter

import (
	"APIGenerator/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func initLoginRooter(r *gin.Engine) {
	r.POST("/login", model.LoginHandler)
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{"link": link, "registerweb": "/register-1"})
	})
	r.GET("/captcha", func(c *gin.Context) {
		model.Captcha(c, 4)
	})
}
