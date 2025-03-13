package validation

import (
	"testing"

	"github.com/gregbrant2/soda/internal/domain/entities"
	"github.com/stretchr/testify/assert"
)

func TestValidateServerNewSuccess(t *testing.T) {
	uow, _, _ := fakeDataAccess()
	server := ValidServer()

	success, errors := ValidateServerNew(uow, *server)

	assert.True(t, success)
	assert.Len(t, errors, 0)
}

func TestValidateServerNewEmptyErrors(t *testing.T) {
	server := entities.Server{}
	uow, _, _ := fakeDataAccess()

	success, errors := ValidateServerNew(uow, server)

	assert.False(t, success)
	assert.Len(t, errors, 6)
}

func TestValidateServerNewExistingNameErrors(t *testing.T) {
	server := ValidServer()
	uow, _, sr := fakeDataAccess()

	sr.getServerByNameResult = &entities.Server{
		Name: server.Name,
	}

	success, errors := ValidateServerNew(uow, *server)

	assert.False(t, success)
	assert.Len(t, errors, 1)
	assert.Contains(t, errors, "Name")
}

type FakeServerRepository struct {
	addServerCalled       bool
	addServerResult       int64
	addServerError        error
	getServerByIdCalled   bool
	getServerByIdResult   *entities.Server
	getServerByIdError    error
	getServerByNameCalled bool
	getServerByNameResult *entities.Server
	getServerByNameError  error
	getServersCalled      bool
	getServersResult      []entities.Server
	getServersError       error
}

func (r FakeServerRepository) AddServer(Server entities.Server) (int64, error) {
	r.addServerCalled = true
	return r.addServerResult, r.addServerError
}

func (r FakeServerRepository) GetServerById(id int64) (*entities.Server, error) {
	r.getServerByIdCalled = true
	return r.getServerByIdResult, r.getServerByIdError
}

func (r FakeServerRepository) GetServerByName(name string) (*entities.Server, error) {
	r.getServerByNameCalled = true
	return r.getServerByNameResult, r.getServerByNameError
}

func (r FakeServerRepository) GetServers() ([]entities.Server, error) {
	r.getServersCalled = true
	return r.getServersResult, r.getServersError
}
