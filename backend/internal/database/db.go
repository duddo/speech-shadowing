package database

import (
	"database/sql"
	"log"
	"speech-shadowing/internal"

	_ "github.com/mattn/go-sqlite3"
)

func NewDb(config *internal.Configuration) (*Db, error) {
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
