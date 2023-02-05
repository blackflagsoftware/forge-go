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
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
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
		log.Infoln("Logging set to debug...")
		e.Debug = true
		e.Use(m.DebugHandler)
	}
	e.Use(
		middleware.Recover(),
		m.Handler,
	)

	// set output based on config value
	m.SetLogOutput(config.LogOutput)

	// set all non-endpoints here
	e.GET("/", Index)
	e.GET("/status", ServerStatus)

	InitializeRoutes(e)

	e.Start(fmt.Sprintf(":%s", restPort))
}

func setPidFile() {
	// purpose: to set the starting applications pid number to file
	if pidFile, err := os.Create(config.PidPath); err != nil {
		log.Panicln("Unable to create pid file...")
	} else if _, err := pidFile.Write([]byte(fmt.Sprintf("%d", os.Getpid()))); err != nil {
		log.Panicln("Unable to write pid to file...")
	}
}

func Index(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to the FORGE_GO_BASE API")
}

func ServerStatus(c echo.Context) error {
	return c.String(http.StatusOK, "live")
}

func InitializeRoutes(e *echo.Echo) {
	// initialize all routes here
	// --- replace server once text - do not remove ---
	// --- replace server text - do not remove ---

}
