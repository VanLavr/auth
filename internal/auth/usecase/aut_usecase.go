package usecase

import (
	"context"

	"github.com/VanLavr/auth/internal/auth/delivery"
	"github.com/VanLavr/auth/internal/models"
	"github.com/VanLavr/auth/internal/pkg/config"
)

type authUsecase struct {
	Repository
}

// Repository for working with MongoDB
type Repository interface {
	Connect(context.Context, *config.Config) error
	CloseConnetion(context.Context) error

	// StoreToken() Saves new refresh token and marks it as unused.
	StoreToken(context.Context, models.RefreshToken) error
	// GetToken() requires provided token for getting the token from mongo and check if it was used.
	GetToken(context.Context, models.RefreshToken) (*models.RefreshToken, error)
}

func New(r Repository) delivery.Usecase {
	return &authUsecase{Repository: r}
}

// Check if provided token exists.
// Check if token is already have been used.
// Mark provided token as used.
// Generate new token pair.
// Save new refresh token and mark it as unused.
// Return the pair.
func (a *authUsecase) RefreshTokenPair(ctx context.Context, token models.RefreshToken) (map[string]string, error) {
	panic("not implemented")
}

func (a *authUsecase) checkIfTokenIsUsed(providedToken models.RefreshToken) bool {
	panic("not implemented")
}
