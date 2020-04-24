package server

import (
	"context"
	"go-blog/pkg/util/log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Config represents server specific config
type Config struct {
	ServiceName         string
	Port                string
	ReadTimeoutSeconds  int
	WriteTimeoutSeconds int
}

// New instantates new Echo server
func New() *echo.Echo {

	e := echo.New()

	e.Use(
		middleware.Logger(),    // default echo logger
		middleware.Recover(),   // recover from panics
		middleware.RequestID(), // generate ID for requests --> TODO: chequear skips, por ej /health
	)

	// default validator
	e.Validator = &CustomValidator{V: validator.New()}
	e.Binder = &CustomBinder{b: &echo.DefaultBinder{}}

	// health check
	e.GET("/health", healthCheckHandler)

	return e
}

// Start starts echo server
func Start(e *echo.Echo, cfg *Config, log *log.Log) {

	s := &http.Server{
		Addr:         cfg.Port,
		ReadTimeout:  time.Duration(cfg.ReadTimeoutSeconds) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeoutSeconds) * time.Second,
	}
	s.SetKeepAlivesEnabled(false)

	e.Debug = false
	e.HideBanner = true
	e.HidePort = true

	// start server
	log.Info("starting "+cfg.ServiceName, nil)
	go func() {
		if err := e.StartServer(s); err != nil {
			log.Error("error starting the server:", err, nil)
			return
		}
	}()

	// wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Error("error stopping server", err, nil)
	} else {
		log.Info(cfg.ServiceName+" stoped!", nil)
	}
}
