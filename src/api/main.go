package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Film struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Director    string `json:"director"`
	ReleaseYear int    `json:"release_year"`
	Genre       string `json:"genre"`
}

var films = []Film{
	{
		Id:          "1",
		Title:       "Inception",
		Director:    "Christopher Nolan",
		ReleaseYear: 2010,
		Genre:       "Science Fiction",
	},
	{
		Id:          "2",
		Title:       "The Godfather",
		Director:    "Francis Ford Coppola",
		ReleaseYear: 1972,
		Genre:       "Crime",
	},
	{
		Id:          "3",
		Title:       "Spirited Away",
		Director:    "Hayao Miyazaki",
		ReleaseYear: 2001,
		Genre:       "Animation",
	},
}

func main() {
	router := gin.Default()
	router.GET("/films", getFilms)

	router.Run("localhost:8080")
}

func getFilms(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, films)
}
