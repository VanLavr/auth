package jwt

import (
	"log"
	"net/http"
	"time"

	"github.com/VanLavr/tz1/internal/pkg/config"
	"github.com/golang-jwt/jwt/v5"
)

type JwtMiddleware struct {
	secret string
}

func New(cfg config.Config) *JwtMiddleware {
	return &JwtMiddleware{
		secret: cfg.Secret,
	}
}

func (j *JwtMiddleware) GenerateToken(id string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"guid": id,
		"exp":  time.Now().Add(time.Minute).Unix(),
	})

	stringToken, err := token.SignedString(j.secret)
	if err != nil {
		log.Fatal(err)
	}

	return stringToken
}

func (j *JwtMiddleware) ValidateToken(next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		next(w, r)
	}
}
