package validation

import (
	"regexp"
	"strings"

	"github.com/gregbrant2/soda/internal/domain/dataaccess"
	"github.com/gregbrant2/soda/internal/domain/entities"
)

func ValidateServerNew(r dataaccess.ServerRepository, server entities.Server) (bool, map[string]string) {
	var errors = make(map[string]string)
	var namePattern = regexp.MustCompile(`\w`)
	var ipPattern = regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`)
	var portPattern = regexp.MustCompile(`\d+`)

	if strings.TrimSpace(server.Name) == "" {
		errors["Name"] = "Please enter a server name"
	}

	match := namePattern.MatchString(server.Name)
	if !match {
		errors["Name"] = "Please enter a valid server name"
	}

	existing, _ := r.GetServerByName(server.Name)

	if existing != nil {
		errors["Name"] = "A server with this name already exists"
	}

	if strings.TrimSpace(server.Type) == "" {
		errors["Type"] = "Please enter a server type"
	}

	match = namePattern.MatchString(server.Type)
	if !match {
		errors["Type"] = "Please enter a valid server type"
	}

	match = ipPattern.MatchString(server.IpAddress)
	if !match {
		errors["IpAddress"] = "Please enter a valid IPv4 Address"
	}

	if strings.TrimSpace(server.Port) == "" {
		errors["Port"] = "Please enter a server port"
	}

	match = portPattern.MatchString(server.Port)
	if !match {
		errors["Port"] = "Please enter a valid port number"
	}

	if strings.TrimSpace(server.Username) == "" {
		errors["Username"] = "Please enter a username"
	}

	if strings.TrimSpace(server.Password) == "" {
		errors["Password"] = "Please enter a password"
	}

	return len(errors) < 1, errors
}
