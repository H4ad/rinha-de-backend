package main

import (
	"database/sql"
	"os"
	"time"

	_ "github.com/lib/pq"
)

//go:generate sqlboiler --wipe psql

func GetDatabase() *sql.DB {
	dbUrl := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dbUrl)

	if err != nil {
		panic("failed to connect database")
	}

	err = db.Ping()

	if err != nil {
		panic("failed to ping database: " + err.Error())
	}

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(2)
	db.SetConnMaxIdleTime(time.Hour)

	return db
}
