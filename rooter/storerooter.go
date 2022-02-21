package rooter

import (
	"APIGenerator/model"
	"github.com/gin-gonic/gin"
)

func initStoreRooter(r *gin.Engine) {
	r.GET("/download/:uuid/:filename", model.DownloadByAPI)
	r.GET("/home", model.UserCookieCheck, func(c *gin.Context) {
		model.HomeHandler(c, redisPool)
	})
	r.GET("/delete/:uuid/:filename", func(c *gin.Context) {
		model.FileDeleteHandler(c, redisPool.Get())
	})
	r.POST("/upload", func(c *gin.Context) {
		model.FileUploadHandler(c, redisPool.Get())
	})
	r.GET("/check/:uuid/:name", model.CheckHandler)
	r.GET("/json/:uuid/:name", model.GetJson)
	r.GET("/pic/:uuid/:name", model.GetPic)
}
