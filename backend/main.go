package main

import "log"

func main() {
	log.Printf("Speech Shadowing %s", Version)

	config := LoadConfig()

	db, err := Connect(config)
	if err != nil {
		log.Fatal(err)
		return
	}

	// StartServer(config)

	err = db.Close()
	if err != nil {
		return
	}
}
