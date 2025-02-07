package dataaccess

import (
	"database/sql"
	"fmt"

	"github.com/gregbrant2/soda/internal/domain/entities"
)

type DatabaseRepository interface {
	AddDatabase(database entities.Database) (int64, error)
	GetDatabaseById(id int64) (entities.Database, error)
	GetDatabaseByName(name string) (*entities.Database, error)
	GetDatabases() ([]entities.Database, error)
}

type MySqlDatabaseRepository struct {
	db *sql.DB
}

func NewMySqlDatabaseRepository(db *sql.DB) *MySqlDatabaseRepository {
	return &MySqlDatabaseRepository{db: db}
}

func (r MySqlDatabaseRepository) AddDatabase(database entities.Database) (int64, error) {
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

func (r MySqlDatabaseRepository) GetDatabaseById(id int64) (entities.Database, error) {
	row := db.QueryRow("SELECT id, name, server_name FROM soda_databases WHERE id=?", id)

	var d entities.Database
	err := row.Scan(&d.Id, &d.Name, &d.Server)
	if err != nil {
		return d, err
	}

	return d, nil
}

func (r MySqlDatabaseRepository) GetDatabaseByName(name string) (*entities.Database, error) {
	row := db.QueryRow("SELECT id, name, server_name FROM soda_databases WHERE name=?", name)

	var db entities.Database
	err := row.Scan(&db.Id, &db.Name, &db.Server)
	if err != nil {
		return nil, err
	}

	return &db, nil
}

func (r MySqlDatabaseRepository) GetDatabases() ([]entities.Database, error) {
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
