package validation

import (
	"testing"

	"github.com/gregbrant2/soda/internal/domain/entities"
)

type FakeDatabaseRepository struct {
	addDatabaseCalled       bool
	addDatabaseResult       int64
	addDatabaseError        error
	getDatabaseByIdCalled   bool
	getDatabaseByIdResult   entities.Database
	getDatabaseByIdError    error
	getDatabaseByNameCalled bool
	getDatabaseByNameResult entities.Database
	getDatabaseByNameError  error
	getDatabasesCalled      bool
	getDatabasesResult      []entities.Database
	getDatabasesError       error
}

func (r FakeDatabaseRepository) AddDatabase(database entities.Database) (int64, error) {
	r.addDatabaseCalled = true
	return r.addDatabaseResult, r.addDatabaseError
}

func (r FakeDatabaseRepository) GetDatabaseById(id int64) (*entities.Database, error) {
	r.getDatabaseByIdCalled = true
	return &r.getDatabaseByIdResult, r.getDatabaseByIdError
}

func (r FakeDatabaseRepository) GetDatabaseByName(name string) (*entities.Database, error) {
	r.getDatabaseByNameCalled = true
	return &r.getDatabaseByNameResult, r.getDatabaseByNameError
}

func (r FakeDatabaseRepository) GetDatabases() ([]entities.Database, error) {
	r.getDatabasesCalled = true
	return r.getDatabasesResult, r.getDatabasesError
}

func TestValidateDatabaseNewSuccess(t *testing.T) {
	r := &FakeDatabaseRepository{}

	database := entities.Database{
		Name:   "Foo",
		Server: "Bar",
	}

	valid, errors := ValidateDatabaseNew(r, database)
	if !valid || len(errors) > 0 {
		t.Fatal("Database should have been valid")
	}
}

func TestValidateDatabaseNewEmptyErrors(t *testing.T) {
	r := &FakeDatabaseRepository{}

	database := entities.Database{}

	valid, errors := ValidateDatabaseNew(r, database)
	if valid || len(errors) < 1 {
		t.Fatal("Database should have been invalid", errors)
	}
}

func TestValidateDatabaseNewExistingNameErrors(t *testing.T) {
	r := &FakeDatabaseRepository{
		getDatabaseByNameResult: entities.Database{
			Name:   "Name1",
			Server: "Server1",
		},
	}

	database := entities.Database{
		Name:   "Name1",
		Server: "Server1",
	}

	valid, errors := ValidateDatabaseNew(r, database)
	if valid || len(errors) < 1 {
		t.Fatal("Database should have been invalid due to name clash")
	}
}

func TestValidateDatabaseNewExistingNameDifferentServerValid(t *testing.T) {
	r := &FakeDatabaseRepository{
		getDatabaseByNameResult: entities.Database{
			Name:   "Name1",
			Server: "Server1",
		},
	}

	database := entities.Database{
		Name:   "Name1",
		Server: "Server2",
	}

	valid, errors := ValidateDatabaseNew(r, database)
	if !valid || len(errors) > 0 {
		t.Fatal("Database should have been valid due to different sever name")
	}
}
