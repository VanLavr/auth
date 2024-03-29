package delivery

func (s *Server) BindRoutes() {
	s.httpMux.HandleFunc("GET /getToken/{id}", s.getTokenPair)
	s.httpMux.Handle("GET /restricted", s.jwt.ValidateToken(s.restricted))
}
