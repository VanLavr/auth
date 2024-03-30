package usecase

import (
	"log"
	"time"

	"github.com/VanLavr/auth/internal/pkg/config"
	e "github.com/VanLavr/auth/internal/pkg/errors"
	"github.com/golang-jwt/jwt/v5"
)

// Stands for generating and token pairs and validating refresh token.
type tokenManager struct {
	secret string
	acExp  time.Duration
	refExp time.Duration
}

func newTokenManager(cfg *config.Config) *tokenManager {
	return &tokenManager{
		secret: cfg.Secret,
		acExp:  cfg.AccessExpTime,
		refExp: cfg.RefreshExpTime,
	}
}

func (j *tokenManager) GenerateTokenPair(id string) map[string]string {
	return map[string]string{
		"access_token":  j.generateAccessToken(id),
		"refresh_token": j.generateRefreshToken(id),
	}
}

// Create claims.
// Sign token.
// Return it.
func (j *tokenManager) generateRefreshToken(id string) string {
	// Create claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"guid": id,
		"exp":  time.Now().Add(time.Second * j.refExp).Unix(),
	})

	// Sign token.
	stringToken, err := token.SignedString([]byte(j.secret))
	if err != nil {
		log.Fatal(err)
	}

	// Return it.
	return stringToken
}

func (j *tokenManager) generateAccessToken(id string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"guid": id,
		"exp":  time.Now().Add(time.Second * j.acExp).Unix(),
	})

	stringToken, err := token.SignedString([]byte(j.secret))
	if err != nil {
		log.Fatal(err)
	}

	return stringToken
}

// Parse token from provided string.
// Check if it is valid.
// Extract guid from claims.
func (j *tokenManager) ValidateRefreshToken(tokenString string) (string, bool) {
	// Parse token from provided string.
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, e.ErrInvalidSigningMethod
		}
		return []byte(j.secret), nil
	})

	// Check if it is valid.
	if err != nil || !token.Valid {
		return "", false
	}

	// Extract guid from claims.
	claims := token.Claims.(jwt.MapClaims)

	id := claims["guid"]
	guid, ok := id.(string)
	if !ok {
		return "", false
	}

	return guid, true
}
