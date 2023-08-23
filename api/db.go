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

tryagain:
	err = db.Ping()

	if err != nil {
		time.Sleep(2 * time.Second)
		goto tryagain
	}

	db.SetMaxIdleConns(10)
	db.SetConnMaxIdleTime(time.Hour)

	return db
}
