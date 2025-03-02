package validation

import (
	"testing"

	"github.com/gregbrant2/soda/internal/domain/entities"
	"github.com/stretchr/testify/assert"
)

type FakeDatabaseRepository struct {
	addDatabaseCalled       bool
	addDatabaseResult       int64
	addDatabaseError        error
	getDatabaseByIdCalled   bool
	getDatabaseByIdResult   *entities.Database
	getDatabaseByIdError    error
	getDatabaseByNameCalled bool
	getDatabaseByNameResult *entities.Database
	getDatabaseByNameError  error
	getDatabasesCalled      bool
	getDatabasesResult      []entities.Database
	getDatabasesError       error
}

func (r *FakeDatabaseRepository) AddDatabase(database entities.Database) (int64, error) {
	r.addDatabaseCalled = true
	return r.addDatabaseResult, r.addDatabaseError
}

func (r *FakeDatabaseRepository) GetDatabaseById(id int64) (*entities.Database, error) {
	r.getDatabaseByIdCalled = true
	return r.getDatabaseByIdResult, r.getDatabaseByIdError
}

func (r *FakeDatabaseRepository) GetDatabaseByName(name string) (*entities.Database, error) {
	r.getDatabaseByNameCalled = true
	return r.getDatabaseByNameResult, r.getDatabaseByNameError
}

func (r *FakeDatabaseRepository) GetDatabases() ([]entities.Database, error) {
	r.getDatabasesCalled = true
	return r.getDatabasesResult, r.getDatabasesError
}

func TestValidateDatabaseNewSuccess(t *testing.T) {
	dbr := &FakeDatabaseRepository{}
	sr := &FakeServerRepository{}
	sr.getServerByNameResult = ValidServer()

	database := entities.Database{
		Name:   "Foo",
		Server: "Bar",
	}

	valid, errors := ValidateDatabaseNew(dbr, sr, database)
	assert.True(t, valid)
	assert.Len(t, errors, 0)
}

func TestValidateDatabaseNewEmptyErrors(t *testing.T) {
	r := &FakeDatabaseRepository{}
	sr := &FakeServerRepository{}

	database := entities.Database{}

	valid, errors := ValidateDatabaseNew(r, sr, database)
	if valid || len(errors) < 1 {
		t.Fatal("Database should have been invalid", errors)
	}
}

func TestValidateDatabaseNewExistingNameErrors(t *testing.T) {
	dbr := &FakeDatabaseRepository{
		getDatabaseByNameResult: &entities.Database{
			Name:   "Name1",
			Server: "Server1",
		},
	}
	sr := FakeServerRepository{}

	database := entities.Database{
		Name:   "Name1",
		Server: "Server1",
	}

	valid, errors := ValidateDatabaseNew(dbr, sr, database)
	if valid || len(errors) < 1 {
		t.Fatal("Database should have been invalid due to name clash")
	}
}

func TestValidateDatabaseNewExistingNameDifferentServerValid(t *testing.T) {
	dbr := &FakeDatabaseRepository{
		getDatabaseByNameResult: &entities.Database{
			Name:   "Name1",
			Server: "Server1",
		},
	}
	sr := FakeServerRepository{
		getServerByNameResult: ValidServer(),
	}

	database := entities.Database{
		Name:   "Name1",
		Server: "Server2",
	}

	valid, errors := ValidateDatabaseNew(dbr, sr, database)

	assert.True(t, valid, "Database should have been valid due to different sever name")
	assert.Len(t, errors, 0)
}
