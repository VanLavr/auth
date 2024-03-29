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
		authHeaders := r.Header.Values("Authorization")
		if len(authHeaders) == 0 {
			if _, err := fmt.Fprint(w, e.ErrTokenWasNotProvided.Error()); err != nil {
				log.Println(err)
			}
			return
		}

		jwtHeader := authHeaders[0]
		tokenString := jwtHeader[len("Bearer "):]
		if len(tokenString) == 0 {
			if _, err := fmt.Fprint(w, e.ErrTokenWasNotProvided.Error()); err != nil {
				log.Println(err)
			}
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
			if _, er := fmt.Fprint(w, e.ErrInternal.Error()); er != nil {
				log.Println(er)
			}
			log.Println(err)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			if _, err := fmt.Fprint(w, e.ErrTokenWasNotProvided.Error()); err != nil {
				log.Println(err)
			}
			return
		}

		exp := claims["exp"]
		expTime, err := j.validateExpTime(exp)
		if err != nil {
			if _, err := fmt.Fprint(w, err.Error()); err != nil {
				log.Println(err)
			}
			return
		}

		if expTime < time.Now().Unix() {
			if _, err := fmt.Fprint(w, e.ErrTokenExpired.Error()); err != nil {
				log.Println(err)
			}
			return
		}

		next(w, r)
	})
}

func (j *JwtMiddleware) validateExpTime(exp any) (int64, error) {
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
