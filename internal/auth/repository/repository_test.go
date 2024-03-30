package repository_test

import (
	"context"
	"testing"

	"github.com/VanLavr/auth/internal/auth/repository"
	"github.com/VanLavr/auth/internal/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	testCase := struct {
		ctx           context.Context
		cfg           *config.Config
		ExpectedError error
	}{
		ctx:           context.TODO(),
		cfg:           &config.Config{DBName: "auth", CollectionName: "refreshtokens", Mongo: "mongodb://127.0.0.1:27017"},
		ExpectedError: nil,
	}

	assert := assert.New(t)

	repo := repository.New(testCase.cfg)

	err := repo.Connect(testCase.ctx, testCase.cfg)
	assert.Equal(err, testCase.ExpectedError)
}

// Mocks here?
func TestGetToken(t *testing.T) {

}

func TestStoreToken(t *testing.T) {

}

func TestUpdateToken(t *testing.T) {

}
