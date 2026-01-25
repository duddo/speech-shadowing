package database

import (
	"database/sql"
	"log"
	"speech-shadowing/internal"

	_ "github.com/mattn/go-sqlite3"
)

func createExercises(db *sql.DB) error {
	log.Println("Creating speech_exercises")
	createExercisesStmt := `CREATE TABLE IF NOT EXISTS speech_exercises (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL
	);`
	result, err := db.Exec(createExercisesStmt)
	if err != nil {
		log.Printf("CREATE speech_exercises %q: %s\n", err, createExercisesStmt)
		return err
	}

	log.Printf("Created speech_exercises: %s\n", result)

	return nil
}

func (db *Db) GetExercises() ([]internal.SpeechExercise, error) {
	rows, err := db.connection.Query("SELECT id, title FROM speech_exercises;")
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	var exercises []internal.SpeechExercise

	for rows.Next() {
		var id int64
		var title string

		err := rows.Scan(&id, &title)
		if err != nil {
			return nil, err
		}

		exercises = append(exercises, internal.SpeechExercise{ID: id, Title: title})
	}

	return exercises, nil
}

func (db *Db) InsertExercise(exercise internal.SpeechExercise) (*int64, error) {
	stmt, err := db.connection.Prepare("INSERT INTO speech_exercises(id, title) VALUES (NULL, ?);")
	if err != nil {
		return nil, err
	}

	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(stmt)

	log.Printf("Inserting exercise: \"%s\"\n", exercise.Title)
	res, err := stmt.Exec(exercise.Title)
	if err != nil {
		return nil, err
	}

	insertedID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	log.Printf("Inserted exercise id %d, title \"%s\"\n", insertedID, exercise.Title)

	return &insertedID, nil
}

func (db *Db) UpdateExercise(exercise internal.SpeechExercise) error {
	stmt, err := db.connection.Prepare("UPDATE speech_exercises SET title = ? WHERE id = ?;")
	if err != nil {
		return err
	}

	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(stmt)

	log.Printf("Updating exercise: %d %s\n", exercise.ID, exercise.Title)
	res, err := stmt.Exec(exercise.Title, exercise.ID)
	if err != nil {
		return err
	}
	log.Printf("Updated exercise %d: %s\n", exercise.ID, res)

	return nil
}

func (db *Db) DeleteExercise(exerciseID int64) error {
	stmt, err := db.connection.Prepare("DELETE FROM speech_exercises WHERE id = ?;")
	if err != nil {
		return err
	}

	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(stmt)

	log.Printf("Deleting exercise: %d\n", exerciseID)
	res, err := stmt.Exec(exerciseID)
	if err != nil {
		return err
	}
	log.Printf("Deleted exercise %d: %s\n", exerciseID, res)

	return nil
}
