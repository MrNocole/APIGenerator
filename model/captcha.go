package model

import (
	"bytes"
	"fmt"
	"github.com/dchest/captcha"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func Captcha(c *gin.Context, length ...int) {
	fmt.Println("set captcha")
	l := captcha.DefaultLen
	w, h := 107, 36
	if len(length) == 1 {
		l = length[0]
	}
	if len(length) == 2 {
		w = length[1]
	}
	if len(length) == 3 {
		h = length[2]
	}
	captchaId := captcha.NewLen(l)
	session := sessions.Default(c)
	session.Set("captcha", captchaId)
	if err := session.Save(); err != nil {
		fmt.Println(err)
	}
	_ = Serve(c.Writer, c.Request, captchaId, "zh", false, w, h)
}

func Serve(w http.ResponseWriter, r *http.Request, id, lang string, download bool, width, height int) error {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	var content bytes.Buffer
	w.Header().Set("Content-Type", "image/png")
	_ = captcha.WriteImage(&content, id, width, height)

	if download {
		w.Header().Set("Content-Type", "application/octet-stream")
	}
	http.ServeContent(w, r, id+".png", time.Time{}, bytes.NewReader(content.Bytes()))
	return nil
}

func checkCaptcha(code string, c *gin.Context) bool {
	session := sessions.Default(c)
	if captchaId := session.Get("captcha"); captchaId != nil {
		fmt.Println("captchaId got")
		session.Delete("captcha")
		_ = session.Save()
		if captcha.VerifyString(captchaId.(string), code) {
			return true
		} else {
			return false
		}
	} else {
		fmt.Println("captcha session not found")
		return false
	}

}
