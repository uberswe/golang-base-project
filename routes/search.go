package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/uberswe/golang-base-project/models"
	"log"
	"net/http"
)

// SearchData holds additional data needed to render the search HTML page
type SearchData struct {
	PageData
	Results []models.Website
}

// Search renders the search HTML page and any search results
func (controller Controller) Search(c *gin.Context) {
	pdS := controller.DefaultPageData(c)
	pdS.Title = pdS.Trans("Search")
	pd := SearchData{
		PageData: pdS,
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
			Content: pdS.Trans("No results found"),
		})
		log.Println(res.Error)
		c.HTML(http.StatusOK, "search.html", pd)
		return
	}

	pd.Results = results

	c.HTML(http.StatusOK, "search.html", pd)
}
