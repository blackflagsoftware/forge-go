package middleware

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func Handler(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		// start time for duration
		start := time.Now()

		// handle request, i.e. go down the stack
		if err := h(c); err != nil {
			c.Error(err)
		}

		latencyStr := fmt.Sprintf("%v", time.Since(start))

		url := c.Request().URL.String()

		logRequest := true
		if (strings.Contains(url, "/status") && c.Request().Method == "HEAD") || strings.Contains(url, "/metrics") {
			logRequest = false
		}

		if logRequest {
			Default.WithFields(
				logrus.Fields{
					"method":      c.Request().Method,
					"status_code": c.Response().Status,
					"status_text": http.StatusText(c.Response().Status),
					"request_url": url,
					"latency":     latencyStr,
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
				Default.Println("/******** Request Parameters ********/")
				Default.Printf("%s", string(requestDump))
				Default.Printf("\n/******** End ********/")
			} else {
				Default.Errorf("[DebugHandler]: %s", err)
			}
		}
		return h(c)
	}
}
