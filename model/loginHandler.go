package model

import (
	"fmt"
	"github.com/gin-gonic/gin"
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
		loginByJson(c)
	case "application/x-www-form-urlencoded":
		loginByForm(c)
	}
}

func loginByForm(c *gin.Context) {
	var info Login
	if err := c.Bind(&info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !checkInfo(info.User, info.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"status": "304"})
		return
	}
	c.SetCookie("userName", info.User, 60, "", "localhost", false, true)
	c.SetCookie("password", info.Password, 60, "", "localhost", false, true)
	fmt.Println("User Login")
	fmt.Println("UserName--" + info.User + "\nPassword--" + info.Password)
	c.JSON(http.StatusOK, gin.H{"status": "200"})
}

func loginByJson(c *gin.Context) {
	var json Login
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !checkInfo(json.User, json.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"status": "304"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "200"})
}

func checkInfo(userName string, password string) bool {
	isVaild := false
	if userName == "root" && password == "admin" {
		isVaild = true
	}
	return isVaild
}
