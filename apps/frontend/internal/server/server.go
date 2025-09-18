package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Server wraps the Echo instance and startup/shutdown lifecycle.
type Server struct {
	app *echo.Echo
}

// New creates a configured Echo server with middleware and routes.
func New() *Server {
	e := echo.New()

	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.RequestID())
	e.Pre(middleware.RemoveTrailingSlash())

	registerStatic(e)
	registerRoutes(e)

	return &Server{app: e}
}

// Start runs the HTTP server and handles graceful shutdown.
func (s *Server) Start(ctx context.Context) error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)

	// Server start in a goroutine
	go func() {
		if err := s.app.Start(addr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("error starting server: %v", err)
		}
	}()

	// Wait for termination signal or context cancellation
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
	case <-stop:
	}

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return s.app.Shutdown(shutdownCtx)
}
