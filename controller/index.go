package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)


func IndexSave(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{})
}