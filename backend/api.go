package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func StartServer(config *Configuration, db *Db) {
	router := gin.Default()

	router.GET("/api/exercises", getExercises(db))
	router.POST("/api/exercises", postExercises(db))

	router.GET("/api/segments/:exercise_id", getSegments(db))
	router.POST("/api/segments/:exercise_id", postSegments(db))

	err := router.Run(config.Addr())
	if err != nil {
		log.Fatal(err)
	}
}

/*************************************
 * Exercise
 */

func getExercises(db *Db) gin.HandlerFunc {
	return func(c *gin.Context) {
		exercises, err := db.GetExercises()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err})
			return
		}

		c.JSON(http.StatusOK, gin.H{"exercises": exercises})
	}
}

func postExercises(db *Db) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newExercise SpeechExercise

		err := c.BindJSON(&newExercise)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err})
			return
		}

		err = db.InsertExercise(newExercise)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Exercise added successfully", "created": newExercise})
	}
}

/*************************************
 * Segments
 */

func getSegments(db *Db) gin.HandlerFunc {
	return func(c *gin.Context) {
		exerciseId := c.Param("exercise_id")

		segments, err := db.GetSegments(exerciseId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err})
			return
		}

		c.JSON(http.StatusOK, gin.H{"segments": segments})
	}
}

func postSegments(db *Db) gin.HandlerFunc {
	return func(c *gin.Context) {
		exerciseId := c.Param("exercise_id")

		var newSegment SpeechSegment
		err := c.BindJSON(&newSegment)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err})
			return
		}

		err = db.InsertSegment(exerciseId, newSegment)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Segment added successfully", "created": newSegment})
	}
}
