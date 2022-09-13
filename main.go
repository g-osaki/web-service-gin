package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

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

func postAlbums(c *gin.Context) {
	var alb album
	if err := c.BindJSON(&alb); err != nil {
		fmt.Errorf(err.Error())
		return
	}

	result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?,?,?)", alb.Title, alb.Artist, alb.Price)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	id, err := result.LastInsertId()
	if err != nil {
		fmt.Errorf(err.Error())
	}
	alb.ID = strconv.FormatInt(id, 10)

	c.IndentedJSON(http.StatusCreated, alb)
}

func getAlbumByID(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	var alb album

	row := db.QueryRow("SELECT * FROM album WHERE id = ?", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		fmt.Errorf(err.Error())
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, alb)
}

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
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)
	err = router.Run("localhost:8080")
	if err != nil {
		log.Fatalln(err.Error())
	}
}
