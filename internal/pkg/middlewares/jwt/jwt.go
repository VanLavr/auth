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

func (j *JwtMiddleware) GenerateTokenPair(id string) map[string]string {
	return map[string]string{
		"access":  j.generateAccessToken(id),
		"refresh": j.generateRefreshToken(id),
	}
}

func (j *JwtMiddleware) generateRefreshToken(id string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"guid": id,
		"exp":  time.Now().Add(time.Second * j.refExp).Unix(),
	})

	stringToken, err := token.SignedString([]byte(j.secret))
	if err != nil {
		log.Fatal(err)
	}
	return stringToken
}

func (j *JwtMiddleware) generateAccessToken(id string) string {
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

func (j *JwtMiddleware) ValidateAccessToken(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := j.extractTokenString(r)
		if err != nil {
			fmt.Fprint(w, err.Error())
			log.Println(err)
			return
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			_, ok := t.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, e.ErrInvalidSigningMethod
			}
			return []byte(j.secret), nil
		})
		if err != nil || !token.Valid {
			fmt.Fprint(w, e.ErrInvalidToken.Error())
			log.Println(err)
			return
		}

		next(w, r)
	})
}

func (j *JwtMiddleware) extractTokenString(r *http.Request) (string, error) {
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

func (j *JwtMiddleware) ValidateRefreshToken(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, e.ErrInvalidSigningMethod
		}
		return []byte(j.secret), nil
	})
	if err != nil || !token.Valid {
		return false
	}

	return true
}
