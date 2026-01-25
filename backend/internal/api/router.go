package api

import (
	"log"
	"speech-shadowing/internal"

	"github.com/gin-gonic/gin"
)

func StartServer(config *internal.Configuration, db *internal.Db) {
	router := gin.Default()

	router.POST("/api/do-exercise", doExercise(db, config))

	router.GET("/api/exercises", getExercises(db))
	router.POST("/api/exercises", postExercise(db))
	router.PUT("/api/exercises/:exercise_id", putExercise(db))
	router.DELETE("/api/exercises/:exercise_id", deleteExercise(db))

	router.GET("/api/exercises/:exercise_id/segments", getSegments(db))
	router.POST("/api/exercises/:exercise_id/segments", postSegment(db))
	router.PUT("/api/exercises/:exercise_id/segments/:segment_id", putSegment(db))
	router.DELETE("/api/exercises/:exercise_id/segments/:segment_id", deleteSegment(db))

	router.Static("/audio", config.AudioPath)

	router.NoRoute(func(c *gin.Context) {
		c.File(config.StaticPath + "/index.html")
	})

	err := router.Run(config.Addr())
	if err != nil {
		log.Fatal(err)
	}
}
