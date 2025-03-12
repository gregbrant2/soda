package dataaccess

import (
	"database/sql"
	"log/slog"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/gregbrant2/soda/internal/plumbing/utils"
)

var db *sql.DB

func Initialize() *sql.DB {
	slog.Info("Initializing System Database")
	cfg := mysql.Config{
		User:   os.Getenv("SODA_SYSTEM_DB_USER"),
		Passwd: os.Getenv("SODA_SYSTEM_DB_PASSWORD"),
		Net:    "tcp",
		Addr:   os.Getenv("SODA_SYSTEM_DB_ADDR"),
		DBName: os.Getenv("SODA_SYSTEM_DB_NAME"),
	}

	slog.Info("Connecting to system database", "address", cfg.Addr, "database", cfg.DBName, "user", cfg.User)

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		utils.Fatal("Connection to system database failed:", err)
	}

	err = db.Ping()
	if err != nil {
		utils.Fatal("Pinging system database failed:", err)
	}

	slog.Info("Connected to system database successfully")
	return db
}
