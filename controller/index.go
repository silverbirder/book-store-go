package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/algolia/algoliasearch-client-go/algoliasearch"
)
var ALGOLIA_APP_KEY = os.Getenv("ALGOLIA_APP_KEY")
var ALGOLIA_ADMIN_KEY = os.Getenv("ALGOLIA_ADMIN_KEY")
var ALGOLIA_INDEX_KEY = os.Getenv("ALGOLIA_INDEX_KEY")
var client = algoliasearch.NewClient(ALGOLIA_APP_KEY, ALGOLIA_ADMIN_KEY)
var index = client.InitIndex(ALGOLIA_INDEX_KEY)

func AddBook(c *gin.Context) {
	isbn := c.Query("isbn")
	content, _ := ioutil.ReadFile("./controller/contacts.json")
	var objects []algoliasearch.Object
	if err := json.Unmarshal(content, &objects); err != nil {
		fmt.Println(err)
		return
	}
	_, _ = index.AddObjects(objects)
	c.JSON(http.StatusOK, gin.H{"status": isbn})
}
