package main

type SpeechSegment struct {
	ID         int64  `json:"id"`
	Title      string `json:"title"`
	ExerciseID int64  `json:"exercise_id"`
}

type SpeechExercise struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}
