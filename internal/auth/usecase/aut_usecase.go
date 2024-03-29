package usecase

import (
	"github.com/VanLavr/auth/internal/auth/delivery"
	"github.com/VanLavr/auth/internal/models"
)

type authUsecase struct {
	Repository
}

// Repository for working with MongoDB
type Repository interface {
	StoreToken(models.RefreshToken)
	GetToken(string) *models.RefreshToken
}

func New(r Repository) delivery.Usecase {
	return &authUsecase{Repository: r}
}

// check if token is already have been used
// mark provided token as used
// generate token pair
// save new refresh token and mark it as unused
// return the pair
func (a *authUsecase) RefreshTokenPair(token models.RefreshToken) (map[string]string, error) {
	panic("not implemented")
}

func (a *authUsecase) checkIfTokenIsUsed(providedToken models.RefreshToken) bool {
	panic("not implemented")
}
