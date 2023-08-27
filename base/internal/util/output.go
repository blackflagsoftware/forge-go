package util

import (
	"net/http"
	"strconv"

	ae "github.com/blackflagsoftware/forge-go/base/internal/api_error"
	m "github.com/blackflagsoftware/forge-go/base/internal/middleware"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type (
	Output struct {
		Payload interface{} `json:"data,omitempty"`
		*Error  `json:"error,omitempty"`
		*Meta   `json:"meta,omitempty"`
	}

	Error struct {
		Id     string `json:"Id,omitempty"`
		Title  string `json:"Title,omitempty"`
		Detail string `json:"Detail,omitempty"`
		Status string `json:"Status,omitempty"`
	}

	Meta struct {
		TotalCount int `json:"total_count"`
	}
)

func NewOutput(c echo.Context, payload interface{}, apiError *ae.ApiError, totalCount *int) Output {
	var err *Error
	var meta *Meta
	if apiError != nil {
		LogError(c, apiError)
		err = &Error{Id: apiError.ApiErrorCode, Title: apiError.Title, Detail: apiError.Detail, Status: strconv.Itoa(apiError.StatusCode)}
	}
	if totalCount != nil {
		meta = &Meta{TotalCount: *totalCount}
	}
	output := Output{
		Payload: payload,
		Error:   err,
		Meta:    meta,
	}
	return output
}

func LogError(c echo.Context, apiError *ae.ApiError) {
	url := c.Request().URL.String()

	m.Default.WithFields(
		logrus.Fields{
			"method":      c.Request().Method,
			"status_code": apiError.StatusCode,
			"status_text": http.StatusText(c.Response().Status),
			"request_url": url,
			"referer":     c.Request().Referer(),
			"user_agent":  c.Request().UserAgent(),
			"remote":      c.Request().RemoteAddr,
			"detail":      apiError.Detail,
		},
	).Errorln("error")
}
