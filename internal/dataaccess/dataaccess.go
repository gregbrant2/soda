package dataaccess

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/gregbrant2/soda/internal/entities"
)

var db *sql.DB

func Initialize() *sql.DB {
	cfg := mysql.Config{
		User:   "root",
		Passwd: "password",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "soda",
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(err)
	}

	return db
}

func AddDatabase(database entities.Database) (int64, error) {
	res, err := db.Exec("INSERT INTO soda_databases (name, server_name) VALUES (?, ?)", database.Name, database.Server)
	if err != nil {
		return 0, fmt.Errorf("addDatabase: %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addDatabase: %v", err)
	}

	return id, nil
}

func GetDatabaseById(id int64) (entities.Database, error) {
	row := db.QueryRow("SELECT id, name, server_name FROM soda_databases WHERE id=?", id)

	var d entities.Database
	err := row.Scan(&d.Id, &d.Name, &d.Server)
	if err != nil {
		return d, err
	}

	return d, nil
}

func GetDatabaseByName(name string) (entities.Database, error) {
	row := db.QueryRow("SELECT id, name, server_name FROM soda_databases WHERE name=?", name)

	var d entities.Database
	err := row.Scan(&d.Id, &d.Name, &d.Server)
	if err != nil {
		return d, err
	}

	return d, nil
}

func GetDatabases() ([]entities.Database, error) {
	rows, err := db.Query("SELECT id, name, server_name FROM soda_databases")
	if err != nil {
		return nil, err
	}

	var databases []entities.Database
	for rows.Next() {
		var d entities.Database
		if err := rows.Scan(&d.Id, &d.Name, &d.Server); err != nil {
			return databases, err
		}
		databases = append(databases, d)
	}

	if err = rows.Err(); err != nil {
		return databases, err
	}

	return databases, nil
}

func AddServer(server entities.Server) (int64, error) {
	res, err := db.Exec("INSERT INTO soda_servers (name, ip_address, port, type, username, password) VALUES (?, ?, ?, ?, ?, ?)", server.Name, server.IpAddress, server.Port, server.Type, server.Username, server.Password)
	if err != nil {
		return 0, fmt.Errorf("addServer: %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addServer: %v", err)
	}

	return id, nil
}

func GetServers() ([]entities.Server, error) {
	log.Println("Getting all servers")
	rows, err := db.Query("SELECT id, name, type, ip_address, port FROM soda_servers")
	if err != nil {
		return nil, err
	}

	var servers []entities.Server
	for rows.Next() {
		var s entities.Server
		if err := rows.Scan(&s.Id, &s.Name, &s.Type, &s.IpAddress, &s.Port); err != nil {
			return servers, err
		}
		servers = append(servers, s)
	}

	if err = rows.Err(); err != nil {
		log.Println(err)
		return servers, err
	}

	return servers, nil
}

func GetServerById(id int64) (entities.Server, error) {
	row := db.QueryRow("SELECT id, name, type, ip_address, port, username, password, (select COUNT(1) from soda_databases WHERE soda_databases.id = soda_servers.id) as 'db_count' FROM soda_servers WHERE id=?", id)

	var server entities.Server
	err := row.Scan(&server.Id, &server.Name, &server.Type, &server.IpAddress, &server.Port, &server.Username, &server.Password, &server.Databases)
	if err != nil {
		return server, err
	}

	return server, nil
}

func GetServerByName(name string) (entities.Server, error) {
	row := db.QueryRow("SELECT id, name, type, ip_address, port, username, password, (select COUNT(1) from soda_databases WHERE server_name = name) as 'db_count' FROM soda_servers WHERE name=?", name)

	var server entities.Server
	err := row.Scan(&server.Id, &server.Name, &server.Type, &server.IpAddress, &server.Port, &server.Username, &server.Password, &server.Databases)
	if err != nil {
		return server, err
	}

	return server, nil
}
