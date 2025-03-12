package dataaccess

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/gregbrant2/soda/internal/domain/entities"
)

type ServerRepository interface {
	AddServer(server entities.Server) (int64, error)
	GetServerById(id int64) (*entities.Server, error)
	GetServerByName(name string) (*entities.Server, error)
	GetServers() ([]entities.Server, error)
}

type MySqlServerRepository struct {
	db *sql.DB
}

func NewMySqlServerRepository(db *sql.DB) *MySqlServerRepository {
	return &MySqlServerRepository{db: db}
}

func (r MySqlServerRepository) AddServer(server entities.Server) (int64, error) {
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

func (r MySqlServerRepository) GetServers() ([]entities.Server, error) {
	slog.Debug("Getting all servers")
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
		slog.Error("Errror reading server rows", err)
		return servers, err
	}

	return servers, nil
}

func (r MySqlServerRepository) GetServerById(id int64) (*entities.Server, error) {
	row := db.QueryRow("SELECT id, name, type, ip_address, port, username, password, (select COUNT(1) from soda_databases WHERE soda_databases.id = soda_servers.id) as 'db_count' FROM soda_servers WHERE id=?", id)

	var server entities.Server
	if err := row.Scan(&server.Id, &server.Name, &server.Type, &server.IpAddress, &server.Port, &server.Username, &server.Password, &server.Databases); err != nil {
		return nil, err
	}

	return &server, nil
}

func (r MySqlServerRepository) GetServerByName(name string) (*entities.Server, error) {
	row := db.QueryRow("SELECT id, name, type, ip_address, port, username, password, (select COUNT(1) from soda_databases WHERE server_name = name) as 'db_count' FROM soda_servers WHERE name=?", name)

	var server entities.Server
	if err := row.Scan(&server.Id, &server.Name, &server.Type, &server.IpAddress, &server.Port, &server.Username, &server.Password, &server.Databases); err != nil {
		return nil, err
	}

	return &server, nil
}
