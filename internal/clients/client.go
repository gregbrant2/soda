package clients

import (
	"database/sql"
	"log"

	"github.com/gregbrant2/soda/internal/entities"
)

type DbClient interface {
	CreateDatabase(targetServer entities.Server, name string) error
	CreateUser(targetServer entities.Server, name string, password string) error
}

func query(conn *sql.DB, query string) (sql.Result, error) {
	r, err := conn.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	return r, nil
}
