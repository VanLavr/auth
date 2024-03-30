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
	MarkToken(models.RefreshToken) error
}

func New(r Repository, cfg *config.Config) delivery.Usecase {
	tokenManager := newTokenGenerator(cfg)
	return &authUsecase{repository: r, tokenManager: tokenManager}
}

// Check if provided token exists.
// Check if token is already have been used.
// Check if this token owned by provided user.
// Mark provided token as used.
// Generate new token pair.
// Save new refresh token and mark it as unused.
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

	a.repository.MarkToken(provided)

	panic("not implemented")
}

// Generate new token pair.
// Save refresh token in mongo.
// Return token pair.
func (a *authUsecase) GetNewTokenPair(ctx context.Context, id string) (map[string]any, error) {
	// Generate new token pair.
	tokens := a.tokenManager.GenerateTokenPair(id)
	refresh := models.RefreshToken{
		GUID:        id,
		TokenString: tokens["refresh"],
	}

	// Save refresh token in mongo.
	if err := a.repository.StoreToken(ctx, refresh); err != nil {
		return nil, err
	}

	// Return token pair.
	return map[string]any{
		"access":  tokens["access"],
		"refresh": refresh,
	}, nil
}

func (a *authUsecase) checkIfTokenIsUsed(providedToken models.RefreshToken) bool {
	panic("not implemented")
}
