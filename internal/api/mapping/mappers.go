package mapping

import (
	"github.com/gregbrant2/soda/internal/api/dtos"
	"github.com/gregbrant2/soda/internal/domain/entities"
	"github.com/gregbrant2/soda/internal/plumbing/utils"
)

func MapServers(servers []entities.Server) []dtos.Server {
	return utils.Map(servers, func(s entities.Server) dtos.Server { return MapServer(s) })
}

func MapServer(entity entities.Server) dtos.Server {
	return dtos.Server{
		Id:        entity.Id,
		Name:      entity.Name,
		Type:      entity.Type,
		Databases: entity.Databases,
		IpAddress: entity.IpAddress,
		Port:      entity.Port,
		Status:    entity.Status,
	}
}

func MapNewServer(dto dtos.NewServer) entities.Server {
	return entities.Server{
		Name:      dto.Name,
		Type:      dto.Type,
		IpAddress: dto.IpAddress,
		Port:      dto.Port,
		Username:  dto.Username,
		Password:  dto.Password,
	}
}
func MapDatabases(servers []entities.Database) []dtos.Database {
	return utils.Map(servers, func(s entities.Database) dtos.Database { return MapDatabase(s) })
}

func MapDatabase(entity entities.Database) dtos.Database {
	return dtos.Database{
		Id:     entity.Id,
		Name:   entity.Name,
		Server: entity.Server,
	}
}

func MapNewDatabase(entity dtos.NewDatabase) entities.Database {
	return entities.Database{
		Name:   entity.Name,
		Server: entity.Server,
	}
}
