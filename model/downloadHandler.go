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

// DownloadByAPI 根据文件名获取md5然后下载
func DownloadByAPI(c *gin.Context) {
	// 文件名和用户上传/修改后的文件名有关，对应一个唯一md5
	filename := c.Param("filename")
	uuid, err := c.Cookie("uuid")
	if err != nil {
		fmt.Println("未知访问者")
		c.Redirect(http.StatusTemporaryRedirect, "/404")
	}
	// 拿到用户的文件列表遍历找到对应其文件名的md5
	md5, err := getFileMd5ByUserFile(uuid, filename)
	fmt.Println("md5:", md5)
	if err != nil || md5 == "" {
		fmt.Println(err)
		c.Redirect(http.StatusTemporaryRedirect, "/404")
		return
	}
	// 获取文件的仓库地址
	filePath := "store/" + md5
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")

	// 锁住md5对应文件，等待修改文件名并恢复
	GetMd5Lock(md5).Lock()
	// 修改文件名为用户需要的
	err = os.Rename(filePath, "store/"+filename)
	if err != nil {
		fmt.Println("仓库文件操作失败", err)
		c.Redirect(http.StatusTemporaryRedirect, "/404")
		return
	}
	// 下载文件
	c.File("store/" + filename)
	// 恢复文件名，释放文件锁
	go resumeFile(md5, filename)
}

func resumeFile(md5, filename string) {
	err := os.Rename("store/"+filename, "store/"+md5)
	if err != nil {
		fmt.Println("仓库文件恢复失败", err)
	}
	GetMd5Lock(md5).Unlock()
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
	return "", errors.New("文件不存在")
}
