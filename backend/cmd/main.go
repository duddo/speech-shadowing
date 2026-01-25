package main

import (
	"log"
	"speech-shadowing/internal"
	"speech-shadowing/internal/api"
)

func main() {
	log.Printf("Speech Shadowing %s", internal.Version)

	log.Println("Loading configuration")
	config := internal.NewConfiguration()

	log.Println("Starting database")
	db, err := internal.NewDb(config)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println("Starting REST API")
	api.StartServer(config, db)

	log.Println("Exiting...")
	err = db.Close()
	if err != nil {
		log.Fatal(err)
		return
	}
}
