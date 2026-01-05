package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func StartServer(config *Configuration, db *Db) {
	router := gin.Default()

	//router.Static("/", config.StaticPath)
	//router.Static("/audio", config.AudioPath)

	router.POST("/api/do-exercise", doExercise(db, config))

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

func doExercise(db *Db, config *Configuration) gin.HandlerFunc {
	return func(c *gin.Context) {
		/* typescript:

		...
		formData.append("audio", file)
		formData.append("payload", JSON.stringify({
		  exercise_id: 2,
		  segment_id: 1,
		  audio_format: "ogg",
		}))
		...
		*/
		payloadStr := c.PostForm("payload")
		if payloadStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing Payload"})
			return
		}
		var payload SegmentSubmit
		err := json.Unmarshal([]byte(payloadStr), &payload)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		file, err := c.FormFile("audio")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upload file"})
			return
		}

		guid := uuid.New()
		filePath := config.TmpPath + "/" + guid.String() + ".mp4" //TODO capire l'estensione da mimetype?

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}

		segmentAnswer, err := Shadowing(payload, filePath, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Segment submitted successfully", "answer": segmentAnswer})
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

		segments, err := db.GetSegments(exerciseIdInt, nil)
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
