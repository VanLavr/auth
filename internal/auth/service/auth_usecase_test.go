package usecase_test

import (
	"context"
	"testing"

	auth_repo_mocks "github.com/VanLavr/auth/internal/mocks/auht/repo"
	"github.com/VanLavr/auth/internal/models"
	e "github.com/VanLavr/auth/internal/pkg/errors"
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
	}{
		{
			providedContext:    context.Background(),
			providedID:         "67a23ff3-20be-4420-9274-d16f2833d595",
			expectedResultKeys: []string{"access_token", "refresh_token"},
			expectedError:      nil,
		},
		{
			providedContext:    context.Background(),
			providedID:         "adsf",
			expectedResultKeys: nil,
			expectedError:      e.ErrInvalidGUID,
		},
		{
			providedContext:    context.Background(),
			providedID:         "67a23ff3-20be-4420-9274-d16f2833d656",
			expectedResultKeys: []string{"access_token", "refresh_token"},
			expectedError:      nil,
		},
	}
}
