package main

import (
	"APIGenerator/rooter"
	"github.com/gin-gonic/gin"
	"io"
	"os"
)

func main() {

	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	r := rooter.SetupRouter()
	r.Run(":8000")
}
