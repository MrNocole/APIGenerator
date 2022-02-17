package model

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func GetJson(c *gin.Context) {
	uuid := c.Param("uuid")
	uuidCookie, err := c.Cookie("uuid")
	if err != nil || uuid != uuidCookie {
		fmt.Println("未知访问者")
		c.Redirect(http.StatusTemporaryRedirect, "/404")
	}
	jsonName := c.Param("name")
	md5, err := getFileMd5ByUserFile(uuid, jsonName)
	if err != nil || md5 == "" {
		fmt.Println(err)
		c.Redirect(http.StatusTemporaryRedirect, "/404")
		return
	}
	file, err := getFileByMd5(md5)
	if err != nil {
		fmt.Println(err)
		c.Redirect(http.StatusTemporaryRedirect, "/404")
		return
	}
	var data []byte
	data, err = io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		c.Redirect(http.StatusTemporaryRedirect, "/404")
		return
	}
	c.AsciiJSON(http.StatusOK, gin.H{
		"user": uuid,
		"data": string(data),
	})
}
