package model

import (
	"APIGenerator/util"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"io"
	"io/ioutil"
	"os"
	"sync"
)

// 这个设计比较复杂，留了一堆问题等着设计。
// 流程：1.验证上传者信息 （目前只有UUID，且没有验证，以后要加入Token）
//		2.把上传的文件存到中转缓存
//		3.把中转缓存的文件存到仓库，更新仓库清单 （用文件MD5做的清单，清单有严重的性能问题，下面注释里有）
//		4.把上传成功的全部文件文件名和MD5对应写到上传者的Json文件里

// 仓库清单锁
var jsonLock sync.Mutex

// 用户清单锁
var lockMap = make(map[string]*sync.Mutex)

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

// FileMD5 仓库清单Json格式，存硬盘
type FileMD5 struct {
	MD5   string `json:"md5"`
	Count int64  `json:"count"`
}

func FileDeleteHandler(c *gin.Context, Conn redis.Conn) {
	filename := c.Param("filename")
	uuid := c.Param("uuid")
	jsonfileName := "fileownerinfo/" + uuid + ".json"
	lock := GetOwnerLock(uuid)
	lock.Lock()
	data, err := ioutil.ReadFile(jsonfileName)
	if err != nil {
		c.JSON(200, gin.H{
			"status":  "error",
			"message": "读取用户 json 失败",
		})
		lock.Unlock()
		return
	}
	var ownerInfo OwnerInfo
	err = json.Unmarshal(data, &ownerInfo)
	if err != nil {
		c.JSON(200, gin.H{
			"status":  "error",
			"message": "解析用户 json 失败",
		})
		lock.Unlock()
		return
	}
	var deleteMd5 string
	for i, v := range ownerInfo.FileName {
		if v == filename {
			ownerInfo.FileName = append(ownerInfo.FileName[:i], ownerInfo.FileName[i+1:]...)
			ownerInfo.MD5 = append(ownerInfo.MD5[:i], ownerInfo.MD5[i+1:]...)
			deleteMd5 = ownerInfo.MD5[i]
			break
		}
	}
	deletedData, err := json.Marshal(ownerInfo)
	if err != nil {
		c.JSON(200, gin.H{
			"status":  "error",
			"message": "修改用户 json 失败",
		})
		lock.Unlock()
		return
	}
	go DeleteMd5(deleteMd5)
	go util.RedisRemoveH(Conn, uuid, filename)
	err = os.WriteFile(jsonfileName, deletedData, 0666)
	if err != nil {
		c.JSON(200, gin.H{
			"status":  "error",
			"message": "写入用户 json 失败",
		})
	}
	lock.Unlock()
}

func DeleteMd5(md5 string) {
	filePath := "Documentation.json"
	jsonLock.Lock()
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("读取仓库清单失败", err)
		jsonLock.Unlock()
		return
	}
	var fileList []FileMD5
	err = json.Unmarshal(data, &fileList)
	if err != nil {
		fmt.Println("解析仓库清单失败", err)
		jsonLock.Unlock()
		return
	}
	for i := range fileList {
		if md5 == fileList[i].MD5 {
			fileList[i].Count--
			if fileList[i].Count == 0 {
				fileList = append(fileList[:i], fileList[i+1:]...)
				removeFileInStore(fileList[i].MD5)
			}
			break
		}
	}
	data, err = json.Marshal(fileList)
	if err != nil {
		fmt.Println("更新仓库解析清单失败", err)
		jsonLock.Unlock()
		return
	}
	err = os.WriteFile(filePath, data, 0666)
	if err != nil {
		fmt.Println("写入仓库清单失败", err)
	}
	jsonLock.Unlock()
}

func removeFileInStore(md5 string) {
	filePath := "store/" + md5
	err := os.Remove(filePath)
	if err != nil {
		fmt.Println(err)
	}
}

// FileUploadHandler 处理上传post
func FileUploadHandler(c *gin.Context, Conn redis.Conn) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(400, gin.H{
			"message": "上传失败",
		})
		return
	}
	// 上传完成的文件信息存到filesJson切片里
	filesJson := make([]FileJson, 1)
	files := form.File["files"]
	// 获取上传者信息（仅用 uuid 识别，以后加入Token验证）
	uuid, err := c.Cookie("uuid")
	if err != nil {
		c.JSON(400, gin.H{
			"message": "未知上传来源",
			"error":   err,
		})
		return
	}
	OwnerInfo, err := GetOwnerInfo(uuid)
	defer GetOwnerLock(uuid).Unlock()
	if err != nil {
		c.JSON(400, gin.H{
			"message": "用户信息获取失败",
			"error":   err,
		})
		return
	}

	// 用户新上传的文件名和对应md5
	wg := sync.WaitGroup{}
	newFileName := make([]string, len(files))
	newMd5Name := make([]string, len(files))

	for i, file := range files {
		wg.Add(2)
		md5h := md5.New()
		// 把用户上传的文件写到中转目录，中转目录下的文件名为用户上传的文件名，md5为文件内容的md5，可以考虑用这里做一个缓存，应对错误
		err := c.SaveUploadedFile(file, "filetransit/"+file.Filename)
		if err != nil {
			c.JSON(400, gin.H{
				"file":    file.Filename,
				"size":    file.Size,
				"message": "上传失败",
			})
			return
		}

		// 把传上来的文件内容读取到md5h里。
		pfile, err := os.Open("filetransit/" + file.Filename)
		if err != nil {
			c.JSON(400, gin.H{
				"file":    file.Filename,
				"size":    file.Size,
				"message": "上传成功，解析失败",
			})
			return
		}
		io.Copy(md5h, pfile)
		fileName := file.Filename
		md5Code := hex.EncodeToString(md5h.Sum(nil))

		// 编码md5完成后，把文件名和md5写到待更新切片里
		go updateSlice(newFileName, i, fileName, &wg)
		go updateSlice(newMd5Name, i, md5Code, &wg)
		// 这里 1.把文件从中转移动到仓库 2.更新仓库清单
		go transitFile(pfile, md5Code)
		go util.RedisInsertH(Conn, uuid, fileName, md5Code)
		// 调试用的
		filetmp := FileJson{
			FileName: file.Filename,
			MD5:      md5Code,
			Size:     file.Size,
		}
		filesJson = append(filesJson, filetmp)
	}
	// 调试
	fmt.Println(filesJson)
	jsonByte, err := json.Marshal(filesJson)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "上传失败",
		})
		return
	}

	//等切片处理完成后一起更新用户Json信息
	wg.Wait()
	OwnerInfo.FileName = append(OwnerInfo.FileName, newFileName...)
	OwnerInfo.MD5 = append(OwnerInfo.MD5, newMd5Name...)
	modifyedOwnerInfo, err := json.Marshal(OwnerInfo)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "上传成功，但用户信息json在内存中修改失败",
		})
		return
	}
	err = os.WriteFile("fileownerinfo/"+uuid+".json", modifyedOwnerInfo, 0666)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "上传成功，但用户信息json在服务器中c修改失败",
			"error":   err,
		})
		return
	}

	// 上传成功之后要更新用户在redis中的缓存。

	c.JSON(200, gin.H{
		"message": "上传成功",
		"data":    string(jsonByte),
	})
}

// 更新切片的指定位置
// 这里的设计是为了做并行，用户上传文件并行处理会导致很多问题，slice的并发安全很难保证，但是上传的文件是可以得到序号的，因此这里直接指定修改slice的指定序号
// 当然这个设计也是是有问题的，如果用户上传的多个文件中间出了问题，会导致最终用户json中存在一个空隙，但是这个空隙造成的危害现在看来不大。如果后期有需要，修复Json也是容易的
// 衡量了并发安全的开支和可能出现的问题，最终采用了这个方案。
func updateSlice(slice []string, index int, value string, wg *sync.WaitGroup) []string {
	if index < 0 || index > len(slice) {
		return slice
	}
	wg.Done()
	slice[index] = value
	return slice
}

// GetOwnerInfo 获取用户拥有文件的 Json，如果没有就新建。这里会取锁，但不会主动释放
func GetOwnerInfo(uuid string) (*OwnerInfo, error) {
	fileName := "fileownerinfo/" + uuid + ".json"
	fmt.Println(fileName)
	lock := GetOwnerLock(uuid)
	lock.Lock()
	var file *os.File
	var ownerInfo OwnerInfo
	if !util.CheckFileIsExist(fileName) {
		fmt.Println("first time upload! create a json for owner")
		var err error
		file, err = os.Create(fileName)
		if err != nil {
			fmt.Println("create file error:", err)
			return nil, err
		}
		ownerInfo = OwnerInfo{
			UUid:     uuid,
			FileName: []string{},
			MD5:      []string{},
		}
		initinfo, err := json.Marshal(ownerInfo)
		if err != nil {
			fmt.Println("初始化Json生成失败", err)
			os.Remove(fileName)
			return nil, err
		}
		_, err = file.Write(initinfo)
		if err != nil {
			fmt.Println("初始化Json写入失败", err)
			os.Remove(fileName)
			return nil, err
		}
		return &ownerInfo, nil
	} else {
		var err error
		file, err = os.Open(fileName)
		if err != nil {
			fmt.Println("open file error:", err)
			return nil, err
		}
	}
	buffer, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("read file error:", err)
		return nil, err
	}
	err = json.Unmarshal(buffer, &ownerInfo)
	if err != nil {
		fmt.Println("unmarshal error:", err)
		return nil, err
	}
	return &ownerInfo, nil
}

func GetOwnerLock(uuid string) *sync.Mutex {
	if _, ok := lockMap[uuid]; !ok {
		lockMap[uuid] = &sync.Mutex{}
	}
	return lockMap[uuid]
}

func InitStore() error {
	if !util.CheckFileIsExist("filetransit") {
		err := os.Mkdir("filetransit", os.ModePerm)
		if err != nil {
			return err
		}
		fmt.Println("初始化中转仓库")
	} else {
		fmt.Println("中转仓库已存在")
	}
	if !util.CheckFileIsExist("fileownerinfo") {
		err := os.Mkdir("fileownerinfo", os.ModePerm)
		if err != nil {
			return err
		}
		fmt.Println("初始化用户清单")
	} else {
		fmt.Println("用户清单已存在")
	}
	if !util.CheckFileIsExist("store") {
		err := os.Mkdir("store", os.ModePerm)
		if err != nil {
			return err
		}
		fmt.Println("初始化仓库")
	} else {
		fmt.Println("仓库已存在")
	}
	if !util.CheckFileIsExist("Documentation.json") {
		_, err := os.Create("Documentation.json")
		if err != nil {
			return err
		}
		fmt.Println("初始化仓库清单")
	} else {
		fmt.Println("仓库清单已存在！")
	}
	return nil
}

// 验证文件对应 md5 是否已经存在，并更新文件json
// transitFile 这个功能有严重性能问题，每次处理文件都扫描全部文件的Json，在解析、遍历过程都 O(n) ，如果出现性能瓶颈优先改这里，目前的思路大概是二分Json文件/用树结构组织md码
func transitFile(file *os.File, md5 string) {
	jsonLock.Lock()
	allFiles, err := ioutil.ReadFile("Documentation.json")
	if err != nil {
		fmt.Println("读取文件记录失败", err)
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
	fmt.Println("后台 Json 更新完毕!")
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
