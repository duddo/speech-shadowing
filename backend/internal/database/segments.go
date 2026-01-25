package database

import (
	"database/sql"
	"log"
	"speech-shadowing/internal"

	_ "github.com/mattn/go-sqlite3"
)

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

func (db *Db) GetSegments(exerciseID int64, segmentID *int64) ([]internal.SpeechSegment, error) {
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

	var segments []internal.SpeechSegment

	for rows.Next() {
		var id int64
		var title string
		var exerciseID int64
		var spokenText string
		var audioFile string

		err := rows.Scan(&id, &title, &exerciseID, &spokenText, &audioFile)
		if err != nil {
			return nil, err
		}

		segments = append(segments, internal.SpeechSegment{
			ID:         id,
			Title:      title,
			ExerciseID: exerciseID,
			SpokenText: spokenText,
			AudioFile:  audioFile})
	}

	return segments, nil
}

func (db *Db) InsertSegment(exerciseId int64, segment internal.SpeechSegment) (*int64, error) {
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
	res, err := stmt.Exec(segment.Title, exerciseId, segment.SpokenText, segment.AudioFile)
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

func (db *Db) UpdateSegment(segment internal.SpeechSegment) error {
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
	res, err := stmt.Exec(segment.Title, segment.SpokenText, segment.AudioFile, segment.ID)
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

	log.Printf("Deleting segment: %d\n", segmentID)
	res, err := stmt.Exec(segmentID)
	if err != nil {
		return err
	}
	log.Printf("Deleted segment: %d %s\n", segmentID, res)

	return nil
}
