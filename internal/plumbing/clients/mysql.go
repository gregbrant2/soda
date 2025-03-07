package clients

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/gregbrant2/soda/internal/domain/entities"
)

type MySqlClient struct {
	DbClient
}

func (c *MySqlClient) CreateDatabase(targetServer entities.Server, name string) error {
	conn := connect(targetServer)
	q := "CREATE DATABASE " + name + ";"
	log.Println("Creating database with ", q)
	_, err := c.Query(conn, q)
	return err
}

func (c *MySqlClient) CreateUser(targetServer entities.Server, database string, name string, password string) error {
	conn := connect(targetServer)
	q := "CREATE USER '" + name + "'@'" + targetServer.IpAddress + "' IDENTIFIED BY '" + password + "';"
	log.Println("Creating user with ", q)
	_, err := c.Query(conn, q)
	if err != nil {
		log.Fatal(err)
	}

	// GRANT ALL ON db1.* TO 'jeffrey'@'localhost';
	q = "GRANT ALL ON " + database + ".* TO '" + name + "'@'" + targetServer.IpAddress + "' WITH GRANT OPTION;"
	log.Println("Granting with ", q)
	_, err = c.Query(conn, q)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err) // ?
		return false, err
	}

	c.connection = db
	return true, nil
}

func connect(targetServer entities.Server) *sql.DB {
	cfg := mysql.Config{
		User:   targetServer.Username,
		Passwd: targetServer.Password,
		Net:    "tcp",
		Addr:   targetServer.IpAddress + ":" + targetServer.Port,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db
}
