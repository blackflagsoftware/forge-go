package middleware

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

var logOutput io.Writer

func init() {
	// set the default log output file to be stdout
	// this can be overridden by calling SetLogOutput
	logOutput = os.Stdout
}

func SetLogOutput(out io.Writer) {
	logOutput = out
}

func Handler(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		// start time for duration
		start := time.Now()

		// handle request, i.e. go down the stack
		if err := h(c); err != nil {
			c.Error(err)
		}

		latency := fmt.Sprintf("%v", time.Since(start))

		// set up log format
		logFormat := &log.JSONFormatter{}
		logFormat.TimestampFormat = time.RFC3339Nano

		log.SetFormatter(logFormat)
		log.SetOutput(logOutput)

		url := c.Request().URL.String()

		logRequest := true
		if strings.Contains(url, "server_status") && c.Request().Method == "HEAD" {
			logRequest = false
		}

		if logRequest {
			log.WithFields(
				log.Fields{
					"method":      c.Request().Method,
					"status_code": c.Response().Status,
					"status_text": http.StatusText(c.Response().Status),
					"request_url": url,
					"latency":     latency,
					"referer":     c.Request().Referer(),
					"user_agent":  c.Request().UserAgent(),
					"remote":      c.Request().RemoteAddr,
				},
			).Infoln("completed")
		}

		return nil
	}
}

func DebugHandler(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		var requestDump []byte
		getBody := true
		if c.Request().Method == "POST" || c.Request().Method == "PUT" {
			contentType := c.Request().Header.Get("Content-Type")
			if strings.Contains(contentType, "multipart/form-data") {
				getBody = false
			}
			if requestDump, err = httputil.DumpRequest(c.Request(), getBody); err == nil {
				log.Println("/******** Request Parameters ********/")
				log.Printf("%s", string(requestDump))
				log.Printf("\n/******** End ********/")
			} else {
				log.Errorf("[DebugHandler]: %s", err)
			}
		}
		return h(c)
	}
}
