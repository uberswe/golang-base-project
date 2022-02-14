package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/uberswe/golang-base-project/models"
	"gorm.io/gorm"
	"log"
	"math"
	"net/http"
	"sort"
	"time"
)

type AdminData struct {
	PageData
	Chart Chart
}

type Chart struct {
	Title   string
	XLabel  string
	YLabel  string
	XValues []ChartValue
	YValues []ChartValue
	Points  string
}

type ChartValue struct {
	Value string
	X     int
	Y     int
}

type MonthlySignUps struct {
	Year  int
	Month int
	Count int
}

// Admin renders the admin dashboard
func (controller Controller) Admin(c *gin.Context) {
	pd := controller.DefaultPageData(c)
	pd.Title = pd.Trans("Admin")

	ad := AdminData{
		PageData: pd,
		Chart: Chart{
			Title:  "Monthly User Sign Ups",
			XLabel: "Month",
			YLabel: "Users",
		},
	}

	var msu []MonthlySignUps

	res := controller.db.Model(&models.User{}).
		Select("YEAR(created_at) as year, MONTH(created_at) AS month, COUNT(*) AS count").
		Where("created_at >= CURDATE() - INTERVAL 1 YEAR").
		Group("YEAR(created_at), MONTH(created_at)").
		Find(&msu)

	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		pd.Messages = append(pd.Messages, Message{
			Type:    "error",
			Content: "Something went wrong while fetching user data",
		})
		log.Println(res.Error)
		c.HTML(http.StatusInternalServerError, "admin.html", ad)
		return
	}

	// TODO make this into a graph package or find a go package to replace this manual code
	// Add 0 values
	y, m, _ := time.Now().AddDate(-1, 0, 0).Date()
	ny, nm, _ := time.Now().Date()
	for y < ny || m < nm {
		found := false
		for _, ms := range msu {
			if ms.Year == y && ms.Month == int(m) {
				found = true
			}
		}
		if !found {
			msu = append(msu, MonthlySignUps{
				Year:  y,
				Month: int(m),
				Count: 0,
			})
		}
		m++
		if m > 12 {
			m = 1
			y++
		}
	}

	// Sort our values so that the graph shows left to right going from earliest to latest date
	sort.Slice(msu, func(i, j int) bool {
		if msu[i].Year < msu[j].Year {
			return true
		} else if msu[i].Year == msu[j].Year && msu[i].Month < msu[j].Month {
			return true
		}
		return false
	})

	rightAdjust := 50
	base := 350
	lowest := math.MaxInt32
	highest := 0
	count := len(msu)
	for _, m := range msu {
		if lowest > m.Count {
			lowest = m.Count
		}
		if highest < m.Count {
			highest = m.Count
		}
	}
	diff := highest - lowest
	if diff < 5 {
		diff = 5
	}
	// We multiply counts by this ratio after subtracting our lowest
	ratio := base / diff

	// y line is 90 x to 705 x at 370 y

	maxX := 705
	minX := 90

	seg := (maxX - minX) / count

	var yValues []ChartValue
	for i, m := range msu {
		yValues = append(yValues, ChartValue{
			Value: fmt.Sprintf("%d-%d", m.Year, m.Month),
			X:     rightAdjust + 20 + minX + (i * seg),
			Y:     base + 40,
		})
	}

	// x line is 5 y to 371 y on 90 x

	var xValues []ChartValue
	seg2 := diff / 5

	for i := 0; i <= 5; i++ {
		xValues = append(xValues, ChartValue{
			Value: fmt.Sprintf("%d", lowest+(seg2*i)+1),
			X:     minX - 20,
			Y:     (base - 10) - (i * seg2 * ratio),
		})
	}

	points := ""
	for i, m := range msu {
		if i > 0 {
			points += " "
		}
		points += fmt.Sprintf("%d,%d", rightAdjust+minX+(i*seg), (base+20)-(m.Count*ratio))
	}

	ad.Chart.XValues = xValues
	ad.Chart.YValues = yValues
	ad.Chart.Points = points

	// The chart is inverted so we need to subtract the base value to our calculated values

	c.HTML(http.StatusOK, "admin.html", ad)
}
