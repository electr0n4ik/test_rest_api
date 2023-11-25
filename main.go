package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
)

var db *gorm.DB

type Album struct {
	gorm.Model
	Title  string
	Artist string
	Price  float64
}

func main() {
	defer func() {
		if db != nil {
			sqlDB, err := db.DB()
			if err != nil {
				log.Fatal(err)
			}
			sqlDB.Close()
		}
	}()

	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)

	router.Run("localhost:8080")
}

func getAlbums(c *gin.Context) {
	var result []Album
	if err := db.Find(&result).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	var a Album
	if err := db.Where("id = ?", id).First(&a).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		}
		return
	}

	c.IndentedJSON(http.StatusOK, a)
}

// postAlbums добавляет альбом из JSON, полученного в теле запроса.
func postAlbums(c *gin.Context) {
	var newAlbum Album

	// Вызовите BindJSON, чтобы привязать полученный JSON к
	// newAlbum.
	if err := c.BindJSON(&newAlbum); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	result := db.WithContext(context.Background()).Exec("INSERT INTO albums (title, artist, price) VALUES (?, ?, ?)", newAlbum.Title, newAlbum.Artist, newAlbum.Price)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func init() {
	var err error
	db, err = connectDB()
	if err != nil {
		log.Fatal(err)
	}
}

func connectDB() (*gorm.DB, error) {
	connString := "postgresql://postgres:12345@localhost:5432/postgres?sslmode=disable"
	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
