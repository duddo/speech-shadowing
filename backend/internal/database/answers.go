package database

import (
	"database/sql"
	"log"
	"speech-shadowing/internal"

	_ "github.com/mattn/go-sqlite3"
)

func createAnswers(db *sql.DB) error {
	log.Println("Creating shadowing_answers")
	createAnswersStmt := `CREATE TABLE IF NOT EXISTS shadowing_answers(
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	exercise_id INTEGER,
    	segment_id INTEGER,
    	transcript TEXT,
    	rating REAL,
    	datetime DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (exercise_id) REFERENCES speech_exercises(id),
        FOREIGN KEY (segment_id) REFERENCES speech_segments(id)
    );`
	result, err := db.Exec(createAnswersStmt)
	if err != nil {
		log.Printf("Error creating shadowing_answers: %s\n", err)
		return err
	}

	log.Printf("Created shadowing_answers: %s\n", result)

	return nil
}

func (db *Db) InsertAnswer(answer internal.SegmentAnswer) error {
	stmt, err := db.connection.Prepare(`INSERT INTO shadowing_answers(
		id, exercise_id, segment_id, transcript, rating)
		VALUES (NULL, ?, ?, ?, ?);`)
	if err != nil {
		return err
	}

	log.Printf("Inserting answer")
	_, err = stmt.Exec(answer.ExerciseID, answer.SegmentID, answer.Transcript, answer.Rating)
	if err != nil {
		return err
	}

	return nil
}
