package validation

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/gregbrant2/soda/internal/domain/dataaccess"
	"github.com/gregbrant2/soda/internal/domain/entities"
)

func ValidateDatabaseNew(dbr dataaccess.DatabaseRepository, sr dataaccess.ServerRepository, database entities.Database) (bool, map[string]string) {
	var errors = make(map[string]string)
	var namePattern = regexp.MustCompile(`\w`)

	if strings.TrimSpace(database.Name) == "" {
		errors["Name"] = "Please enter a database name"
	}

	match := namePattern.MatchString(database.Name)
	if !match {
		errors["Name"] = "Please enter a valid database name"
	}

	if strings.TrimSpace(database.Server) == "" {
		errors["Server"] = "Please specify a server name"
	} else {
		server, err := sr.GetServerByName(database.Server)
		if err != nil {
			log.Println("Error fetching server")
			log.Println(err)
		}

		if server != nil {
			existing, _ := dbr.GetDatabaseByName(database.Name)
			// TODO: Figure out the pattern here.
			//       err would be !nil but that's right when existing is nil
			//       but what if there's a real error, not just "database not found"?
			// if err != nil {
			// 	log.Fatal(err)
			// }

			if existing != nil {
				if existing.Server == database.Server {
					errors["Name"] = fmt.Sprintf("Database '%v' already exists on server '%v", database.Name, database.Server)
				}
			}
		} else {
			errors["Server"] = "Server not found"
		}
	}

	return len(errors) < 1, errors
}
