package usecase

import (
	"context"
	"testing"
	"time"

	auth_repo_mocks "github.com/VanLavr/auth/internal/mocks/auht/repo"
	"github.com/VanLavr/auth/internal/models"
	"github.com/VanLavr/auth/internal/pkg/config"
	e "github.com/VanLavr/auth/internal/pkg/errors"
	"github.com/VanLavr/auth/internal/pkg/hasher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Testcases:
// 1) provide valid id and unexisting token
// 2) provide an invalid id
// 3) provide valid id and an existing token
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

	service := New(repo, &config.Config{
		Secret:         "asdf",
		AccessExpTime:  3 * time.Second,
		RefreshExpTime: 5 * time.Second,
	})

	for _, tc := range testcases {
		t.Log(tc.name)
		assert := assert.New(t)

		tokens, err := service.GetNewTokenPair(context.Background(), tc.providedID)
		assert.Equal(tc.expectedError, err)
		for k := range tokens {
			if !(k == tc.expectedResultKeys[0] || k == tc.expectedResultKeys[1]) {
				t.Fatalf("unexpected key: %v, test %s", k, tc.name)
			}
		}
	}
}

// Testcases:
// 1) provide valid refresh and valid access tokens
// 2) provide expired refresh token
// 3) provide invalid refresh token
// 4) provide used and not expired refresh token
func TestRefreshTokenPair(t *testing.T) {
	tokenMngr := newTokenManager(&config.Config{
		Secret:         "ggg",
		AccessExpTime:  3 * time.Second,
		RefreshExpTime: 5 * time.Second,
	})

	// 1) provide valid refresh and valid access tokens
	tokens := tokenMngr.GenerateTokenPair("67a23ff3-20be-4420-9274-d16f2833d595")
	actk, ok := tokens["access_token"]
	if !ok {
		t.Fatal("can not generate")
	}
	reftk, ok := tokens["refresh_token"]
	if !ok {
		t.Fatal("can not generate")
	}

	hashedRefreshToken595 := hasher.Hshr.Encrypt(reftk)

	// 2) provide expired refresh token
	secondTokens := tokenMngr.GenerateTokenPair("67a23ff3-20be-4420-9274-d16f2833d656")
	actk2, ok := secondTokens["access_token"]
	if !ok {
		t.Fatal("can not generate")
	}
	reftk2, ok := tokens["refresh_token"]
	if !ok {
		t.Fatal("can not generate")
	}

	hashedRefreshToken656 := hasher.Hshr.Encrypt(reftk2)

	// 3) provide invalid refresh token
	invalidRefreshToken := "asdfv"
	invalidTokenHash := hasher.Hshr.Encrypt(invalidRefreshToken)

	// 4) provide used and not expired refresh token
	thirdTokens := tokenMngr.GenerateTokenPair("67a23ff3-20be-4420-9274-d16f2833d656")
	actk3, ok := thirdTokens["access_token"]
	if !ok {
		t.Fatal("can not generate")
	}
	reftk3, ok := tokens["refresh_token"]
	if !ok {
		t.Fatal("can not generate")
	}

	hashedRefreshTokenUsed := hasher.Hshr.Encrypt(reftk3)

	repo := &auth_repo_mocks.Repository{}

	repo.On("GetToken", context.Background(), models.RefreshToken{
		GUID:        "67a23ff3-20be-4420-9274-d16f2833d595",
		TokenString: reftk,
	}).
		Return(&models.RefreshToken{
			GUID:        "67a23ff3-20be-4420-9274-d16f2833d595",
			TokenString: hashedRefreshToken595,
		}, nil).Once()

	repo.On("GetToken", context.Background(), models.RefreshToken{
		GUID:        "67a23ff3-20be-4420-9274-d16f2833d656",
		TokenString: reftk2,
	}).
		Return(&models.RefreshToken{
			GUID:        "67a23ff3-20be-4420-9274-d16f2833d656",
			TokenString: hashedRefreshToken656,
		}, nil).Once()

	repo.On("GetToken", context.Background(), models.RefreshToken{
		GUID:        "67a23ff3-20be-4420-9274-d16f2833d777",
		TokenString: invalidRefreshToken,
	}).
		Return(&models.RefreshToken{
			GUID:        "67a23ff3-20be-4420-9274-d16f2833d777",
			TokenString: invalidTokenHash,
		}, nil).Once()

	repo.On("GetToken", context.Background(), models.RefreshToken{
		GUID:        "67a23ff3-20be-4420-9274-d16f2833d656",
		TokenString: reftk3,
	}).
		Return(&models.RefreshToken{
			GUID:        "67a23ff3-20be-4420-9274-d16f2833d656",
			TokenString: hashedRefreshTokenUsed,
		}, nil).Once()

	repo.On("UpdateToken", context.Background(), mock.AnythingOfType("models.RefreshToken")).Return(nil).Twice()

	service := New(repo, &config.Config{
		Secret:         "ggg",
		AccessExpTime:  3 * time.Second,
		RefreshExpTime: 5 * time.Second,
	})

	testcases := []struct {
		providedContext           context.Context
		providedRefreshToken      models.RefreshToken
		providedAccessTokenString string
		expectedResultKeys        []string
		expectedError             error
		timeToWaitTillExpires     time.Duration
		name                      string
	}{
		{
			providedContext: context.Background(),
			providedRefreshToken: models.RefreshToken{
				GUID:        "67a23ff3-20be-4420-9274-d16f2833d595",
				TokenString: reftk,
			},
			providedAccessTokenString: actk,
			expectedResultKeys:        []string{"access_token", "refresh_token"},
			expectedError:             nil,
			timeToWaitTillExpires:     time.Duration(1 * time.Second),
			name:                      "1",
		},
		{
			providedContext: context.Background(),
			providedRefreshToken: models.RefreshToken{
				GUID:        "67a23ff3-20be-4420-9274-d16f2833d656",
				TokenString: reftk2,
			},
			providedAccessTokenString: actk2,
			expectedResultKeys:        nil,
			expectedError:             e.ErrInvalidToken,
			timeToWaitTillExpires:     time.Duration(7 * time.Second),
			name:                      "2",
		},
		{
			providedContext: context.Background(),
			providedRefreshToken: models.RefreshToken{
				GUID:        "67a23ff3-20be-4420-9274-d16f2833d777",
				TokenString: invalidRefreshToken,
			},
			providedAccessTokenString: actk,
			expectedResultKeys:        nil,
			expectedError:             e.ErrInvalidToken,
			timeToWaitTillExpires:     time.Duration(1 * time.Second),
			name:                      "3",
		},
		{
			providedContext: context.Background(),
			providedRefreshToken: models.RefreshToken{
				GUID:        "67a23ff3-20be-4420-9274-d16f2833d656",
				TokenString: reftk3,
			},
			providedAccessTokenString: actk3,
			expectedResultKeys:        nil,
			expectedError:             e.ErrInvalidToken,
			timeToWaitTillExpires:     time.Duration(1 * time.Second),
			name:                      "4",
		},
	}

	for i := 0; i < 3; i++ {
		t.Log("||||||||||||||||||||", testcases[i].name, testcases[i].expectedError)
		assert := assert.New(t)
		<-time.After(testcases[i].timeToWaitTillExpires)

		tokens, err := service.RefreshTokenPair(testcases[i].providedContext, testcases[i].providedRefreshToken, testcases[i].providedAccessTokenString)

		assert.Equal(testcases[i].expectedError, err)
		for k := range tokens {
			if !(k == testcases[i].expectedResultKeys[0] || k == testcases[i].expectedResultKeys[1]) {
				t.Fatalf("unexpected key: %v, test %s", k, testcases[i].name)
			}
		}
	}
}
