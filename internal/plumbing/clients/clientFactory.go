package clients

import (
	"errors"

	"github.com/gregbrant2/soda/internal/domain/entities"
)

func CreateServer(server entities.Server) (IDbClient, error) {
	switch server.Type {
	case "mysql":
		client := &MySqlClient{DbClient{Server: server}}
		client.Connect()
		return client, nil
	}

	return nil, errors.New("Borked")
}
