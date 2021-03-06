package util

import (
	"fmt"
)

type RegisterPostFrom struct {
	UserName string `from:"user" binding:"required"`
	Password string `from:"password" binding:"required"`
	Email    string
}

func RegisterServer(ch chan *RegisterPostFrom) {
	db, err := GetSQLX()
	if err != nil {
		fmt.Println("注册机链接数据库出错！", err)
		return
	}
	for info := range ch {
		go func(info *RegisterPostFrom) {
			err := NewUserToMySQL(db, info)
			if err != nil {
				fmt.Println("注册失败！", err)
				fmt.Println(info, "没有注册！")
				return
			}
		}(info)
	}
}
