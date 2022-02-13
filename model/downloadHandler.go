package model

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"sync"
)

var md5lock = make(map[string]*sync.Mutex)

func DownloadByAPI(c *gin.Context) {
	filename := c.Param("filename")
	uuid, err := c.Cookie("uuid")
	if err != nil {
		fmt.Println("未知访问者")
		c.Redirect(http.StatusTemporaryRedirect, "/404")
	}
	md5, err := getFileMd5ByUserFile(uuid, filename)
	fmt.Println("md5:", md5)
	if err != nil {
		fmt.Println(err)
		c.Redirect(http.StatusTemporaryRedirect, "/404")
	}
	filePath := "store/" + md5
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	GetMd5Lock(md5).Lock()
	defer GetMd5Lock(md5).Unlock()
	err = os.Rename(filePath, "store/"+filename)
	defer os.Rename("store/"+filename, "/store/"+md5)
	if err != nil {
		fmt.Println("仓库文件操作失败", err)
		c.Redirect(http.StatusTemporaryRedirect, "/404")
		return
	}
	c.File("store/" + filename)
}

func GetMd5Lock(md5 string) *sync.Mutex {
	if _, ok := md5lock[md5]; !ok {
		md5lock[md5] = new(sync.Mutex)
	}
	return md5lock[md5]
}

func getFileMd5ByUserFile(uuid, file string) (string, error) {
	ownerinfo, err := GetOwnerInfo(uuid)
	defer GetOwnerLock(uuid).Unlock()
	if err != nil {
		fmt.Println("用户信息仓库信息没有找到")
		return "", err
	}
	fileNames, fileMd5s := ownerinfo.FileName, ownerinfo.MD5
	for i, filename := range fileNames {
		if filename == file {
			return fileMd5s[i], nil
		}
	}
	return "", errors.New(file + "文件不存在")
}
