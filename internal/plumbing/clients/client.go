package clients

import (
	"database/sql"
	"errors"

	"github.com/gregbrant2/soda/internal/domain/entities"
	"github.com/gregbrant2/soda/internal/plumbing/utils"
)

type IDbClient interface {
	CreateDatabase(targetServer entities.Server, name string) error
	CreateUser(targetServer entities.Server, database string, name string, password string) error
	Connect() (bool, error)
	Ping() error
}

type DbClient struct {
	connection *sql.DB
	Server     entities.Server
}

func (dbc *DbClient) Ping() error {
	if dbc.connection == nil {
		return errors.New("no connection available")
	}

	return dbc.connection.Ping()
}

func (dbc *DbClient) Query(conn *sql.DB, query string) (sql.Result, error) {
	r, err := conn.Exec(query)
	if err != nil {
		utils.Fatal("Error executing query", err, "query", query)
	}

	return r, nil
}
