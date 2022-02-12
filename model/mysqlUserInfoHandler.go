package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
)

type sqlConnectInfo struct {
	Host     string `json:"host"`
	User     string `json:"username"`
	Password string `json:"password"`
	Database string `json:"dbName"`
}

type NewUserInfo struct {
	Username string `db:"username"`
	Password string `db:"password"`
	Email    string `db:"email"`
}
type NewUserInfoInMysql struct {
	NewUserInfo *NewUserInfo
	UUid        string `db:"uuid"`
	Permission  int    `db:"permission"`
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
	db, err := sqlx.Open("mysql", sqlConfig.User+":"+sqlConfig.Password+"@tcp("+sqlConfig.Host+")/"+sqlConfig.Database)
	if err != nil {
		return &sqlx.DB{}, errors.New("connect error: " + err.Error())
	}
	db.MustExec(newUerInfoTable)
	fmt.Println("init userinfo table...")
	var res NewUserInfoInMysql
	if res, err = SelectUserInfoByPermission(db, -1); err != nil {
		fmt.Println("select super user error:", err)
		db.MustExec(dropUserInfoTable)
		return &sqlx.DB{}, err
	}
	if res == (NewUserInfoInMysql{}) {
		fmt.Println("init super user")
		db.MustExec(initSuperUser)
	}
	if res, err = SelectUserInfoByPermission(db, -1); err != nil {
		fmt.Println("select super user error:", err)
		db.MustExec(dropUserInfoTable)
		return db, err
	}
	fmt.Println("super user info:", res)
	fmt.Println("init mysql success")
	return db, nil
}

func SelectUserInfoByPermission(db *sqlx.DB, uuid int) (NewUserInfoInMysql, error) {
	var userinfo NewUserInfoInMysql
	err := db.Get(&userinfo, "select * from userinfo where permission=?", uuid)
	if err != nil {
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
