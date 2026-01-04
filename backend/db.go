package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Db struct {
	connection *sql.DB
}

func NewDb(config *Configuration) (*Db, error) {
	db, err := sql.Open("sqlite3", config.DataPath+"/data.db")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	log.Println("Creating speech_exercises")
	createExercisesStmt := `CREATE TABLE IF NOT EXISTS speech_exercises (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL
	);`
	result, err := db.Exec(createExercisesStmt)
	if err != nil {
		log.Printf("CREATE speech_exercises %q: %s\n", err, createExercisesStmt)
		return nil, err
	}
	log.Printf("Created speech_exercises: %s\n", result)

	log.Println("Creating speech_segments")
	createSegmentsStmt := `CREATE TABLE IF NOT EXISTS speech_segments (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		exercise_id TEXT NOT NULL,
		FOREIGN KEY (exercise_id) REFERENCES speech_exercises(id) ON DELETE RESTRICT ON UPDATE CASCADE
	);`
	result, err = db.Exec(createSegmentsStmt)
	if err != nil {
		log.Printf("CREATE speech_segments %q: %s\n", err, createSegmentsStmt)
		return nil, err
	}
	log.Printf("Created speech_segments: %s\n", result)

	return &Db{
		connection: db,
	}, nil
}

func (db *Db) Close() error {
	return db.connection.Close()
}

/*************************************
 * Exercise
 */

func (db *Db) GetExercises() ([]SpeechExercise, error) {
	rows, err := db.connection.Query("SELECT * FROM speech_exercises;")
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	var exercises []SpeechExercise

	for rows.Next() {
		var id string
		var title string
		var exerciseID string
		var exercise SpeechExercise

		err := rows.Scan(&id, &title, &exerciseID, &exercise)
		if err != nil {
			return nil, err
		}

		exercises = append(exercises, exercise)
	}

	return exercises, nil
}

func (db *Db) InsertExercise(exercise SpeechExercise) error {
	stmt, err := db.connection.Prepare("INSERT INTO speech_exercises(id, title) VALUES (?, ?);")
	if err != nil {
		return err
	}

	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(stmt)

	log.Printf("Inserting exercise: %s %s\n", exercise.ID, exercise.Title)
	res, err := stmt.Exec(exercise.ID, exercise.Title)
	if err != nil {
		return err
	}
	log.Printf("Inserted exercise: %s %s\n", exercise.ID, res)

	return nil
}

func (db *Db) UpdateExercise(exercise SpeechExercise) error {
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

	log.Printf("Updating exercise: %s %s\n", exercise.ID, exercise.Title)
	res, err := stmt.Exec(exercise.Title, exercise.ID)
	if err != nil {
		return err
	}
	log.Printf("Updated exercise %s: %s\n", exercise.ID, res)

	return nil
}

func (db *Db) DeleteExercise(exercise SpeechExercise) error {
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

	log.Printf("Deleting exercise: %s\n", exercise.ID)
	res, err := stmt.Exec(exercise.ID)
	if err != nil {
		return err
	}
	log.Printf("Deleted exercise %s: %s\n", exercise.ID, res)

	return nil
}

/*************************************
 * Segments
 */

func (db *Db) GetSegments(id int) ([]SpeechSegment, error) {
	rows, err := db.connection.Query("SELECT * FROM speech_segments WHERE id=?;", id)
	if err != nil {

		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	var segments []SpeechSegment

	for rows.Next() {
		var id string
		var title string
		var exerciseID string
		var exercise SpeechSegment

		err := rows.Scan(&id, &title, &exerciseID, &exercise)
		if err != nil {
			return nil, err
		}

		segments = append(segments, exercise)
	}

	return segments, nil
}

func (db *Db) InsertSegment(segment SpeechSegment) error {
	stmt, err := db.connection.Prepare("INSERT INTO speech_segments(id, title)")
	if err != nil {
		return err
	}

	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(stmt)

	log.Printf("Inserting segment: %s %s\n", segment.ID, segment.Title)
	res, err := stmt.Exec(segment.ID, segment.Title)
	if err != nil {
		return err
	}
	log.Printf("Inserted segment: %s %s\n", segment.ID, res)

	return nil
}

func (db *Db) UpdateSegment(segment SpeechSegment) error {
	stmt, err := db.connection.Prepare("UPDATE speech_segments SET title = ? WHERE id = ?")
	if err != nil {
		return err
	}

	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(stmt)

	log.Printf("Updating segment: %s %s\n", segment.ID, segment.Title)
	res, err := stmt.Exec(segment.Title, segment.ID)
	if err != nil {
		return err
	}
	log.Printf("Updated segment: %s %s\n", segment.ID, res)

	return nil
}

func (db *Db) DeleteSegment(segment SpeechSegment) error {
	stmt, err := db.connection.Prepare("DELETE FROM speech_segments WHERE id = ?")
	if err != nil {
		return err
	}

	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(stmt)

	log.Printf("Deleting segment: %s %s\n", segment.ID, segment.Title)
	res, err := stmt.Exec(segment.ID)
	if err != nil {
		return err
	}
	log.Printf("Deleted segment: %s %s\n", segment.ID, res)

	return nil
}
