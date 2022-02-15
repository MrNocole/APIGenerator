package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"strconv"
)

type sqlConnectInfo struct {
	Host     string `json:"host"`
	User     string `json:"username"`
	Password string `json:"password"`
	Database string `json:"dbName"`
}

//type NewUserInfo struct {
//	Username string `db:"username"`
//	Password string `db:"password"`
//	Email    string `db:"email"`
//}
type NewUserInfoInMysql struct {
	//NewUserInfo *NewUserInfo
	Username   string `db:"username"`
	Password   string `db:"password"`
	Email      string `db:"email"`
	UUid       string `db:"uuid"`
	Permission int    `db:"permission"`
}

var newUerInfoTable = `
create table if not exists userinfo(
    uuid char(36) primary key ,
    username text ,
    password text ,
    email text ,
    permission int default 0
)   
`
var initSuperUser = `
insert into userinfo (uuid, username, password, email,permission) values (uuid(),'root','123','null@api.com',-1) ;
`
var dropUserInfoTable = `
drop table if exists userinfo ;
`

func loadSQLConfig() *sqlConnectInfo {
	data, err := ioutil.ReadFile("sqlconfig.json")
	sqlConfig := &sqlConnectInfo{}
	if err = json.Unmarshal(data, &sqlConfig); err != nil {
		fmt.Println("json.Unmarshal error:", err)
	}
	fmt.Println("sqlconfig:", sqlConfig)
	return sqlConfig
}

func Init() (*sqlx.DB, error) {
	sqlConfig := loadSQLConfig()
	db, err := sqlx.Connect("mysql", sqlConfig.User+":"+sqlConfig.Password+"@tcp("+sqlConfig.Host+")/"+sqlConfig.Database)
	if err != nil {
		return &sqlx.DB{}, errors.New("connect error: " + err.Error())
	}
	_, err = db.Exec(newUerInfoTable)
	if err != nil {
		fmt.Println("Table init error! ", err)
		return nil, err
	}
	fmt.Println("init userinfo table...")
	var res NewUserInfoInMysql
	res, err = SelectUserInfoByPermission(db, -1)
	if res == (NewUserInfoInMysql{}) {
		fmt.Println("init super user")
		db.MustExec(initSuperUser)
	} else if err != nil {
		fmt.Println("super user is not found:", err)
		return nil, err
	}
	fmt.Println("super user info:", res)
	return db, nil
}

func SelectUserInfoByPermission(db *sqlx.DB, permission int) (NewUserInfoInMysql, error) {
	userinfo := NewUserInfoInMysql{}
	err := db.Get(&userinfo, "select uuid,username,password,email,permission from userinfo where permission=?", strconv.Itoa(permission))
	if err != nil {
		fmt.Println("select userinfo error:", err)
		return NewUserInfoInMysql{}, err
	}
	return userinfo, nil
}

func SelectPasswordAndUUidByUserName(db *sqlx.DB, username string) (string, string, error) {
	var retInfo struct {
		Password string `db:"password"`
		UUid     string `db:"uuid"`
	}
	err := db.Get(&retInfo, "select password,uuid from userinfo where username=?", username)
	if err != nil {
		return "", "", err
	}
	return retInfo.Password, retInfo.UUid, nil
}
