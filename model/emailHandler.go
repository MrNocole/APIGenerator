package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/smtp"
	"strings"
	"sync"
)

type userInfo struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
}

var userinfo userInfo
var config_file_locker sync.Mutex

func MailConfig() error {
	config_file_locker.Lock()
	data, err := ioutil.ReadFile("mailconfig.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &userinfo)
	if err != nil {
		return err
	}
	return nil
}
func SendToMail(user, password, host, to, subject, body, mailtype string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var content_type string
	if mailtype == "html" {
		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	msg := []byte("To: " + to + "\r\nFrom: " + user + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, send_to, msg)
	return err
}

func SendAnEmailDefault(to, subject, body string) error {
	if err := MailConfig(); err != nil {
		return err
	}
	fmt.Println(userinfo)
	return SendToMail(userinfo.User, userinfo.Password, userinfo.Host, to, subject, body, "html")
}

func SendAnEmail(to, subject, body, mailtype string) error {
	if err := MailConfig(); err != nil {
		return err
	}
	return SendToMail(userinfo.User, userinfo.Password, userinfo.Host, to, subject, body, mailtype)
}

func SendVerify(to, code string) error {
	body := "verify code :" + code
	fmt.Println(body + "  " + to)
	return SendAnEmailDefault(to, "verify code", body)
	//user := "1926276913@qq.com"
	//password := "hiuppgxofkbpejjb"
	//host := "smtp.qq.com:25"
	//to := "1926276913@qq.com"
	//body := "your verify code is :" + fmt.Sprintf("%d", rand.Intn(8999)+1000)
	//subject := "test"
	//	body := `
	//	<html>
	//		<body>
	//			<h3>
	//				your verify code is :
	//			</h3>
	//			<h3>
	//
	//			</h3>
	//		</body>
	//	</html>
	//`
	//err := SendToMail(user, password, host, to, subject, body, "html")
	//if err != nil {
	//	fmt.Println("Error!")
	//	return err
	//} else {
	//	fmt.Println("Send success")
	//}
}
