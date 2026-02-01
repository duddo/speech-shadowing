package api

import (
	"net/http"
	"os"
	"os/exec"
	"speech-shadowing/internal"
	"speech-shadowing/internal/database"

	"github.com/gin-gonic/gin"
)

func postBulk(db *database.Db) gin.HandlerFunc {
	return func(c *gin.Context) {
		var bulk internal.BulkSubmit

		err := c.BindJSON(&bulk)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		insertedID, err := db.InsertExercise(bulk.Exercise)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		for _, segment := range bulk.Segments {
			if segment.GenerateAudio {
				err := generateAudioSample(segment)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
					return
				}
			}

			_, err = db.InsertSegment(*insertedID, segment)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Bulk submission added successfully"})
	}
}

func generateAudioSample(segment internal.SpeechSegment) error {
	audioFilePath := "./public" + segment.AudioFile

	if audioFilePath != "" && fileExists(audioFilePath) {
		return nil
	}

	sayCmd := "say -v Fred -o ./tmp/sample.aiff \"" + segment.SpokenText + "\"" +
		" && lame ./tmp/sample.aiff " + audioFilePath +
		" && rm ./tmp/sample.aiff"

	execCmd := exec.Command("zsh", "-c", sayCmd)
	err := execCmd.Run()

	if err != nil {
		return err
	}

	return nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}
