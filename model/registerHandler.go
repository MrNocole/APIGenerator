package model

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
)

func CheckEmail(c *gin.Context) {
	//var info register1Info
	if code := c.PostForm("captcha"); code != "" {
		if checkCaptcha(code, c) {
			fmt.Println("captcha right")
			email := c.PostForm("email")
			emailVerifyCode := rand.Intn(9000) + 1000
			if err := SendVerify(email, fmt.Sprintf("%d", emailVerifyCode)); err != nil {
				fmt.Println(err)
				c.Abort()
			}
		} else {
			fmt.Println("captcha wrong!")
			c.Abort()
		}
	}
}
