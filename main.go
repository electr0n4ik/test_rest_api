package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

var db *pgxpool.Pool

type Album struct {
	ID     uint `gorm:"primaryKey"`
	Title  string
	Artist string
	Price  float64
}

func main() {
	defer db.Close()

	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)

	router.Run("localhost:8080")

}

func getAlbums(c *gin.Context) {
	var result []Album
	if err := db.Model(&Album{}).Find(&result).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	var a Album
	if err := db.First(&a, id).Error; err != nil {
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
	var newAlbum album

	// Вызовите BindJSON, чтобы привязать полученный JSON к
	// newAlbum.
	if err := c.BindJSON(&newAlbum); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Вставьте новый альбом в базу данных.
	_, err := db.Exec(context.Background(), "INSERT INTO albums (title, artist, price) VALUES ($1, $2, $3)", newAlbum.Title, newAlbum.Artist, newAlbum.Price)
	if err != nil {
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

func connectDB() (*pgxpool.Pool, error) {
	connString := "postgresql://postgres:12345@localhost:5432/postgres?sslmode=disable"
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	db, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	return db, nil
}
