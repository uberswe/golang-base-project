package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/uberswe/golang-base-project/models"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

// SearchData holds additional data needed to render the search HTML page
type SearchData struct {
	PageData
	Results []models.Website
	Prev    bool
	Next    bool
	PrevURL string
	NextURL string
}

// Search renders the search HTML page and any search results
func (controller Controller) Search(c *gin.Context) {
	page := 1
	resultsPerPage := 5
	pdS := controller.DefaultPageData(c)
	pdS.Title = pdS.Trans("Search")
	pd := SearchData{
		PageData: pdS,
	}
	search := ""
	if c.Request.Method == "POST" && c.Request.RequestURI == "/search" {
		search = c.PostForm("search")
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("/search/1/%s", url.QueryEscape(search)))
		return
	} else {
		search = c.Param("query")
		if i, err := strconv.Atoi(c.Param("page")); err == nil {
			page = i
		}
	}

	var results []models.Website

	log.Println(search)
	searchFilter := fmt.Sprintf("%s%s%s", "%", search, "%")
	search2 := fmt.Sprintf("%s%s", "%", search)
	search4 := fmt.Sprintf("%s%s", search, "%")

	res := controller.db.
		Raw(fmt.Sprintf("SELECT * FROM websites WHERE title LIKE ? OR description LIKE ? ORDER BY CASE WHEN title LIKE ? OR description LIKE ? THEN 1 WHEN title LIKE ? OR description LIKE ? THEN 2 WHEN title LIKE ? OR description LIKE ? THEN 4 ELSE 3 END LIMIT %d OFFSET %d", resultsPerPage, resultsPerPage*(page-1)), searchFilter, searchFilter, search, search, search2, search2, search4, search4).
		Find(&results)

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
	if len(pd.Results) >= resultsPerPage {
		pd.Next = true
		pd.NextURL = fmt.Sprintf("/search/%d/%s", page+1, url.QueryEscape(search))
	}
	if page > 1 {
		pd.Prev = true
		pd.PrevURL = fmt.Sprintf("/search/%d/%s", page-1, url.QueryEscape(search))
	}

	c.HTML(http.StatusOK, "search.html", pd)
}
