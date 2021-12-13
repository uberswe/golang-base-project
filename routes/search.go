package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/uberswe/golang-base-project/models"
	"log"
	"net/http"
)

type SearchData struct {
	PageData
	Results []models.Website
}

func (controller Controller) Search(c *gin.Context) {
	pd := SearchData{
		PageData: PageData{
			Title:           "Search",
			IsAuthenticated: isAuthenticated(c),
			CacheParameter:  controller.config.CacheParameter,
		},
	}
	search := c.PostForm("search")

	var results []models.Website

	log.Println(search)
	search = fmt.Sprintf("%s%s%s", "%", search, "%")

	log.Println(search)
	res := controller.db.Where("title LIKE ? OR description LIKE ?", search, search).Find(&results)

	if res.Error != nil || len(results) == 0 {
		pd.Messages = append(pd.Messages, Message{
			Type:    "error",
			Content: "No results found",
		})
		log.Println(res.Error)
		c.HTML(http.StatusOK, "search.html", pd)
		return
	}

	pd.Results = results

	c.HTML(http.StatusOK, "search.html", pd)
}
