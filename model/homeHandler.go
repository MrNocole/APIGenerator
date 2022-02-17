package model

import (
	"APIGenerator/util"
	"fmt"
	"github.com/gin-gonic/gin"
)

type Item struct {
	Name string
	URL  string
	Md5  string
}

func HomeHandler(c *gin.Context) {
	items := getItemList(c)
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

func getItemList(c *gin.Context) []Item {
	var items []Item
	uuid, err := c.Cookie("uuid")
	if err != nil {
		fmt.Println("用户列表获取失败", err)
		errItem := Item{
			Name: "用户列表获取失败",
			URL:  "/404",
		}
		items = append(items, errItem)
	} else {
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
	}
	fmt.Println(items)
	return items
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
	default:
		c.Redirect(302, "/download/"+uuid+"/"+fileName)
	}
}
