package model

import (
	"APIGenerator/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
)

type Item struct {
	Name string
	URL  string
	Md5  string
}

func HomeHandler(c *gin.Context, pool *redis.Pool) {
	items := getItemList(c, pool)
	name, err := c.Cookie("userName")
	if err != nil {
		name = "UnKnown"
	}
	uuid, err := c.Cookie("uuid")
	if err != nil {
		fmt.Println("UUID is not found!")
		c.Redirect(302, "/login")
	}
	c.HTML(200, "home.html", gin.H{
		"items":    items,
		"userName": name,
		"uuid":     uuid,
		"link":     util.GetUrl(),
	})
}

func getItemList(c *gin.Context, pool *redis.Pool) []Item {
	var items []Item
	uuid, err := c.Cookie("uuid")
	if err != nil {
		fmt.Println("用户列表获取失败", err)
		errItem := Item{
			Name: "用户列表获取失败",
			URL:  "/404",
		}
		items = append(items, errItem)
		return items
	} else {
		// 先查redis有没有，没有就到硬盘找
		fileNames, err := util.RedisGetHKeys(pool.Get(), uuid)
		md5s, err := util.RedisGetHVals(pool.Get(), uuid)
		fmt.Println(fileNames, " ", md5s)
		// 以下两个if是redis中没有的情况
		if err != nil {
			fmt.Println("Redis 读取Items失败", err)
			items = getItemListFromDisk(uuid)
			go updateRedisItemsData(uuid, items, pool.Get())
			fmt.Println("redis 更新中")
		} else if fileNames == nil || len(fileNames) == 0 {
			fmt.Println("Redis 未命中")
			items = getItemListFromDisk(uuid)
			go updateRedisItemsData(uuid, items, pool.Get())
			fmt.Println("redis 更新中")
		} else {
			fmt.Println("redis 命中")
			for i := 0; i < Min(len(fileNames), len(md5s)); i++ {
				tmpItem := Item{
					Name: fileNames[i],
					URL:  "/download/" + uuid + "/" + fileNames[i],
					Md5:  md5s[i],
				}
				items = append(items, tmpItem)
			}
		}
		return items
	}
}
func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
func getItemListFromDisk(uuid string) []Item {
	var items []Item
	ownerInfo, err := GetOwnerInfo(uuid)
	defer GetOwnerLock(uuid).Unlock()
	if err != nil {
		errItem := Item{
			Name: "用户列表获取失败",
			URL:  "/404",
		}
		items = append(items, errItem)
	} else {
		for i, v := range ownerInfo.FileName {
			item := Item{
				Name: v,
				URL:  "/download/" + uuid + "/" + v,
				Md5:  ownerInfo.MD5[i],
			}
			items = append(items, item)
		}
	}
	return items
}

func updateRedisItemsData(uuid string, items []Item, Conn redis.Conn) {
	var fileNames []string
	var md5s []string
	for _, item := range items {
		fileNames = append(fileNames, item.Name)
		md5s = append(md5s, item.Md5)
	}
	func(fileNames []string, md5s []string) {
		for i, fileName := range fileNames {
			util.RedisInsertH(Conn, uuid, fileName, md5s[i])
		}
	}(fileNames, md5s)

}

func UserCookieCheck(c *gin.Context) {
	fmt.Println("MiddleWare begin...")
	fmt.Println(c.Cookie("userName"))

	uuid, err := c.Cookie("uuid")
	if err != nil {
		fmt.Println("Cookie not found")
		c.Abort()
	}
	res, err := checkUser(uuid)
	if err != nil {
		fmt.Println("checkUserInfo error")
		c.Abort()
	}
	if res {
		fmt.Println("checkUserInfo success")
		c.Next()
	} else {
		fmt.Println("checkUserInfo fail")
		c.Abort()
	}
}

func checkUser(uuid string) (bool, error) {
	if uuid != "" {
		return true, nil
	}
	return false, nil
}

func CheckHandler(c *gin.Context) {
	uuid := c.Param("uuid")
	fileName := c.Param("name")
	suffix := util.GetSuffix(fileName)
	fmt.Println(suffix)
	switch suffix {
	case "json":
		c.Redirect(302, "/json/"+uuid+"/"+fileName)
	case "jpg":
		c.Redirect(302, "/pic/"+uuid+"/"+fileName)
	case "png":
		c.Redirect(302, "/pic/"+uuid+"/"+fileName)
	default:
		c.Redirect(302, "/download/"+uuid+"/"+fileName)
	}
}
