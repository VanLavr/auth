// Two routes:

// 1) get token pair -> gen pair (usecase) -> save refresh token (repository) -> return pair.

// 2) refresh token pair -> check if provided refresh token is valid (delivery - middleware) ->
// check if token is used (usecase) -> generate new pair (usecase) -> update refresh token (repository) -> return pair.
package delivery

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/VanLavr/auth/internal/models"
	"github.com/VanLavr/auth/internal/pkg/config"
	e "github.com/VanLavr/auth/internal/pkg/errors"
	jwt "github.com/VanLavr/auth/internal/pkg/middlewares/validator"
)

type Server struct {
	httpSrv *http.Server
	httpMux *http.ServeMux
	jwt     *jwt.JwtMiddleware
	u       Usecase
}

// Busyness logic for refreshing tokens e.g.
type Usecase interface {
	RefreshTokenPair(context.Context, models.RefreshToken) (map[string]any, error)
	GetNewTokenPair(context.Context, string) (map[string]any, error)
}

func New(cfg *config.Config, u Usecase) *Server {
	srv := &Server{
		httpSrv: &http.Server{
			Addr:           cfg.Addr,
			ReadTimeout:    cfg.ReadTimeout,
			WriteTimeout:   cfg.WriteTimeout,
			MaxHeaderBytes: cfg.MaxHeaderBytes,
		},
		httpMux: http.NewServeMux(),
		u:       u,
		jwt:     jwt.New(cfg),
	}

	srv.httpSrv.Handler = srv.httpMux
	return srv
}

func (s *Server) Run() error {
	return s.httpSrv.ListenAndServe()
}

func (s *Server) ShutDown(ctx context.Context) error {
	return s.httpSrv.Shutdown(ctx)
}

func (s *Server) refreshToken(w http.ResponseWriter, r *http.Request) {
	var token models.RefreshToken
	s.decodeBody(r, &token)
	guid, valid := s.jwt.ValidateRefreshToken(token.TokenString)
	if !valid {
		fmt.Fprint(w, s.encodeToJSON(Response{
			Error:   e.ErrInvalidToken.Error(),
			Content: nil,
		}))
		return
	}

	if guid != token.GUID {
		fmt.Fprint(w, s.encodeToJSON(Response{
			Error:   e.ErrInvalidToken.Error(),
			Content: nil,
		}))
		return
	}

	data, err := s.u.RefreshTokenPair(r.Context(), token)
	if err != nil {
		fmt.Fprint(w, s.encodeToJSON(Response{
			Error:   err.Error(),
			Content: nil,
		}))
		return
	}

	fmt.Fprint(w, s.encodeToJSON(Response{
		Error:   "",
		Content: data,
	}))
}

func (s *Server) getTokenPair(w http.ResponseWriter, r *http.Request) {
	tokens, err := s.u.GetNewTokenPair(r.Context(), r.PathValue("id"))
	if err != nil {
		fmt.Fprint(w, s.encodeToJSON(Response{
			Error:   e.ErrInternal.Error(),
			Content: nil,
		}))
		return
	}

	fmt.Fprint(w, s.encodeToJSON(Response{
		Error:   "",
		Content: tokens,
	}))
}

func (s *Server) restricted(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, s.encodeToJSON(Response{
		Error:   "",
		Content: "got",
	}))
}

func (s *Server) encodeToJSON(resp Response) string {
	encoded, err := json.Marshal(resp)
	if err != nil {
		log.Fatal(err)
	}

	return string(encoded)
}

func (s *Server) decodeBody(r *http.Request, dest *models.RefreshToken) {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		log.Println(err)
	}
}
