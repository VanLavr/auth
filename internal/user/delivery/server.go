package delivery

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/VanLavr/tz1/internal/pkg/config"
	"github.com/VanLavr/tz1/internal/pkg/middlewares/jwt"
)

type Server struct {
	httpSrv *http.Server
	httpMux *http.ServeMux
	jwt     *jwt.JwtMiddleware
	Usecase
}

type Usecase interface {
	GetNewTokenPair()
	RefreshToken()
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

func (s *Server) getTokenPair(w http.ResponseWriter, r *http.Request) {
	token := s.jwt.GenerateToken(r.PathValue("id"))
	fmt.Fprint(w, s.encodeToJSON(Response{
		Error:   "",
		Content: token,
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
