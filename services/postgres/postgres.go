package chrono

import (
	"log"
	"database/sql"
)

// Db is the database connection pool 
var Db *sql.DB

// InitDB creates the initial connection pool for the database.
func InitDB(connStr string) {
	db, err := sql.Open("postgres", connStr)
	if(err != nil) {
		log.Panicf("Unable to launch postgres with reason: %s", err)
	}
	Db = db
}