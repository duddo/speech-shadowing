package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

func startServer() {
	http.HandleFunc("/audio", audioHandler)

	fs := http.FileServer(http.Dir("../static"))
	http.Handle("/", fs)

	fmt.Println("Listening on :8080")

	http.ListenAndServe(":8080", nil)
}

func audioHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("receiving")

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		fmt.Println(err)
		return
	}

	file, _, err := r.FormFile("audio")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	dst, err := os.Create(filepath.Join(".", "recording.webm"))
	_, err = io.Copy(dst, file)
	if err != nil {
		fmt.Println(err)
		return
	}

	command := exec.Command("ffmpeg", "-i", "recording.webm", "-ar", "16000", "-ac", "1", "-c:a", "pcm_s16le", "recording.wav")
	out, err := command.Output()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(out))

	command = exec.Command("/opt/homebrew/Cellar/whisper-cpp/1.8.2/bin/whisper-cli", "-m", "ggml-base.en.bin", "recording.wav", "--output-txt", "")
	out, err = command.Output()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(out))

	w.Write([]byte("You said: " + string(out)))
}
