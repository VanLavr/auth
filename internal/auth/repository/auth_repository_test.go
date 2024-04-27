package repository_test

import (
	"context"
	"log"
	"testing"

	"github.com/VanLavr/auth/internal/auth/repository"
	"github.com/VanLavr/auth/internal/models"
	"github.com/VanLavr/auth/internal/pkg/config"
	e "github.com/VanLavr/auth/internal/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestStoreToken(t *testing.T) {
	testcases := []struct {
		arg           models.RefreshToken
		expectedError error
	}{
		{
			arg: models.RefreshToken{
				GUID:        "adsf",
				TokenString: "adsf",
			},
			expectedError: nil,
		},
	}

	config := &config.Config{
		Mongo:          "mongodb://127.0.0.1:27017",
		DBName:         "auth_test",
		CollectionName: "users_tokens",
	}
	repo := repository.New(config)
	repo.Connect(context.Background(), config)
	defer func() {
		if err := repo.CloseConnetion(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	for _, tc := range testcases {
		assert := assert.New(t)
		err := repo.StoreToken(context.TODO(), tc.arg)
		log.Println(err)
		assert.Equal(tc.expectedError, err)
	}
}

func TestUpdateToken(t *testing.T) {
	testcases := []struct {
		tokenToStore  models.RefreshToken
		tokenToUpdate models.RefreshToken
		expectedError error
	}{
		{
			tokenToStore: models.RefreshToken{
				GUID:        "adsf",
				TokenString: "adsf",
			},
			tokenToUpdate: models.RefreshToken{
				GUID:        "adsf",
				TokenString: "updated",
			},
			expectedError: nil,
		},
		{
			tokenToStore: models.RefreshToken{
				GUID:        "adsf",
				TokenString: "adsf",
			},
			tokenToUpdate: models.RefreshToken{
				GUID:        "333",
				TokenString: "won't update",
			},
			expectedError: e.ErrUserNotFound,
		},
	}

	config := &config.Config{
		Mongo:          "mongodb://127.0.0.1:27017",
		DBName:         "auth_test",
		CollectionName: "users_tokens",
	}
	repo := repository.New(config)
	repo.Connect(context.Background(), config)
	defer func() {
		if err := repo.CloseConnetion(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	for _, tc := range testcases {
		assert := assert.New(t)
		err := repo.StoreToken(context.Background(), tc.tokenToStore)
		fatalOnErr(err)

		err = repo.UpdateToken(context.Background(), tc.tokenToUpdate)

		assert.Equal(tc.expectedError, err)
	}
}

func TestGetToken(t *testing.T) {
	testcases := []struct {
		tokenToStore   models.RefreshToken
		provided       models.RefreshToken
		expectedResult *models.RefreshToken
		expectedError  error
	}{
		{
			tokenToStore: models.RefreshToken{
				GUID:        "adsf",
				TokenString: "adsf",
			},
			provided: models.RefreshToken{
				GUID:        "adsf",
				TokenString: "",
			},
			expectedResult: &models.RefreshToken{
				GUID:        "adsf",
				TokenString: "adsf",
			},
			expectedError: nil,
		},
		{
			tokenToStore: models.RefreshToken{
				GUID:        "adsf",
				TokenString: "adsf",
			},
			provided: models.RefreshToken{
				GUID:        "333",
				TokenString: "",
			},
			expectedResult: nil,
			expectedError:  e.ErrTokenNotFound,
		},
	}

	config := &config.Config{
		Mongo:          "mongodb://127.0.0.1:27017",
		DBName:         "auth_test",
		CollectionName: "users_tokens",
	}
	repo := repository.New(config)
	repo.Connect(context.Background(), config)
	defer func() {
		if err := repo.CloseConnetion(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	for _, tc := range testcases {
		assert := assert.New(t)
		err := repo.StoreToken(context.Background(), tc.tokenToStore)
		fatalOnErr(err)

		token, err := repo.GetToken(context.Background(), tc.provided)

		assert.Equal(tc.expectedError, err)
		assert.Equal(tc.expectedResult, token)
	}
}

func fatalOnErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
