package dataaccess

import (
	"database/sql"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func Initialize() *sql.DB {
	log.Println("Initializing System Database")
	cfg := mysql.Config{
		User:   os.Getenv("SODA_SYSTEM_DB_USER"),
		Passwd: os.Getenv("SODA_SYSTEM_DB_PASSWORD"),
		Net:    "tcp",
		Addr:   os.Getenv("SODA_SYSTEM_DB_ADDR"),
		DBName: os.Getenv("SODA_SYSTEM_DB_NAME"),
	}

	log.Printf("Connecting to %s/%s as %s", cfg.Addr, cfg.DBName, cfg.User)

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Println("Connection to system database failed:")
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Println("Pinging system database failed:")
		log.Fatal(err)
	}

	log.Println("Connected to system database successfully")
	return db
}
