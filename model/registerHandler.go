package model

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"math/rand"
	"strconv"
	"time"
)

func RandomCode() string {
	return strconv.Itoa(int(time.Now().Unix() ^ 100000%10000))
}

func CheckEmail(c *gin.Context) {
	//var info register1Info
	if code := c.PostForm("captcha"); code != "" {
		if checkCaptcha(code, c) {
			fmt.Println("captcha right")
			email := c.PostForm("email")
			emailVerifyCode := rand.Intn(9000) + 1000
			fmt.Println(emailVerifyCode)
			//if err := SendVerify(email, fmt.Sprintf("%d", emailVerifyCode)); err != nil {
			//	fmt.Println(err)
			//	c.Abort()
			//}
			session := sessions.Default(c)
			session.Set("email", email)
			session.Set("emailVerifyCode", emailVerifyCode)
			err := session.Save()
			if err != nil {
				return
			}
			c.Redirect(302, "/register-2")
		} else {
			fmt.Println("captcha wrong!")
			c.Abort()
		}
	}
}

func CheckUsername(c *gin.Context) (userInfo NewUserInfoInMysql) {

	return userInfo
}
