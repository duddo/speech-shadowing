package internal

type SpeechExercise struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

type SpeechSegment struct {
	ID         int64  `json:"id"`
	Title      string `json:"title"`
	ExerciseID int64  `json:"exercise_id"`
	SpokenText string `json:"spoken_text"`
	AudioFile  string `json:"audio_file"`
}

type SegmentSubmit struct {
	ExerciseID  int64  `json:"exercise_id"`
	SegmentID   int64  `json:"segment_id"`
	AudioFormat string `json:"audio_format"`
}

type SegmentAnswer struct {
	ExerciseID int64   `json:"exercise_id"`
	SegmentID  int64   `json:"segment_id"`
	Transcript string  `json:"transcript"`
	Rating     float32 `json:"rating"`
	Timestamp  int64   `json:"timestamp"`
}

type BulkSubmit struct {
	Exercise SpeechExercise  `json:"exercise"`
	Segments []SpeechSegment `json:"segments"`
}
