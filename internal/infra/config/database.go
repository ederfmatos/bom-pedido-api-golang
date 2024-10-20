package config

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"log/slog"
)

func Database(driver, url string) *sql.DB {
	database, err := sql.Open(driver, url)
	failOnError(err, "Failed to connect to database")
	err = database.Ping()
	failOnError(err, "Failed to ping database")
	slog.Info("Connected to database successfully")
	return database
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
