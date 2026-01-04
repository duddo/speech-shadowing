package main

type SpeechExercise struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

type SpeechSegment struct {
	ID         int64  `json:"id"`
	Title      string `json:"title"`
	ExerciseID int64  `json:"exercise_id"`
	SpokenText string `json:"spoken_text"`
	AudioID    string `json:"audio_id"`
}
