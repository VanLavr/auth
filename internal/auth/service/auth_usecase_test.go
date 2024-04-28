package usecase_test

import (
	"context"
	"testing"

	usecase "github.com/VanLavr/auth/internal/auth/service"
	auth_repo_mocks "github.com/VanLavr/auth/internal/mocks/auht/repo"
	"github.com/VanLavr/auth/internal/models"
	"github.com/VanLavr/auth/internal/pkg/config"
	e "github.com/VanLavr/auth/internal/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetNewTokenPair(t *testing.T) {
	repo := &auth_repo_mocks.Repository{}

	repo.On("GetToken", context.Background(), models.RefreshToken{GUID: "67a23ff3-20be-4420-9274-d16f2833d595"}).Return(nil, e.ErrTokenNotFound).Once()
	repo.On("GetToken", context.Background(), models.RefreshToken{GUID: "67a23ff3-20be-4420-9274-d16f2833d656"}).Return(&models.RefreshToken{
		GUID:        "67a23ff3-20be-4420-9274-d16f2833d656",
		TokenString: "some string",
	}, nil).Once()
	repo.On("StoreToken", context.Background(), mock.AnythingOfType("models.RefreshToken")).Return(nil).Once()
	repo.On("UpdateToken", context.Background(), mock.AnythingOfType("models.RefreshToken")).Return(nil).Twice()

	testcases := []struct {
		providedContext    context.Context
		providedID         string
		expectedResultKeys []string
		expectedError      error
		name               string
	}{
		{
			providedContext:    context.Background(),
			providedID:         "67a23ff3-20be-4420-9274-d16f2833d595",
			expectedResultKeys: []string{"access_token", "refresh_token"},
			expectedError:      nil,
			name:               "1",
		},
		{
			providedContext:    context.Background(),
			providedID:         "adsf",
			expectedResultKeys: nil,
			expectedError:      e.ErrInvalidGUID,
			name:               "2",
		},
		{
			providedContext:    context.Background(),
			providedID:         "67a23ff3-20be-4420-9274-d16f2833d656",
			expectedResultKeys: []string{"access_token", "refresh_token"},
			expectedError:      nil,
			name:               "3",
		},
	}

	service := usecase.New(repo, &config.Config{
		Secret:         "asdf",
		AccessExpTime:  120,
		RefreshExpTime: 240,
	})

	for _, tc := range testcases {
		t.Log(tc.name)
		assert := assert.New(t)

		tokens, err := service.GetNewTokenPair(context.Background(), tc.providedID)
		assert.Equal(err, tc.expectedError)
		for k := range tokens {
			if !(k == tc.expectedResultKeys[0] || k == tc.expectedResultKeys[1]) {
				t.Fatalf("unexpected key: %v, test %s", k, tc.name)
			}
		}
	}
}
