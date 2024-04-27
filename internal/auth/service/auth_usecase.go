// Generate token pair: create tokens (./jwt.go) -> store them in mongo -> return them.
// Refresh token pair: validate provided refresh token (./jwt.go) -> create new token pair -> store them in mongo -> return them.
package usecase

import (
	"context"
	"log/slog"

	"github.com/VanLavr/auth/internal/auth/delivery"
	"github.com/VanLavr/auth/internal/models"
	"github.com/VanLavr/auth/internal/pkg/config"
	e "github.com/VanLavr/auth/internal/pkg/errors"
	"github.com/VanLavr/auth/internal/pkg/hasher"

	"github.com/beevik/guid"
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
	tokenManager := newTokenManager(cfg)
	return &authUsecase{repository: r, tokenManager: tokenManager}
}

// Check if provided token exists.
// Compare provided token with stored hash.
// Validate refresh token jwt.
// Check if this token owned by provided user.
// Validate access and refresh token coherence.
// Generate new token pair.
// Hash refresh token.
// Update token in mongo -> it will replace used tokenstring with new tokenstring.
// Return the pair.
func (a *authUsecase) RefreshTokenPair(ctx context.Context, provided models.RefreshToken, access string) (map[string]any, error) {
	token, err := a.repository.GetToken(ctx, provided)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	// Check if provided token exists.
	if token == nil {
		slog.Error(e.ErrInvalidToken.Error())
		return nil, e.ErrInvalidToken
	}

	// Compare provided token with stored hash.
	if !hasher.Hshr.Validate(token.TokenString, provided.TokenString) {
		slog.Error(e.ErrInvalidToken.Error())
		return nil, e.ErrInvalidToken
	}

	// Validate refresh token jwt. (expired or not)
	guid, valid := a.tokenManager.ValidateRefreshToken(provided.TokenString)
	if !valid {
		slog.Error(e.ErrInvalidToken.Error())
		return nil, e.ErrInvalidToken
	}

	// Check if this token owned by provided user.
	// Check if provided refresh token was already used.
	// TokenString comparison here stands for comparing provided tokenstring and tokenstring from database
	// to inspect if provided token was updated or not. This comparison allow us to check if we can use provided token for refreshing
	// or not even if token was not expired, but was updated (used). So if it not expired, but updated
	// it is an already used token, so we can not use it anymore.
	if token.GUID != guid || token.TokenString != provided.TokenString {
		slog.Error(e.ErrInvalidToken.Error())
		return nil, e.ErrInvalidToken
	}

	// Validate access and refresh token coherence.
	if !a.tokenManager.ValidateTokensCoherence(access, provided.TokenString) {
		slog.Error(e.ErrInvalidToken.Error())
		return nil, e.ErrInvalidToken
	}

	// Generate new token pair.
	tokens := a.tokenManager.GenerateTokenPair(provided.GUID)
	refresh := models.RefreshToken{
		GUID:        provided.GUID,
		TokenString: tokens["refresh_token"],
	}

	// Hash refresh token.
	hash := hasher.Hshr.Encrypt(refresh.TokenString)
	toStoreToken := models.RefreshToken{
		GUID:        provided.GUID,
		TokenString: hash,
	}

	// Update token in mongo -> it will replace used tokenstring with new tokenstring.
	if err := a.repository.UpdateToken(ctx, toStoreToken); err != nil {
		slog.Error(err.Error())
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
// Check if there is old refresh token
// Hash refresh token
// Save hash of refresh token in mongo if there was no old token.
// Update hash of refresh token if there was an old token.
// Return token pair.
func (a *authUsecase) GetNewTokenPair(ctx context.Context, id string) (map[string]any, error) {
	// Validate GUID.
	if !a.validateID(id) {
		slog.Error(e.ErrInvalidGUID.Error())
		return nil, e.ErrInvalidGUID
	}

	// Generate new token pair.
	tokens := a.tokenManager.GenerateTokenPair(id)
	refresh := models.RefreshToken{
		GUID:        id,
		TokenString: tokens["refresh_token"],
	}

	// Check if there is old refresh token
	token, err := a.repository.GetToken(ctx, models.RefreshToken{GUID: id})
	if err != nil && err != e.ErrTokenNotFound {
		slog.Error(err.Error())
		return nil, err
	}

	// Hash refresh token
	hash := hasher.Hshr.Encrypt(refresh.TokenString)
	toStoreToken := models.RefreshToken{
		GUID:        refresh.GUID,
		TokenString: hash,
	}

	// Save refresh token in mongo if there was no old token.
	if token == nil {
		// Save refresh token in mongo.
		if err := a.repository.StoreToken(ctx, toStoreToken); err != nil {
			slog.Error(err.Error())
			return nil, err
		}

		// Return token pair.
		return map[string]any{
			"access_token":  tokens["access_token"],
			"refresh_token": refresh,
		}, nil

		// Update refresh token if there was an old token.
	} else {
		// Update token in mongo -> it will replace used tokenstring with new tokenstring.
		if err := a.repository.UpdateToken(ctx, toStoreToken); err != nil {
			slog.Error(err.Error())
			return nil, err
		}

		// Return the pair.
		return map[string]any{
			"access_token":  tokens["access_token"],
			"refresh_token": refresh,
		}, nil
	}
}

func (a *authUsecase) validateID(id string) bool {
	return guid.IsGuid(id)
}
