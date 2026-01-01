package main

import "log"

func main() {
	log.Printf("Speech Shadowing %s", Version)

	config := LoadConfig()

	Connect(config)

	// StartServer(config)
}
