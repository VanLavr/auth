package delivery

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/VanLavr/auth/internal/pkg/config"
	e "github.com/VanLavr/auth/internal/pkg/errors"
	"github.com/VanLavr/auth/internal/pkg/middlewares/jwt"
)

type Server struct {
	httpSrv *http.Server
	httpMux *http.ServeMux
	jwt     *jwt.JwtMiddleware
	Usecase
}

type Usecase interface {
	RefreshToken(string) (map[string]string, error)
	CheckIfTokenIsUsed(string, string) bool
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
		Usecase: u,
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
	var token Refresh
	s.decodeBody(r, &token)
	guid, valid := s.jwt.ValidateRefreshToken(token.Token)
	if !valid {
		fmt.Fprint(w, s.encodeToJSON(Response{
			Error:   e.ErrInvalidToken.Error(),
			Content: nil,
		}))
	}

	if !s.CheckIfTokenIsUsed(guid, token.Token) {
		fmt.Fprint(w, s.encodeToJSON(Response{
			Error:   e.ErrInvalidToken.Error(),
			Content: nil,
		}))
	}
}

func (s *Server) getTokenPair(w http.ResponseWriter, r *http.Request) {
	tokens := s.jwt.GenerateTokenPair(r.PathValue("id"))
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

func (s *Server) decodeBody(r *http.Request, dest *Refresh) {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		log.Println(err)
	}
}
