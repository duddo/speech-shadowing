package internal

type SpeechExercise struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type SpeechSegment struct {
	ID            int64  `json:"id"`
	Title         string `json:"title"`
	ExerciseID    int64  `json:"exercise_id"`
	SpokenText    string `json:"spoken_text"`
	AudioFile     string `json:"audio_file"`
	GenerateAudio bool   `json:"generate_audio"`
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
	Datetime   int64   `json:"datetime"`
}

type BulkSubmit struct {
	Exercise SpeechExercise  `json:"exercise"`
	Segments []SpeechSegment `json:"segments"`
}
