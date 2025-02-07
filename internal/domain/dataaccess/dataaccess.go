package dataaccess

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/gregbrant2/soda/internal/domain/entities"
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

func GetServerById(id int64) (*entities.Server, error) {
	row := db.QueryRow("SELECT id, name, type, ip_address, port, username, password, (select COUNT(1) from soda_databases WHERE soda_databases.id = soda_servers.id) as 'db_count' FROM soda_servers WHERE id=?", id)

	var server entities.Server
	if err := row.Scan(&server.Id, &server.Name, &server.Type, &server.IpAddress, &server.Port, &server.Username, &server.Password, &server.Databases); err != nil {
		return nil, err
	}

	return &server, nil
}

func GetServerByName(name string) (*entities.Server, error) {
	row := db.QueryRow("SELECT id, name, type, ip_address, port, username, password, (select COUNT(1) from soda_databases WHERE server_name = name) as 'db_count' FROM soda_servers WHERE name=?", name)

	var server entities.Server
	if err := row.Scan(&server.Id, &server.Name, &server.Type, &server.IpAddress, &server.Port, &server.Username, &server.Password, &server.Databases); err != nil {
		return nil, err
	}

	return &server, nil
}
