package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

//var albums = []album{
//	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
//	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
//	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
//}

func getAlbums(c *gin.Context) {
	var albums []album
	rows, err := db.Query("SELECT * FROM album")
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
		return
	}

	for rows.Next() {
		var alb album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message: getAlbums:": err})
			return
		}
		albums = append(albums, alb)
	}

	c.IndentedJSON(http.StatusOK, albums)
}

//func postAlbums(c *gin.Context) {
//	var newAlbum album
//
//	if err := c.BindJSON(&newAlbum); err != nil {
//		return
//	}
//	albums = append(albums, newAlbum)
//	c.IndentedJSON(http.StatusCreated, newAlbum)
//}
//
//func getAlbumByID(c *gin.Context) {
//	id := c.Param("id")
//
//	for _, a := range albums {
//		if a.ID == id {
//			c.IndentedJSON(http.StatusOK, a)
//			return
//		}
//	}
//	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
//}

var db *sql.DB

func main() {
	cfg := mysql.Config{
		User:   "root",
		Passwd: "docker",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "album",
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("connected!")

	router := gin.Default()
	router.GET("/albums", getAlbums)
	//router.GET("/albums/:id", getAlbumByID)
	//router.POST("/albums", postAlbums)
	err = router.Run("localhost:8080")
	if err != nil {
		log.Fatalln(err.Error())
	}
}
