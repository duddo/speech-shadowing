package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Db struct {
	connection *sql.DB
}

func NewDb(config *Configuration) (*Db, error) {
	db, err := sql.Open("sqlite3", config.DataPath+"/data.sqlite")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	err = createExercises(db)
	if err != nil {
		return nil, err
	}

	err = createSegments(db)
	if err != nil {
		return nil, err
	}

	err = createAnswers(db)
	if err != nil {
		return nil, err
	}

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

func (db *Db) GetExercises() ([]SpeechExercise, error) {
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

	var exercises []SpeechExercise

	for rows.Next() {
		var id int64
		var title string

		err := rows.Scan(&id, &title)
		if err != nil {
			return nil, err
		}

		exercises = append(exercises, SpeechExercise{ID: id, Title: title})
	}

	return exercises, nil
}

func (db *Db) InsertExercise(exercise SpeechExercise) (*int64, error) {
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

	log.Printf("Deleting exercise: %d %s\n", exerciseID)
	res, err := stmt.Exec(exerciseID)
	if err != nil {
		return err
	}
	log.Printf("Deleted exercise %d: %s\n", exerciseID, res)

	return nil
}

/*************************************
 * Segments
 */

func createSegments(db *sql.DB) error {
	log.Println("Creating speech_segments")
	createSegmentsStmt := `CREATE TABLE IF NOT EXISTS speech_segments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		exercise_id TEXT NOT NULL,
		spoken_text TEXT NOT NULL,
		audio_id TEXT NOT NULL,
		FOREIGN KEY (exercise_id) REFERENCES speech_exercises(id) ON DELETE CASCADE ON UPDATE CASCADE
	);`
	result, err := db.Exec(createSegmentsStmt)
	if err != nil {
		log.Printf("CREATE speech_segments %q: %s\n", err, createSegmentsStmt)
		return err
	}

	log.Printf("Created speech_segments: %s\n", result)

	return nil
}

func (db *Db) GetSegments(exerciseID int64, segmentID *int64) ([]SpeechSegment, error) {
	var rows *sql.Rows
	var err error

	if segmentID == nil {
		rows, err = db.connection.Query(`SELECT id, title, exercise_id, spoken_text, audio_id
			FROM speech_segments WHERE exercise_id=?;`, exerciseID)
	} else {
		rows, err = db.connection.Query(`SELECT id, title, exercise_id, spoken_text, audio_id
			FROM speech_segments WHERE exercise_id=? AND id=?;`, exerciseID, segmentID)
	}

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
		var id int64
		var title string
		var exerciseID int64
		var spokenText string
		var audioID string

		err := rows.Scan(&id, &title, &exerciseID, &spokenText, &audioID)
		if err != nil {
			return nil, err
		}

		segments = append(segments, SpeechSegment{
			ID:         id,
			Title:      title,
			ExerciseID: exerciseID,
			SpokenText: spokenText,
			AudioID:    audioID})
	}

	return segments, nil
}

func (db *Db) InsertSegment(exerciseId int64, segment SpeechSegment) (*int64, error) {
	stmt, err := db.connection.Prepare(`INSERT INTO 
    	speech_segments(id, title, exercise_id, spoken_text, audio_id) 
		VALUES (NULL, ?, ?, ?, ?);`)
	if err != nil {
		return nil, err
	}

	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(stmt)

	log.Printf("Inserting segment: %d %s\n", segment.ID, segment.Title)
	res, err := stmt.Exec(segment.Title, exerciseId, segment.SpokenText, segment.AudioID)
	if err != nil {
		return nil, err
	}

	insertedID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	log.Printf("Inserted segment id %d, title \"%s\"\n", insertedID, segment.Title)

	return &insertedID, nil
}

func (db *Db) UpdateSegment(segment SpeechSegment) error {
	stmt, err := db.connection.Prepare(`UPDATE speech_segments 
			SET title = ?, spoken_text = ?, audio_id = ?
		    WHERE id = ?`)
	if err != nil {
		return err
	}

	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(stmt)

	log.Printf("Updating segment: %d %s\n", segment.ID, segment.Title)
	res, err := stmt.Exec(segment.Title, segment.SpokenText, segment.AudioID, segment.ID)
	if err != nil {
		return err
	}
	log.Printf("Updated segment: %d %s\n", segment.ID, res)

	return nil
}

func (db *Db) DeleteSegment(segmentID int64) error {
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

	log.Printf("Deleting segment: %d %s\n", segmentID)
	res, err := stmt.Exec(segmentID)
	if err != nil {
		return err
	}
	log.Printf("Deleted segment: %d %s\n", segmentID, res)

	return nil
}

/*************************************
 * Shadowing Answers
 */

func createAnswers(db *sql.DB) error {
	log.Println("Creating shadowing_answers")
	createAnswersStmt := `CREATE TABLE IF NOT EXISTS shadowing_answers(
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	exercise_id INTEGER,
    	segment_id INTEGER,
    	transcript TEXT,
    	rating REAL,
    	datetime DATETIME,
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

func (db *Db) InsertAnswer(answer SegmentAnswer) error {
	stmt, err := db.connection.Prepare(`INSERT INTO shadowing_answers(
		id, exercise_id, segment_id, transcript, rating, datetime)
		VALUES (NULL, ?, ?, ?, ?, ?);`)
	if err != nil {
		return err
	}

	datetime := time.Unix(answer.Timestamp, 0)

	log.Printf("Inserting answer")
	_, err = stmt.Exec(answer.ExerciseID, answer.SegmentID, answer.Transcript, answer.Rating, datetime)
	if err != nil {
		return err
	}

	return nil
}
