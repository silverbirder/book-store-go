// main.go
package main

import (
	"github.com/Silver-birder/book-store-go/controller"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
	}))
	router.GET("/api/book/0.1/add", controller.AddBook)
	router.Run(":3000")
}