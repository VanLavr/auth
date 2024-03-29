package delivery

import (
	"net/http"

	"github.com/VanLavr/tz1/internal/pkg/config"
)

type Server struct {
	httpSrv *http.Server
	httpMux *http.ServeMux
	usecase Usecase
}

type Usecase interface {
	GetNewTokenPair()
	RefreshToken()
}

func New(cfg *config.Config, u Usecase) *Server {
	return &Server{
		httpSrv: &http.Server{
			Addr:           cfg.Addr,
			ReadTimeout:    cfg.ReadTimeout,
			WriteTimeout:   cfg.WriteTimeout,
			MaxHeaderBytes: cfg.MaxHeaderBytes,
		},
		httpMux: http.NewServeMux(),
		usecase: u,
	}
}
