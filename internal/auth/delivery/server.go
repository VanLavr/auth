// Two routes:

// 1) get token pair -> gen pair (usecase) -> save refresh token (repository) -> return pair.

// 2) refresh token pair -> check if provided refresh token is valid (usecase) ->
// check if token is used (usecase) -> generate new pair (usecase) -> update refresh token (repository) -> return pair.
package delivery

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
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
	RefreshTokenPair(context.Context, models.RefreshToken, string) (map[string]any, error)
	GetNewTokenPair(context.Context, string) (map[string]any, error)
}

func New(u Usecase, cfg *config.Config) *Server {
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

// Decode refresh token from body.
// Decode token string from base64.
// Extract access token from header.
// Call usecase to refresh token pair.
// Encode new refresh token to base64.
func (s *Server) refreshToken(w http.ResponseWriter, r *http.Request) {
	slog.Info("refresh token called")

	// Decode refresh token from body.
	var token models.RefreshToken
	s.decodeBody(r, &token)

	// Decode token string from base64.
	tokStr := make([]byte, 1024)
	if _, err := base64.StdEncoding.Decode(tokStr, []byte(token.TokenString)); err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, s.encodeToJSON(Response{
			Error:   e.ErrBadRequest.Error(),
			Content: nil,
		}))
		return
	}

	tokenString := s.stripZeros(tokStr)

	token.TokenString = string(tokenString)
	fmt.Println(token.TokenString)

	// Extract access token from header.
	access, err := s.jwt.ExtractTokenString(r)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, s.encodeToJSON(Response{
			Error:   err.Error(),
			Content: nil,
		}))
		return
	}

	// Call usecase to refresh token pair.
	data, err := s.u.RefreshTokenPair(r.Context(), token, access)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, s.encodeToJSON(Response{
			Error:   err.Error(),
			Content: nil,
		}))
		return
	}

	// Encode new refresh token to base64.
	newRefresh := data["refresh_token"]
	newRefreshToken, ok := newRefresh.(models.RefreshToken)
	if !ok {
		slog.Error("conversion error")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, s.encodeToJSON(Response{
			Error:   e.ErrInternal.Error(),
			Content: nil,
		}))
		return
	}

	newRefreshToken.TokenString = base64.StdEncoding.EncodeToString([]byte(newRefreshToken.TokenString))
	data["refresh_token"] = newRefreshToken

	fmt.Fprint(w, s.encodeToJSON(Response{
		Error:   "",
		Content: data,
	}))
}

// Get guid from path value.
// Call usecase to generate pair.
// Encode new refresh token to base64.
func (s *Server) getTokenPair(w http.ResponseWriter, r *http.Request) {
	slog.Info("get token pair is called")

	// Get guid from path value.
	// Call usecase to generate pair.
	tokens, err := s.u.GetNewTokenPair(r.Context(), r.PathValue("id"))
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, s.encodeToJSON(Response{
			Error:   err.Error(),
			Content: nil,
		}))
		return
	}

	// Encode new refresh token to base64.
	newRefresh := tokens["refresh_token"]
	newRefreshToken, ok := newRefresh.(models.RefreshToken)
	if !ok {
		slog.Error("conversion error")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, s.encodeToJSON(Response{
			Error:   e.ErrInternal.Error(),
			Content: nil,
		}))
		return
	}

	fmt.Println("HERE", tokens["refresh_token"].(models.RefreshToken).TokenString)
	newRefreshToken.TokenString = base64.StdEncoding.EncodeToString([]byte(newRefreshToken.TokenString))
	tokens["refresh_token"] = newRefreshToken

	fmt.Fprint(w, s.encodeToJSON(Response{
		Error:   "",
		Content: tokens,
	}))
}

// Access token testing endpoint.
func (s *Server) restricted(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, s.encodeToJSON(Response{
		Error:   "",
		Content: "got",
	}))
}

func (s *Server) encodeToJSON(resp Response) string {
	encoded, err := json.Marshal(resp)
	if err != nil {
		slog.Error(err.Error())
	}

	return string(encoded)
}

func (s *Server) decodeBody(r *http.Request, dest *models.RefreshToken) {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		slog.Error(err.Error())
	}
}

func (s *Server) stripZeros(token []byte) []byte {
	result := []byte{}
	for _, b := range token {
		if b != 0 {
			result = append(result, b)
		}
	}
	return result
}
