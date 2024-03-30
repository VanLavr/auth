package jwt

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/VanLavr/auth/internal/pkg/config"
	e "github.com/VanLavr/auth/internal/pkg/errors"
	"github.com/golang-jwt/jwt/v5"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

// Secret is a secret string for token encryption. acExp - access token exparation time,
// refExp - refresh token exparation time.
type JwtMiddleware struct {
	secret string
	acExp  time.Duration
	refExp time.Duration
}

func New(cfg *config.Config) *JwtMiddleware {
	return &JwtMiddleware{
		secret: cfg.Secret,
		acExp:  cfg.AccessExpTime,
		refExp: cfg.RefreshExpTime,
	}
}

// Extract token string from request.
// Parse it.
// Check if it valid or not.
// Call the handler if it's allright.
func (j *JwtMiddleware) ValidateAccessToken(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token string from request.
		tokenString, err := j.ExtractTokenString(r)
		if err != nil {
			fmt.Fprint(w, err.Error())
			log.Println(err)
			return
		}

		// Parse it.
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			_, ok := t.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, e.ErrInvalidSigningMethod
			}
			return []byte(j.secret), nil
		})

		// Check if it valid or not.
		if err != nil || !token.Valid {
			fmt.Fprint(w, e.ErrInvalidToken.Error())
			log.Println(err)
			return
		}

		// Call the handler if it's allright.
		next(w, r)
	})
}

func (j *JwtMiddleware) ExtractTokenString(r *http.Request) (string, error) {
	authHeaders := r.Header.Values("Authorization")
	if len(authHeaders) == 0 {
		return "", e.ErrTokenWasNotProvided
	}

	jwtHeader := authHeaders[0]
	tokenString := jwtHeader[len("Bearer "):]
	if len(tokenString) == 0 {
		return "", e.ErrTokenWasNotProvided
	}

	return tokenString, nil
}
