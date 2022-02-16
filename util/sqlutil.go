package util

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

func LoadSQLConfig() *sqlConnectInfo {
	data, err := ioutil.ReadFile("sqlconfig.json")
	sqlConfig := &sqlConnectInfo{}
	if err = json.Unmarshal(data, &sqlConfig); err != nil {
		fmt.Println("json.Unmarshal error:", err)
	}
	fmt.Println("sqlconfig:", sqlConfig)
	return sqlConfig
}

func GetSQLX() (*sqlx.DB, error) {
	sqlConfig := LoadSQLConfig()
	db, err := sqlx.Connect("mysql", sqlConfig.User+":"+sqlConfig.Password+"@tcp("+sqlConfig.Host+")/"+sqlConfig.Database)
	if err != nil {
		return &sqlx.DB{}, errors.New("connect error: " + err.Error())
	}
	return db, nil
}

func NewUserToMySQL(db *sqlx.DB, info *RegisterPostFrom) error {
	_, err := db.Exec("insert into userinfo (uuid, username, password, email,permission) values (uuid(),?,?,?,0)", info.UserName, info.Password, info.Email)
	if err != nil {
		fmt.Println("insert userinfo error:", err)
		return err
	}
	return nil
}
