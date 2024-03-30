package usecase

import (
	"github.com/VanLavr/auth/internal/auth/delivery"
	"github.com/VanLavr/auth/internal/models"
	"github.com/VanLavr/auth/internal/pkg/config"
)

type authUsecase struct {
	Repository
}

// Repository for working with MongoDB
type Repository interface {
	Connect(*config.Config) error
	CloseConnetion() error

	StoreToken(models.RefreshToken)
	GetToken(string) *models.RefreshToken
	SetCollection() // worth it???
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
func (a *authUsecase) RefreshTokenPair(token models.RefreshToken) (map[string]string, error) {
	panic("not implemented")
}

func (a *authUsecase) checkIfTokenIsUsed(providedToken models.RefreshToken) bool {
	panic("not implemented")
}
