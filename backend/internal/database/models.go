package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Db struct {
	connection *sql.DB
}
