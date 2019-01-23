package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/algolia/algoliasearch-client-go/algoliasearch"
)
var ALGOLIA_APP_KEY = os.Getenv("ALGOLIA_APP_KEY")
var ALGOLIA_ADMIN_KEY = os.Getenv("ALGOLIA_ADMIN_KEY")
var ALGOLIA_INDEX_KEY = os.Getenv("ALGOLIA_INDEX_KEY")
var OPEN_BD_URL = "https://api.openbd.jp/v1/get"
var NO_IMAGE_URL = "https://res.cloudinary.com/silverbirder/image/upload/v1548200876/no-image.png"

// https://api.openbd.jp/v1/get?isbn=9784621086025
// https://mholt.github.io/json-to-go/
type OpenBd struct {
	Onix struct {
		RecordReference   string `json:"RecordReference"`
		NotificationType  string `json:"NotificationType"`
		ProductIdentifier struct {
			ProductIDType string `json:"ProductIDType"`
			IDValue       string `json:"IDValue"`
		} `json:"ProductIdentifier"`
		DescriptiveDetail struct {
			ProductComposition string `json:"ProductComposition"`
			ProductForm        string `json:"ProductForm"`
			Measure            []struct {
				MeasureType     string `json:"MeasureType"`
				Measurement     string `json:"Measurement"`
				MeasureUnitCode string `json:"MeasureUnitCode"`
			} `json:"Measure"`
			TitleDetail struct {
				TitleType    string `json:"TitleType"`
				TitleElement struct {
					TitleElementLevel string `json:"TitleElementLevel"`
					TitleText         struct {
						Collationkey string `json:"collationkey"`
						Content      string `json:"content"`
					} `json:"TitleText"`
				} `json:"TitleElement"`
			} `json:"TitleDetail"`
			Contributor []struct {
				SequenceNumber  string   `json:"SequenceNumber"`
				ContributorRole []string `json:"ContributorRole"`
				PersonName      struct {
					Content string `json:"content"`
				} `json:"PersonName"`
			} `json:"Contributor"`
			Language []struct {
				LanguageRole string `json:"LanguageRole"`
				LanguageCode string `json:"LanguageCode"`
				CountryCode  string `json:"CountryCode"`
			} `json:"Language"`
			Extent []struct {
				ExtentType  string `json:"ExtentType"`
				ExtentValue string `json:"ExtentValue"`
				ExtentUnit  string `json:"ExtentUnit"`
			} `json:"Extent"`
		} `json:"DescriptiveDetail"`
		CollateralDetail struct {
			TextContent []struct {
				TextType        string `json:"TextType"`
				ContentAudience string `json:"ContentAudience"`
				Text            string `json:"Text"`
			} `json:"TextContent"`
			SupportingResource []struct {
				ResourceContentType string `json:"ResourceContentType"`
				ContentAudience     string `json:"ContentAudience"`
				ResourceMode        string `json:"ResourceMode"`
				ResourceVersion     []struct {
					ResourceForm           string `json:"ResourceForm"`
					ResourceVersionFeature []struct {
						ResourceVersionFeatureType string `json:"ResourceVersionFeatureType"`
						FeatureValue               string `json:"FeatureValue"`
					} `json:"ResourceVersionFeature"`
					ResourceLink string `json:"ResourceLink"`
				} `json:"ResourceVersion"`
			} `json:"SupportingResource"`
		} `json:"CollateralDetail"`
		PublishingDetail struct {
			Imprint struct {
				ImprintIdentifier []struct {
					ImprintIDType string `json:"ImprintIDType"`
					IDValue       string `json:"IDValue"`
				} `json:"ImprintIdentifier"`
				ImprintName string `json:"ImprintName"`
			} `json:"Imprint"`
			PublishingDate []struct {
				PublishingDateRole string `json:"PublishingDateRole"`
				Date               string `json:"Date"`
			} `json:"PublishingDate"`
		} `json:"PublishingDetail"`
		ProductSupply struct {
			SupplyDetail struct {
				ReturnsConditions struct {
					ReturnsCodeType string `json:"ReturnsCodeType"`
					ReturnsCode     string `json:"ReturnsCode"`
				} `json:"ReturnsConditions"`
				ProductAvailability string `json:"ProductAvailability"`
			} `json:"SupplyDetail"`
		} `json:"ProductSupply"`
	} `json:"onix"`
	Hanmoto struct {
		Datecreated  string `json:"datecreated"`
		Dateshuppan  string `json:"dateshuppan"`
		Datemodified string `json:"datemodified"`
	} `json:"hanmoto"`
	Summary struct {
		Isbn      string `json:"isbn"`
		Title     string `json:"title"`
		Volume    string `json:"volume"`
		Series    string `json:"series"`
		Publisher string `json:"publisher"`
		Pubdate   string `json:"pubdate"`
		Cover     string `json:"cover"`
		Author    string `json:"author"`
	} `json:"summary"`
}

var client = algoliasearch.NewClient(ALGOLIA_APP_KEY, ALGOLIA_ADMIN_KEY)
var index = client.InitIndex(ALGOLIA_INDEX_KEY)

func AddBook(c *gin.Context) {
	isbn := c.Query("isbn")
	settings := algoliasearch.Map{
		"attributesToRetrieve": []string{"isbn"},
	}
	res, err := index.Search(isbn, settings)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "fail"})
		return
	}
	for i := 0; i < len(res.Hits); i++  {
		if res.Hits[i]["isbn"] == isbn {
			c.JSON(http.StatusOK, gin.H{"status": "duplication"})
			return
		}
	}
	openBd, err := fetch(isbn)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "fail"})
		return
	}
	if openBd[0].Summary.Isbn == "" {
		c.JSON(http.StatusOK, gin.H{"status": "not exist"})
		return
	}
	textContent := openBd[0].Onix.CollateralDetail.TextContent
	summary := openBd[0].Summary
	cover := summary.Cover
	object := algoliasearch.Object{
		"isbn":  summary.Isbn,
		"title":  summary.Title,
		"volume":  summary.Volume,
		"series":  summary.Series,
		"publisher":  summary.Publisher,
		"pubdate":  summary.Pubdate,
		"cover":  cover,
		"author":  summary.Author,
		"textContent": textContent,
	}
	_, err = index.AddObject(object)
	if err != nil {
		log.Fatalf("Error!: %v", err)
		c.JSON(http.StatusOK, gin.H{"status": "fail"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func fetch(isbn string) ([]OpenBd, error) {
	res, err := http.Get(OPEN_BD_URL + "?isbn=" + isbn)
	if err != nil {
		return nil, err
	} else if res.StatusCode != 200 {
		return nil, fmt.Errorf("Unable to get this url : http status %d", res.StatusCode)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var OpenBd []OpenBd
	if err := json.Unmarshal(body, &OpenBd); err != nil {
	 	return nil, err
	}
	return OpenBd, nil
}