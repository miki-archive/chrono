package chrono

import (
	"errors"
	models "github.com/mikibot/chrono/models"
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

// InsertEndpoint adds the current endpoint to the database.
func InsertEndpoint(model models.EndpointModel) (error) {
	if(len(model.URL) == 0) {
		return errors.New("cannot insert empty URL in endpoint")
	}

	_, err := Db.Query(
		"INSERT INTO task_endpoints (id, url) VALUES ($1, $2);",
		model.ID, model.URL);
	if(err != nil) {
		return err;
	}
	return nil;
}