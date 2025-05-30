package viewmodels

import "github.com/gregbrant2/soda/internal/domain/entities"

type Dashboard struct {
	Databases []entities.Database
	Servers   []entities.Server
}

type NewDatabase struct {
	Database    entities.Database
	ServerNames []string
	Errors      map[string]string
}

type DatabaseDetails struct {
	Database entities.Database
	Server   entities.Server
}

type NewServer struct {
	Server *entities.Server
	Errors map[string]string
}
