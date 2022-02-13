package model

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"net/http"
)

type Login struct {
	User     string `form:"username" json:"user" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func LoginHandler(c *gin.Context) {
	postContent := c.ContentType()
	fmt.Println(postContent)
	switch postContent {
	case "application/json":
		//loginByJson(c)
	case "application/x-www-form-urlencoded":
		loginByForm(c)
	}
}

func loginByForm(c *gin.Context) {
	var info Login
	if code := c.PostForm("captcha"); code != "" {
		fmt.Println("captcha checking...")
		if checkCaptcha(code, c) {
			fmt.Println("captcha right!")
			if err := c.Bind(&info); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			fmt.Println(info)
			res, uuid, err := checkUserInfo(info.User, info.Password)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if !res {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "304",
					"message": "password error",
				})
				return
			}
			c.SetCookie("uuid", uuid, 86400, "", "localhost", false, true)
			c.SetCookie("userName", info.User, 86400, "", "localhost", false, true)
			c.SetCookie("password", info.Password, 86400, "", "localhost", false, true)
			fmt.Println("User Login")
			fmt.Println("UserName--" + info.User + "\nPassword--" + info.Password)
			c.Redirect(302, "/home")
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"status": 400, "Error": "input captcha"})
	}

}

//
//func loginByJson(c *gin.Context) {
//	var json Login
//	if err := c.ShouldBindJSON(&json); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//	if  checkUserInfo(json.User, json.Password) {
//		c.JSON(http.StatusBadRequest, gin.H{"status": "304"})
//		return
//	}
//	c.JSON(http.StatusOK, gin.H{"status": "200"})
//}

func checkUserInfo(userName string, password string) (bool, string, error) {
	if userName != " " {
		sqlConfig := loadSQLConfig()
		db, err := sqlx.Open("mysql", sqlConfig.User+":"+sqlConfig.Password+"@tcp("+sqlConfig.Host+")/"+sqlConfig.Database)
		if err != nil {
			return false, "", err
		}
		crtPwd, uuid, err := SelectPasswordAndUUidByUserName(db, userName)
		if err != nil {
			return false, uuid, err
		}
		if crtPwd == password {
			return true, uuid, nil
		} else {
			return false, uuid, errors.New("password error")
		}
	} else {
		return false, "", errors.New("input username")
	}
}
