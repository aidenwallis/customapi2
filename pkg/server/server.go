package server

import (
	"context"
	"net/http"
	"time"

	v2 "github.com/aidenwallis/customapi2/api/v2"
	"github.com/aidenwallis/customapi2/pkg/version"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Server struct {
	srv *http.Server
}

func New(addr string, cfg *version.Config) *Server {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(time.Second * 15))

	r.Mount("/api/v2", v2.New(cfg))

	return &Server{
		srv: &http.Server{
			Addr:    addr,
			Handler: r,
		},
	}
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
