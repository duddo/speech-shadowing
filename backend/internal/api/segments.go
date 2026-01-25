package api

import (
	"net/http"
	"speech-shadowing/internal"
	"speech-shadowing/internal/database"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getSegments(db *database.Db) gin.HandlerFunc {
	return func(c *gin.Context) {
		exerciseId := c.Param("exercise_id")

		exerciseIdInt, err := strconv.ParseInt(exerciseId, 10, 64)
		if err != nil || exerciseIdInt <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid exercise ID"})
			return
		}

		segments, err := db.GetSegments(exerciseIdInt, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"segments": segments})
	}
}

func postSegment(db *database.Db) gin.HandlerFunc {
	return func(c *gin.Context) {
		exerciseId := c.Param("exercise_id")
		exerciseIdInt, err := strconv.ParseInt(exerciseId, 10, 64)
		if err != nil || exerciseIdInt <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid exercise ID"})
			return
		}

		var newSegment internal.SpeechSegment
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

func putSegment(db *database.Db) gin.HandlerFunc {
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

		var newSegment internal.SpeechSegment
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

func deleteSegment(db *database.Db) gin.HandlerFunc {
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
