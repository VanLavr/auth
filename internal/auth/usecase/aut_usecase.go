// Generate token pair: create tokens (./jwt.go) -> store them in mongo -> return them.
// Refresh token pair: validate provided refresh token (./jwt.go) -> create new token pair -> store them in mongo -> return them.
package usecase

import (
	"context"

	"github.com/VanLavr/auth/internal/auth/delivery"
	"github.com/VanLavr/auth/internal/models"
	"github.com/VanLavr/auth/internal/pkg/config"
	e "github.com/VanLavr/auth/internal/pkg/errors"
)

type authUsecase struct {
	tokenManager *tokenManager
	repository   Repository
}

// Repository for working with MongoDB
type Repository interface {
	Connect(context.Context, *config.Config) error
	CloseConnetion(context.Context) error

	// StoreToken() Saves new refresh token and marks it as unused.
	StoreToken(context.Context, models.RefreshToken) error
	// GetToken() requires provided token for getting the token from mongo and check if it was used.
	GetToken(context.Context, models.RefreshToken) (*models.RefreshToken, error)
	// UpdateToken() is used to mark tokens as used
	UpdateToken(context.Context, models.RefreshToken) error
}

func New(r Repository, cfg *config.Config) delivery.Usecase {
	tokenManager := newTokenGenerator(cfg)
	return &authUsecase{repository: r, tokenManager: tokenManager}
}

// Check if provided token exists.
// Check if this token owned by provided user.
// Generate new token pair.
// Update token in mongo -> it will replace used tokenstring with new tokenstring.
// Return the pair.
func (a *authUsecase) RefreshTokenPair(ctx context.Context, provided models.RefreshToken) (map[string]any, error) {
	token, err := a.repository.GetToken(ctx, provided)
	if err != nil {
		return nil, err
	}

	// Check if provided token exists.
	if token == nil {
		return nil, e.ErrInvalidToken
	}

	// Check if this token owned by provided user.
	if token.GUID != provided.GUID {
		return nil, e.ErrInvalidToken
	}

	// Generate new token pair.
	tokens := a.tokenManager.GenerateTokenPair(provided.GUID)
	refresh := models.RefreshToken{
		GUID:        provided.GUID,
		TokenString: tokens["refresh_token"],
	}

	// Update token in mongo -> it will replace used tokenstring with new tokenstring.
	if err := a.repository.UpdateToken(ctx, refresh); err != nil {
		return nil, err
	}

	// Return the pair.
	return map[string]any{
		"access_token":  tokens["access_token"],
		"refresh_token": refresh,
	}, nil
}

// Validate GUID.
// Generate new token pair.
// Save refresh token in mongo.
// Return token pair.
func (a *authUsecase) GetNewTokenPair(ctx context.Context, id string) (map[string]any, error) {
	// Validate GUID.
	if !a.validateID(id) {
		return nil, e.ErrInvalidGUID
	}

	// Generate new token pair.
	tokens := a.tokenManager.GenerateTokenPair(id)
	refresh := models.RefreshToken{
		GUID:        id,
		TokenString: tokens["refresh_token"],
	}

	// Save refresh token in mongo.
	if err := a.repository.StoreToken(ctx, refresh); err != nil {
		return nil, err
	}

	// Return token pair.
	return map[string]any{
		"access_token":  tokens["access_token"],
		"refresh_token": refresh,
	}, nil
}

func (a *authUsecase) validateID(id string) bool {
	panic("not implemented")
}
