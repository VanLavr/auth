package delivery

import (
	_ "github.com/VanLavr/auth/docs"
	httpSwagger "github.com/swaggo/http-swagger" // http-swagger middleware
)

func (s *Server) BindRoutes() {
	s.httpMux.HandleFunc("GET /getToken/{id}", s.getTokenPair)
	s.httpMux.Handle("GET /restricted", s.jwt.ValidateAccessToken(s.restricted))
	s.httpMux.HandleFunc("POST /refreshToken", s.refreshToken)
	s.httpMux.HandleFunc("GET /swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))
}
