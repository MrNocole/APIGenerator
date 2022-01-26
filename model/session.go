package model

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func Session(keyPairs string, sessionAge int, sessionSecret string) gin.HandlerFunc {
	store := SessionConfig(sessionAge, sessionSecret)
	return sessions.Sessions(keyPairs, store)
}

func SessionDefault(keyPairs string) gin.HandlerFunc {
	store := SessionConfig(60, "zhangtao")
	return sessions.Sessions(keyPairs, store)
}

func SessionConfig(sessionAge int, sessionSecret string) sessions.Store {
	var store sessions.Store
	//fmt.Println("New session")
	store = cookie.NewStore([]byte(sessionSecret))
	store.Options(sessions.Options{MaxAge: sessionAge, Path: "/"})
	return store
}
