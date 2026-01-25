package api

import (
	"encoding/json"
	"net/http"
	"speech-shadowing/internal"
	"speech-shadowing/internal/database"
	"speech-shadowing/internal/logic"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func doExercise(db *database.Db, config *internal.Configuration) gin.HandlerFunc {
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
		var payload internal.SegmentSubmit
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

		segmentAnswer, err := logic.Shadowing(payload, filePath, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Segment submitted successfully", "answer": segmentAnswer})
	}
}
