package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/imjenal/transaction-service/api"
)

type Config struct {
	Port int
	Host string
}

// Server houses the http.Server and other variables for our HTTP server
type Server struct {
	server    http.Server
	router    *mux.Router
	connClose chan struct{} //connClose channel is closed when the http.Server is shutdown. It can be used to listen when the server closes
	apiParams *api.Params
	cfg       *Config
}

// New creates a new instance of Server
func New(config *Config, params *api.Params) *Server {
	r := mux.NewRouter().StrictSlash(true)

	return &Server{
		router:    r,
		connClose: make(chan struct{}, 1),
		server: http.Server{
			Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
			ReadTimeout:  2 * time.Second,
			WriteTimeout: 20 * time.Second,
			IdleTimeout:  65 * time.Second,
		},
		apiParams: params,
		cfg:       config,
	}
}

func (s *Server) Listen() {
	s.setup()
	log.Printf("Starting server... %s \n\n\n", s.server.Addr)
	if err := s.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Printf("Failed to start HTTP server: %v", err)
	}
}

// WaitForShutdown blocks until the server is shutdown
func (s *Server) WaitForShutdown() {
	// Block until the server has completed shutdown i.e., the connClose channel is closed
	// The channel is closed when the graceFullShutdown is triggered after receiving a signal from the OS
	// Check the graceFullShutdown function for more details
	<-s.connClose
}

func (s *Server) setup() {
	defer s.graceFullShutdown()
	s.routes()
	s.server.Handler = s.router
}

func (s *Server) graceFullShutdown() {
	go func() {
		// Create a channel to listen for OS signals
		sigint := make(chan os.Signal, 1)

		// Listen for SIGINT, SIGABRT, SIGTERM signals
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGABRT, syscall.SIGTERM)

		sig := <-sigint // Block until a signal is received

		// Once a signal is received, log it and shutdown the server

		log.Printf("OS terminate signal received : %v \n", sig)
		log.Printf("Shutting down server\n")

		// Create a context with timeout of 5 seconds
		// The shutdown will wait for 5 seconds before timing out and abruptly closing the connections
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.server.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down server: %v", err)
		}

		close(s.connClose) // Close the channel to notify the shutdown
	}()
}

func (s *Server) routes() {
	api.Routes(s.router.PathPrefix("/api/").Subrouter(), s.apiParams)
}
