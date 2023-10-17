package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/blackflagsoftware/forge-go/base/config"
	ae "github.com/blackflagsoftware/forge-go/base/internal/api_error"
	m "github.com/blackflagsoftware/forge-go/base/internal/middleware"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	// --- replace migration header once text - do not remove ---
	// --- replace main header text - do not remove ---
)

func main() {
	setPidFile()

	// argument flag
	var restPort string
	flag.StringVar(&restPort, "restPort", "", "the port number used for the REST listener")

	flag.Parse()

	if restPort == "" {
		restPort = config.RestPort
	}
	// --- replace migration once text - do not remove ---

	e := echo.New()
	e.HTTPErrorHandler = ae.ErrorHandler // set echo's error handler
	if !strings.Contains(config.Env, "prod") {
		m.Default.Infoln("Logging set to debug...")
		e.Debug = true
		e.Use(m.DebugHandler)
	}
	e.Use(
		middleware.Recover(),
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
		}), // remove if you don't need
		m.Handler,
	)
	if config.EnableMetrics {
		e.Use(echoprometheus.NewMiddleware("FORGE_GO_BASE"))
		e.GET("/metrics", echoprometheus.NewHandler())
	}

	// set all non-endpoints here
	e.GET("/", Index)
	e.HEAD("/status", ServerStatus) // for traditional server check
	e.GET("/liveness", Liveness)    // for k8s liveness

	InitializeRoutes(e)

	go func() {
		if err := e.Start(fmt.Sprintf(":%s", restPort)); err != nil && err != http.ErrServerClosed {
			m.Default.Printf("graceful server stop with error: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		m.Default.Printf("gracefult shutdown with error: %s", err)
	}
}

func setPidFile() {
	// purpose: to set the starting applications pid number to file
	if pidFile, err := os.Create(config.PidPath); err != nil {
		m.Default.Panicln("Unable to create pid file...")
	} else if _, err := pidFile.Write([]byte(fmt.Sprintf("%d", os.Getpid()))); err != nil {
		m.Default.Panicln("Unable to write pid to file...")
	}
}

func Index(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to the FORGE_GO_BASE API")
}

func ServerStatus(c echo.Context) error {
	c.Response().Header().Add("FORGE_GO_BASE", config.AppVersion)
	c.Response().WriteHeader(http.StatusOK)
	return nil
}

func Liveness(c echo.Context) error {
	return c.String(http.StatusOK, "live")
}

func InitializeRoutes(e *echo.Echo) {
	// initialize all routes here
	// --- replace server once text - do not remove ---
	// --- replace server text - do not remove ---

}
