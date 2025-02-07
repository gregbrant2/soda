package validation

import (
	"testing"

	"github.com/gregbrant2/soda/internal/domain/entities"
)

type FakeServerRepository struct {
	addServerCalled       bool
	addServerResult       int64
	addServerError        error
	getServerByIdCalled   bool
	getServerByIdResult   entities.Server
	getServerByIdError    error
	getServerByNameCalled bool
	getServerByNameResult entities.Server
	getServerByNameError  error
	getServersCalled      bool
	getServersResult      []entities.Server
	getServersError       error
}

func (r FakeServerRepository) AddServer(Server entities.Server) (int64, error) {
	r.addServerCalled = true
	return r.addServerResult, r.addServerError
}

func (r FakeServerRepository) GetServerById(id int64) (entities.Server, error) {
	r.getServerByIdCalled = true
	return r.getServerByIdResult, r.getServerByIdError
}

func (r FakeServerRepository) GetServerByName(name string) (*entities.Server, error) {
	r.getServerByNameCalled = true
	return &r.getServerByNameResult, r.getServerByNameError
}

func (r FakeServerRepository) GetServers() ([]entities.Server, error) {
	r.getServersCalled = true
	return r.getServersResult, r.getServersError
}

func TestValidateServerNewSuccess(t *testing.T) {
	t.FailNow()
}

func TestValidateServerNewEmptyErrors(t *testing.T) {
	t.FailNow()
}

func TestValidateServerNewExistingNameErrors(t *testing.T) {
	t.FailNow()
}
