package logic

import (
	"os/exec"
	"speech-shadowing/internal"
	"speech-shadowing/internal/database"
)

func Shadowing(segmentSubmit internal.SegmentSubmit, audioFile string, db *database.Db) (*internal.SegmentAnswer, error) {
	transcript, err := speechToText(audioFile)
	if err != nil {
		return nil, err
	}

	speechSegment, err := db.GetSegments(segmentSubmit.ExerciseID, &segmentSubmit.SegmentID)
	if err != nil || len(speechSegment) == 0 {
		return nil, err
	}

	rating := similarity(speechSegment[0].SpokenText, *transcript)

	answer := internal.SegmentAnswer{
		ExerciseID: segmentSubmit.ExerciseID,
		SegmentID:  segmentSubmit.SegmentID,
		Transcript: *transcript,
		Rating:     rating,
	}

	err = db.InsertAnswer(answer)
	if err != nil {
		return nil, err
	}

	return &answer, nil
}

func speechToText(audioFile string) (*string, error) {
	command := exec.Command("ffmpeg", "-i", audioFile, "-ar", "16000", "-ac", "1", "-c:a", "pcm_s16le", audioFile+".wav")
	out, err := command.Output()
	if err != nil {
		return nil, err
	}

	command = exec.Command("whisper-cli", "-m", "tmp/ggml-base.en.bin", audioFile+".wav", "--output-txt", "tmp")
	out, err = command.Output()
	if err != nil {
		return nil, err
	}

	//command = exec.Command("rm", audioFile, audioFile+".wav", audioFile+".wav.txt")
	//_, err = command.Output()

	stringOut := string(out)
	return &stringOut, nil
}

func similarity(a, b string) float32 {
	if a == "" && b == "" {
		return 1.0
	}
	dist := levenshtein(a, b)
	maxLen := max(len(a), len(b))
	return 1.0 - float32(dist)/float32(maxLen)
}

func levenshtein(a, b string) int {
	la := len(a)
	lb := len(b)

	dp := make([][]int, la+1)
	for i := range dp {
		dp[i] = make([]int, lb+1)
	}

	for i := 0; i <= la; i++ {
		dp[i][0] = i
	}
	for j := 0; j <= lb; j++ {
		dp[0][j] = j
	}

	for i := 1; i <= la; i++ {
		for j := 1; j <= lb; j++ {
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}

			dp[i][j] = min(
				dp[i-1][j]+1,
				dp[i][j-1]+1,
				dp[i-1][j-1]+cost,
			)
		}
	}

	return dp[la][lb]
}
