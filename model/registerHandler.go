package model

import (
	"APIGenerator/util"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func RandomCode() string {
	return strconv.Itoa(int((time.Now().Unix() ^ 100000) % 10000))
}

func CheckEmail(c *gin.Context) {
	//var info register1Info
	if code := c.PostForm("captcha"); code != "" {
		if checkCaptcha(code, c) {
			fmt.Println("captcha right")
			email := c.PostForm("email")
			emailVerifyCode := RandomCode()
			fmt.Println(emailVerifyCode)
			go func(emailVerifyCode string) {
				if err := SendVerify(email, emailVerifyCode); err != nil {
					fmt.Println(err)
				}
			}(emailVerifyCode)
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

func CheckUsername(c *gin.Context) *util.RegisterPostFrom {
	userName := c.PostForm("user")
	passWord := c.PostForm("password")
	fmt.Println(userName + " " + passWord)
	if userName == "" || passWord == "" {
		util.ErrorHtml(c, strconv.Itoa(http.StatusBadGateway), "用户名或密码不能为空")
		return nil
	}
	info := util.RegisterPostFrom{
		UserName: userName,
		Password: passWord,
	}
	info.UserName = userName
	info.Password = passWord
	session := sessions.Default(c)
	verifyCode := c.PostForm("verify")
	verifyCodeInSession := session.Get("emailVerifyCode").(string)
	if verifyCode != verifyCodeInSession {
		util.ErrorHtml(c, strconv.Itoa(http.StatusBadRequest), "验证码错误")
	}
	email := session.Get("email")
	info.Email = email.(string)
	return &info
}
