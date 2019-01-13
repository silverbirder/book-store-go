// main.go
package main

import (
	"github.com/Silver-birder/book-store-go/src/controller"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/", controller.IndexGET)
	router.Run(":3000")
}