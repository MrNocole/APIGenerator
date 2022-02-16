package util

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ErrorHtml(c *gin.Context, errorCode, msg string) {
	c.HTML(http.StatusBadGateway, "error.html", gin.H{"errorCode": errorCode, "info": msg})
}
