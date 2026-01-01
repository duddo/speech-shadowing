package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Db struct {
	connection *sql.DB
}

func Connect(config *Configuration) (*Db, error) {
	db, err := sql.Open("sqlite3", config.DataPath+"/data.db")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	sqlStmt := `create table foo (id integer not null primary key, name text);
		delete from foo;`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return nil, err
	}

	return &Db{
		connection: db,
	}, nil
}

func (db *Db) Close() error {
	return db.connection.Close()
}
