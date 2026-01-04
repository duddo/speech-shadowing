package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func StartServer(config *Configuration, db *Db) {
	router := gin.Default()

	router.GET("/api/exercises", getExercises(db))
	router.POST("/api/exercises", postExercise(db))
	router.PUT("/api/exercises/:exercise_id", putExercise(db))
	router.DELETE("/api/exercises/:exercise_id", deleteExercise(db))

	router.GET("/api/exercises/:exercise_id/segments", getSegments(db))
	router.POST("/api/exercises/:exercise_id/segments", postSegment(db))
	router.PUT("/api/exercises/:exercise_id/segments/:segment_id", putSegment(db))
	router.DELETE("/api/exercises/:exercise_id/segments/:segment_id", deleteSegment(db))

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
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"exercises": exercises})
	}
}

func postExercise(db *Db) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newExercise SpeechExercise

		err := c.BindJSON(&newExercise)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		insertedID, err := db.InsertExercise(newExercise)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		newExercise.ID = *insertedID

		c.JSON(http.StatusCreated, gin.H{"message": "Exercise added successfully", "created": newExercise})
	}
}

func putExercise(db *Db) gin.HandlerFunc {
	return func(c *gin.Context) {
		exerciseId := c.Param("exercise_id")

		exerciseIdInt, err := strconv.ParseInt(exerciseId, 10, 64)
		if err != nil || exerciseIdInt <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid exercise ID"})
			return
		}

		var newExercise SpeechExercise

		err = c.BindJSON(&newExercise)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		newExercise.ID = exerciseIdInt

		err = db.UpdateExercise(newExercise)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Exercise updated successfully", "affected": newExercise})
	}
}

func deleteExercise(db *Db) gin.HandlerFunc {
	return func(c *gin.Context) {
		exerciseId := c.Param("exercise_id")

		exerciseIdInt, err := strconv.ParseInt(exerciseId, 10, 64)
		if err != nil || exerciseIdInt <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid exercise ID"})
			return
		}

		err = db.DeleteExercise(exerciseIdInt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Exercise deleted successfully", "affected": exerciseId})
	}
}

/*************************************
 * Segments
 */

func getSegments(db *Db) gin.HandlerFunc {
	return func(c *gin.Context) {
		exerciseId := c.Param("exercise_id")

		exerciseIdInt, err := strconv.ParseInt(exerciseId, 10, 64)
		if err != nil || exerciseIdInt <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid exercise ID"})
			return
		}

		segments, err := db.GetSegments(exerciseIdInt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"segments": segments})
	}
}

func postSegment(db *Db) gin.HandlerFunc {
	return func(c *gin.Context) {
		exerciseId := c.Param("exercise_id")
		exerciseIdInt, err := strconv.ParseInt(exerciseId, 10, 64)
		if err != nil || exerciseIdInt <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid exercise ID"})
			return
		}

		var newSegment SpeechSegment
		err = c.BindJSON(&newSegment)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		insertedID, err := db.InsertSegment(exerciseIdInt, newSegment)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		newSegment.ID = *insertedID

		c.JSON(http.StatusCreated, gin.H{"message": "Segment added successfully", "created": newSegment})
	}
}

func putSegment(db *Db) gin.HandlerFunc {
	return func(c *gin.Context) {
		exerciseId := c.Param("exercise_id")
		exerciseIdInt, err := strconv.ParseInt(exerciseId, 10, 64)
		if err != nil || exerciseIdInt <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid exercise ID"})
			return
		}

		segmentId := c.Param("segment_id")
		segmentIdInt, err := strconv.ParseInt(segmentId, 10, 64)
		if err != nil || segmentIdInt <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid segment ID"})
			return
		}

		var newSegment SpeechSegment
		err = c.BindJSON(&newSegment)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		newSegment.ID = segmentIdInt
		newSegment.ExerciseID = exerciseIdInt
		err = db.UpdateSegment(newSegment)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Segment updated successfully", "updated": newSegment})
	}
}

func deleteSegment(db *Db) gin.HandlerFunc {
	return func(c *gin.Context) {
		exerciseId := c.Param("exercise_id")
		exerciseIdInt, err := strconv.ParseInt(exerciseId, 10, 64)
		if err != nil || exerciseIdInt <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid exercise ID"})
			return
		}

		segmentId := c.Param("segment_id")
		segmentIdInt, err := strconv.ParseInt(segmentId, 10, 64)
		if err != nil || segmentIdInt <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid segment ID"})
			return
		}

		err = db.DeleteSegment(segmentIdInt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Segment deleted successfully", "deleted": segmentId})
	}
}
