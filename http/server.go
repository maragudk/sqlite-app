package http

import (
	"context"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/maragudk/litefs-app/sql"
)

type Server struct {
	address  string
	database *sql.Database
	log      *log.Logger
	mux      chi.Router
	region   string
	server   *http.Server
}

type NewServerOptions struct {
	Database *sql.Database
	Host     string
	Log      *log.Logger
	Port     int
	Region   string
}

// NewServer returns an initialized, but unstarted Server.
// If no logger is provided, logs are discarded.
func NewServer(opts NewServerOptions) *Server {
	if opts.Log == nil {
		opts.Log = log.New(io.Discard, "", 0)
	}

	address := net.JoinHostPort(opts.Host, strconv.Itoa(opts.Port))
	mux := chi.NewMux()

	return &Server{
		address:  address,
		database: opts.Database,
		log:      opts.Log,
		mux:      mux,
		region:   opts.Region,
		server: &http.Server{
			Addr:              address,
			Handler:           mux,
			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      5 * time.Second,
			IdleTimeout:       5 * time.Second,
			ErrorLog:          opts.Log,
		},
	}
}

func (s *Server) Start() error {
	s.log.Println("Starting")

	s.setupRoutes()

	s.log.Println("Listening on http://" + s.address)
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) Stop() error {
	s.log.Println("Stopping")

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	return s.server.Shutdown(ctx)
}
