package main

import (
	"APIGenerator/rooter"
	"APIGenerator/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
)

func main() {
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	r := rooter.SetupRouter()
	err := r.Run(":" + util.GetPort())
	if err != nil {
		fmt.Println(err)
		return
	}
}
