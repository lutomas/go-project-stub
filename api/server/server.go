package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/lutomas/go-project-stub/pkg/cors"
	"github.com/urfave/negroni"
	"go.uber.org/zap"
)

type Options struct {
	// Public configurations
	HttpHost   string
	HttpPort   int
	EnableCors bool
	Logger     *zap.Logger
}

type Server struct {
	// Private configurations
	httpServer *http.Server
	httpAddr   string
	router     *mux.Router
	started    bool
	logger     *zap.Logger
	enableCors bool
}

func New(opts *Options) (*Server, error) {
	return &Server{
		httpAddr:   fmt.Sprintf("%s:%d", opts.HttpHost, opts.HttpPort),
		logger:     opts.Logger,
		enableCors: opts.EnableCors,
	}, nil
}

func (s *Server) ServeHTTP() error {

	// configuring middleware, router
	s.configureServer(s.httpAddr)

	errCh := make(chan error)

	go func() {
		// log.Info("starting server")
		s.started = true
		s.logger.Info("starting HTTP server...", zap.String("address", s.httpAddr))
		defer s.logger.Info("HTTP server stopped")
		err := s.httpServer.ListenAndServe()
		s.started = false
		if err != nil {
			if strings.Contains(err.Error(), "closed") {
				errCh <- nil
				return
			}
		}
		errCh <- err
	}()

	err := <-errCh

	return err
}

func (s *Server) configureServer(httpAddr string) {
	n := negroni.New()

	recoveryMiddleware := negroni.NewRecovery()
	n.Use(recoveryMiddleware)

	if s.enableCors {
		n.Use(negroni.HandlerFunc(cors.CorsHeadersMiddleware))
		s.logger.Warn("CORS middleware enabled! Unset environment variable `MAIN_APP_ENABLE_CORS` to disable it!")
	}

	// authentication middleware
	//n.Use(negroni.HandlerFunc(s.authenticationMiddleware))

	s.router = mux.NewRouter()
	subRouter := s.router.PathPrefix("/api").Subrouter()

	s.infoHandlers(subRouter)

	n.UseHandler(s.router)

	// HTTP server
	s.httpServer = &http.Server{
		Addr:              httpAddr,
		Handler:           n,
		IdleTimeout:       time.Second * 30,
		ReadTimeout:       time.Second * 60,
		ReadHeaderTimeout: time.Second * 60,
		WriteTimeout:      time.Second * 60,
	}
}

func (s *Server) Stop() error {
	if s.httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		err := s.httpServer.Shutdown(ctx)
		if err != nil {
			s.logger.Error("failed to stop HTTP server", zap.Error(err))
			return err
		}
	}

	return nil
}
