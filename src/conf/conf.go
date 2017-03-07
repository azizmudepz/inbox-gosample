package conf

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type DBBundle struct {
	Database  *sql.DB
	Database2 *sql.DB
}

var DB DBBundle

func InitDB(
	database *string,
	database2 *string,
) {
	db, err := sql.Open("postgres", *database)
	if err != nil {
		log.Fatal("database not available", err)
	}

	db2, err := sql.Open("postgres", *database2)
	if err != nil {
		log.Fatal("database2 not available")
	}

	DB = DBBundle{
		db,
		db2,
	}
}
