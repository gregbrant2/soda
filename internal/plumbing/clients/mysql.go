package clients

import (
	"database/sql"
	"log/slog"

	"github.com/go-sql-driver/mysql"
	"github.com/gregbrant2/soda/internal/domain/entities"
	"github.com/gregbrant2/soda/internal/plumbing/utils"
)

type MySqlClient struct {
	DbClient
}

func (c *MySqlClient) CreateDatabase(targetServer entities.Server, name string) error {
	q := "CREATE DATABASE " + name + ";"
	slog.Debug("Creating database with ", "query", q)
	_, err := c.Query(c.connection, q)
	return err
}

func (c *MySqlClient) CreateUser(targetServer entities.Server, database string, name string, password string) error {
	q := "CREATE USER '" + name + "'@'" + targetServer.IpAddress + "' IDENTIFIED BY '" + password + "';"
	slog.Debug("Creating user with ", "query", q)
	_, err := c.Query(c.connection, q)
	if err != nil {
		utils.Fatal("Error creating user", err)
	}

	// GRANT ALL ON db1.* TO 'jeffrey'@'localhost';
	q = "GRANT ALL ON " + database + ".* TO '" + name + "'@'" + targetServer.IpAddress + "' WITH GRANT OPTION;"
	slog.Debug("Granting with ", "query", q)
	_, err = c.Query(c.connection, q)
	if err != nil {
		utils.Fatal("Error granting permissions", err)
	}
	return err
}

func (c *MySqlClient) Connect() (bool, error) {
	cfg := mysql.Config{
		User:   c.Server.Username,
		Passwd: c.Server.Password,
		Net:    "tcp",
		Addr:   c.Server.IpAddress + ":" + c.Server.Port,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		utils.Fatal("Error opening connection", err)
	}

	err = db.Ping()
	if err != nil {
		utils.Fatal("Error pinging connection", err)
		return false, err
	}

	c.connection = db
	return true, nil
}
