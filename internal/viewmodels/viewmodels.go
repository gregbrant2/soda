package viewmodels

import "github.com/gregbrant2/soda/internal/entities"

type Dashboard struct {
	Databases []entities.Database
	Servers   []entities.Server
}
