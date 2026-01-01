package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func StartServer(config *Configuration) {
	router := gin.Default()
	router.GET("/segments", getSegments)
	router.GET("/segments/:id", getSegmentByID)
	router.POST("/segments", postSegments)

	err := router.Run(config.Addr())
	if err != nil {
		log.Fatal(err)
	}
}

func getSegments(c *gin.Context) {
	segments := []SpeechSegment{
		{ID: "1", Title: "ciao"},
		{ID: "2", Title: "ciao2"},
	}

	c.IndentedJSON(http.StatusOK, segments)
}

func getSegmentByID(c *gin.Context) {
	id := c.Param("id")

	segments := []SpeechSegment{
		{ID: "1", Title: "ciao"},
		{ID: "2", Title: "ciao2"},
	}

	for _, a := range segments {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Segment " + id + " not found"})
}

func postSegments(c *gin.Context) {
	var newSeg SpeechSegment

	if err := c.BindJSON(&newSeg); err != nil {
		return
	}

	c.IndentedJSON(http.StatusCreated, newSeg)
}
