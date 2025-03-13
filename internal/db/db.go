package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

var db *sql.DB

func SetDb() error {
	dbUrl := os.Getenv("DB_URL")
	dbToken := os.Getenv("DB_TOKEN")

	var err error
	db, err = sql.Open("libsql", fmt.Sprintf("%s?authToken=%s", dbUrl, dbToken))

	if err != nil {
		return err
	}

	return db.Ping()
}
