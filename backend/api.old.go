package main

/*
import (
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

func startServer2(config Configuration) {
	http.HandleFunc("/api/upload", audioUpload)
	http.HandleFunc("/api/speech-to-text", speechToText)

	fs := http.FileServer(http.Dir("../static"))
	http.Handle("/", fs)

	log.Printf("Listening on port: %s\n", config.Port)

	err := http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func audioUpload(w http.ResponseWriter, r *http.Request) {
	log.Println("receiving data")

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Println(err)
		return
	}

	file, _, err := r.FormFile("recordedAudio")
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	dst, err := os.Create(filepath.Join("tmp", "recording.mp4"))
	_, err = io.Copy(dst, file)
	defer dst.Close()
	if err != nil {
		log.Println(err)
		return
	}
}

func speechToText(w http.ResponseWriter, r *http.Request) {
	command := exec.Command("ffmpeg", "-i", "tmp/recording.mp4", "-ar", "16000", "-ac", "1", "-c:a", "pcm_s16le", "tmp/recording.wav")
	out, err := command.Output()
	if err != nil {
		log.Println(err)
		return
	}

	command = exec.Command("whisper-cli", "-m", "tmp/ggml-base.en.bin", "tmp/recording.wav", "--output-txt", "tmp")
	out, err = command.Output()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(out))

	command = exec.Command("rm", "tmp/recording.mp4", "tmp/recording.wav", "tmp/recording.wav.txt")
	_, err = command.Output()

	w.Write([]byte("You said: " + string(out)))
}
*/
