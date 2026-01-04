package main

import "log"

func main() {
	log.Printf("Speech Shadowing %s", Version)

	log.Println("Loading configuration")
	config := NewConfiguration()

	log.Println("Starting database")
	db, err := NewDb(config)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println("Starting REST API")
	//StartServer(config)

	log.Println("Exiting...")
	err = db.Close()
	if err != nil {
		log.Fatal(err)
		return
	}
}
