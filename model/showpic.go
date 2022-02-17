package model

import (
	"APIGenerator/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func GetPic(c *gin.Context) {
	uuid := c.Param("uuid")
	link := util.GetUrl()
	uuidCookie, err := c.Cookie("uuid")
	if err != nil || uuid != uuidCookie {
		fmt.Println("未知访问者")
		c.Redirect(http.StatusTemporaryRedirect, "/404")
	}
	picName := c.Param("name")
	md5, err := getFileMd5ByUserFile(uuid, picName)
	if err != nil || md5 == "" {
		fmt.Println(err)
		c.Redirect(http.StatusTemporaryRedirect, "/404")
		return
	}
	filePath := "store/" + md5
	newfileName := "store/" + picName
	GetMd5Lock(md5).Lock()
	err = os.Rename(filePath, newfileName)
	if err != nil {
		fmt.Println("仓库文件操作失败", err)
		c.Redirect(http.StatusTemporaryRedirect, "/404")
		return
	}
	c.HTML(http.StatusOK, "picshow.html", gin.H{
		"link":     link,
		"uuid":     "store",
		"title":    picName,
		"filename": picName,
	})
	go resumeFile(md5, picName)
}
