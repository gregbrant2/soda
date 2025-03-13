package dataaccess

import "database/sql"

type UnitOfWork struct {
	db      *sql.DB
	DBs     DatabaseRepository
	Servers ServerRepository
}

func NewUow() UnitOfWork {
	db := Initialize()
	defer db.Close()

	dbr := NewMySqlDatabaseRepository(db)
	sr := NewMySqlServerRepository(db)

	return UnitOfWork{
		db:      db,
		DBs:     dbr,
		Servers: sr,
	}
}

func (u *UnitOfWork) BeginTran() (*sql.Tx, error) {
	tx, err := u.db.Begin()
	return tx, err
}
