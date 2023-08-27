package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/blackflagsoftware/forge-go/base/config"
	ae "github.com/blackflagsoftware/forge-go/base/internal/api_error"
	m "github.com/blackflagsoftware/forge-go/base/internal/middleware"
	p "github.com/labstack/echo-contrib/prometheus"
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
		m.Handler,
	)
	if config.EnableMetrics {
		prom := p.NewPrometheus("echo", nil)
		prom.Use(e)
	}

	// set all non-endpoints here
	e.GET("/", Index)
	e.HEAD("/status", ServerStatus) // for traditional server check
	e.GET("/liveness", Liveness)    // for k8s liveness

	InitializeRoutes(e)

	e.Start(fmt.Sprintf(":%s", restPort))
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
