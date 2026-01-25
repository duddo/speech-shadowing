package api

import (
	"net/http"
	"speech-shadowing/internal"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getExercises(db *internal.Db) gin.HandlerFunc {
	return func(c *gin.Context) {
		exercises, err := db.GetExercises()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"exercises": exercises})
	}
}

func postExercise(db *internal.Db) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newExercise internal.SpeechExercise

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

func putExercise(db *internal.Db) gin.HandlerFunc {
	return func(c *gin.Context) {
		exerciseId := c.Param("exercise_id")

		exerciseIdInt, err := strconv.ParseInt(exerciseId, 10, 64)
		if err != nil || exerciseIdInt <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid exercise ID"})
			return
		}

		var newExercise internal.SpeechExercise

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

func deleteExercise(db *internal.Db) gin.HandlerFunc {
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
