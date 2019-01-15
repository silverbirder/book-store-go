// main.go
package main

import (
	"github.com/Silver-birder/book-store-go/controller"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Static("/css/", "./public/css")
	router.Static("/js/", "./public/js/")
	router.GET("/api/book/0.1/add", controller.AddBook)
	router.Run(":3000")
}