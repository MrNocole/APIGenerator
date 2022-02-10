package model

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"os"
	"sync"
)

var jsonLock sync.Mutex

// FileJson 文件Json格式 (完整信息)
type FileJson struct {
	FileName string `json:"fileName"`
	MD5      string `json:"md5"`
	Size     int64  `json:"size"`
	count    int64
}

// OwnerInfo 用户Json格式
type OwnerInfo struct {
	UUid     string   `json:"uuid"`
	FileName []string `json:"fileName"`
	MD5      []string `json:"md5"`
}

// FileMD5 文件Json格式，存硬盘
type FileMD5 struct {
	MD5   string `json:"md5"`
	Count int64  `json:"count"`
}

// FileUploadHandler 处理上传post，并存到中转文件夹。
func FileUploadHandler(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(400, gin.H{
			"message": "上传失败",
		})
		return
	}
	filesJson := make([]FileJson, 1)
	files := form.File["files"]
	for _, file := range files {
		md5h := md5.New()
		err := c.SaveUploadedFile(file, "filetransit/"+file.Filename)
		if err != nil {
			c.JSON(400, gin.H{
				"file":    file.Filename,
				"size":    file.Size,
				"message": "上传失败",
			})
		}
		pfile, err := os.Open("filetransit/" + file.Filename)
		if err != nil {
			c.JSON(400, gin.H{
				"file":    file.Filename,
				"size":    file.Size,
				"message": "上传成功，解析失败",
			})
		}
		io.Copy(md5h, pfile)
		md5Code := hex.EncodeToString(md5h.Sum(nil))
		go transitFile(pfile, md5Code)
		filetmp := FileJson{
			FileName: file.Filename,
			MD5:      md5Code,
			Size:     file.Size,
		}
		filesJson = append(filesJson, filetmp)
	}
	fmt.Println(filesJson)
	jsonByte, err := json.Marshal(filesJson)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "上传失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "上传成功",
		"data":    string(jsonByte),
	})
}

// 验证文件对应 md5 是否已经存在，并更新文件json
// transitFile 这个功能有严重性能问题，每次处理文件都扫描全部文件的Json，在解析、遍历过程都 O(n) ，如果出现性能瓶颈优先改这里，目前的思路大概是二分Json文件/用树结构组织md码
func transitFile(file *os.File, md5 string) {
	jsonLock.Lock()
	allFiles, err := ioutil.ReadFile("Documentation.json")
	if err != nil {
		fmt.Println("读取文件记录失败")
		jsonLock.Unlock()
		return
	}
	var filesJson []FileMD5
	json.Unmarshal(allFiles, &filesJson)
	for i, filetmp := range filesJson {
		if filetmp.MD5 == md5 {
			filesJson[i].Count++
			fmt.Println("该文件已经存储，更新json")
			err := os.Remove(file.Name())
			if err != nil {
				fmt.Println("删除中转文件失败", err)
				return
			}
			newJson, err := json.Marshal(filesJson)
			if err != nil {
				fmt.Println("Json更新失败")
				jsonLock.Unlock()
				return
			}
			ioutil.WriteFile("Documentation.json", newJson, 777)
			jsonLock.Unlock()
			fmt.Println("Json更新完毕")
			return
		}
	}
	fmt.Println("该文件没有存储，中转到仓库")
	err = os.Rename(file.Name(), "store/"+md5)
	if err != nil {
		fmt.Println("中转文件移动失败", err)
		return
	}
	var newFile = FileMD5{
		MD5:   md5,
		Count: 1,
	}
	filesJson = append(filesJson, newFile)
	newJson, err := json.Marshal(filesJson)
	if err != nil {
		fmt.Println("json更新失败")
		jsonLock.Unlock()
		return
	}
	ioutil.WriteFile("Documentation.json", newJson, 777)
	fmt.Println("Json 更新完毕!")
	jsonLock.Unlock()
}

// checkFileExist 提前优化，解耦检查文件和处理
func checkFileExist(md string) (bool, error) {
	jsonLock.Lock()
	allFiles, err := ioutil.ReadFile("Documentation.json")
	if err != nil {
		fmt.Println("读取文件记录失败")
		return false, err
	}
	var filesJson []FileMD5
	json.Unmarshal(allFiles, &filesJson)
	for _, filetmp := range filesJson {
		if filetmp.MD5 == md {
			fmt.Println("该文件已经存储")
			return true, nil
		}
	}
	fmt.Println("该文件没有存储，中转到仓库")
	return false, nil
}
