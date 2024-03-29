package jwt

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/VanLavr/tz1/internal/pkg/config"
	e "github.com/VanLavr/tz1/internal/pkg/errors"
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

func (j *JwtMiddleware) ValidateToken(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
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
			return j.secret, nil
		})
		if err != nil {
			fmt.Fprint(w, e.ErrInternal.Error())
			log.Println(err)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			fmt.Fprint(w, e.ErrTokenWasNotProvided.Error())
			log.Println(err)
			return
		}

		if !j.validateExpTime(claims["exp"]) {
			fmt.Fprint(w, e.ErrInvalidToken)
		}

		next(w, r)
	})
}

func (j *JwtMiddleware) validateExpTime(exp any) bool {
	expTime, err := j.extractExpTime(exp)
	if err != nil {
		log.Println(err)
		return false
	}

	if expTime < time.Now().Unix() {
		return false
	}

	return true
}

func (j *JwtMiddleware) extractExpTime(exp any) (int64, error) {
	t, ok := exp.(float64)
	if !ok {
		return 0, e.ErrInvalidToken
	}

	time := int64(t)
	if time == 0 {
		return 0, e.ErrInvalidToken
	}

	return time, nil
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
