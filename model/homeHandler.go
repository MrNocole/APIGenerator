package model

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HomeHandler(c *gin.Context) {
	c.JSON(200, gin.H{"Status": "200"})
}

func UserCookieCheck(c *gin.Context) {
	fmt.Println("MiddleWare begin...")
	fmt.Println(c.Cookie("userName"))

	if userName, err := c.Cookie("userName"); err == nil {
		password, _ := c.Cookie("password")
		fmt.Println("User found!--" + userName)
		if userName == "root" && password == "admin" {
			c.Next()
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Login info error!"})
		c.Abort()
		return
	}
	fmt.Println("Cookie is not found")
	c.Abort()
}
