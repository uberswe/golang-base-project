package baseproject

import (
	"fmt"
	"github.com/uberswe/golang-base-project/config"
	"github.com/uberswe/golang-base-project/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

func connectToDatabase(c config.Config) (db *gorm.DB, err error) {
	return connectLoop(c, 0)
}

func connectLoop(c config.Config, count int) (db *gorm.DB, err error) {
	db, err = attemptConnection(c)
	if err != nil {
		if count > 300 {
			return db, fmt.Errorf("could not connect to database after 300 seconds")
		}
		time.Sleep(1 * time.Second)
		return connectLoop(c, count+1)
	}
	return db, err
}

func attemptConnection(c config.Config) (db *gorm.DB, err error) {
	if c.Database == "sqlite" {
		// In-memory sqlite if no database name is specified
		dsn := "file::memory:?cache=shared"
		if c.DatabaseName != "" {
			dsn = fmt.Sprintf("%s.db", c.DatabaseName)
		}
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	} else if c.Database == "mysql" {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.DatabaseUsername, c.DatabasePassword, c.DatabaseHost, c.DatabasePort, c.DatabaseName)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	} else if c.Database == "postgres" {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC", c.DatabaseHost, c.DatabaseUsername, c.DatabasePassword, c.DatabaseName, c.DatabasePort)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	} else {
		return db, fmt.Errorf("no database specified: %s", c.Database)
	}
	return db, err
}

func migrateDatabase(db *gorm.DB) error {
	err := db.AutoMigrate(&models.User{}, &models.Token{}, &models.Session{}, &models.Website{})
	seed(db)
	return err
}

func seed(db *gorm.DB) {
	// We seed some websites for our search results
	websites := []models.Website{
		{
			Title:       "A Tour of Go",
			Description: "A Tour of Go has several interactive examples of how Go which you can learn from. There is a menu available if you would like to skip to different sections.",
			URL:         "https://go.dev/tour/welcome/1",
		},
		{
			Title:       "Go by Example",
			Description: "As described on the website: Go by Example is a hands-on introduction to Go using annotated example programs. I have used this site many times as a reference when I need to look something up.",
			URL:         "https://gobyexample.com/",
		},
		{
			Title:       "Go.dev",
			Description: "Learn how to install Go on your machine and read the documentation on the Go website.",
			URL:         "https://go.dev/learn/",
		},
	}

	for _, w := range websites {
		res := db.Where(&w).First(&w)
		// If no record exists we insert
		if res.Error != nil && res.Error == gorm.ErrRecordNotFound {
			db.Save(&w)
		}
	}
}
