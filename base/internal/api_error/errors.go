package api_error

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/labstack/echo/v4"
)

type (
	ApiError struct {
		ApiErrorCode string `json:"api_error_code"`
		StatusCode   int    `json:"status_code"`
		Title        string `json:"title"`
		Detail       string `json:"detail"`
		Severe       bool   `json:"severe"`
		ProgramData  `json:"program_data"`
	}

	BodyError struct {
		StatusCode int    `json:"status_code"`
		Title      string `json:"title"`
		Detail     string `json:"detail"`
	}

	ProgramData struct {
		InternalError string `json:"internal_error"`
		FileName      string `json:"file_name"`
		FuncName      string `json:"func_name"`
		LineNumber    int    `json:"line_number"`
	}

	RequestData struct {
		User      string `json:"user"`
		IPAddr    string `json:"ip_addr"`
		UserAgent string `json:"user_agent"`
		Method    string `json:"method"`
		URI       string `json:"uri"`
		Body      string `json:"body"`
	}
)

// error handler for echo to handle
func ErrorHandler(err error, c echo.Context) {
	var apiErr ApiError
	switch err.(type) {
	case ApiError:
		apiErr = err.(ApiError)
	case *echo.HTTPError:
		httpErr := err.(*echo.HTTPError)
		if httpErr.Code == http.StatusNotFound || httpErr.Code == http.StatusUnauthorized {
			c.NoContent(httpErr.Code)
			return
		}
		if httpErr.Code == http.StatusMethodNotAllowed {
			apiErr = InvalidMethodError(c.Request().Method, c.Request().RequestURI, err)
		} else {
			apiErr = NewApiError(httpErr.Code, httpErr.Error(), c.Request().RequestURI, false, err)
		}
	default:
		apiErr = GeneralError("Internal Error", nil)
	}
	c.JSON(apiErr.StatusCode, apiErr.Error())
}

// base function call
func NewApiError(statusCode int, title string, detail string, severe bool, err error) ApiError {
	if err == nil {
		err = errors.New("")
	}
	if detail == "" {
		detail = err.Error() // use the err message
	}
	if err.Error() != "" {
		detail = fmt.Sprintf("%s - %s", detail, err.Error())
	}
	pg := SetCaller(err, 4) // by default, the level is 4 of which stack trace level we want to show
	return ApiError{
		ApiErrorCode: getApiErrorCode(),
		StatusCode:   statusCode,
		Title:        title,
		Detail:       detail,
		Severe:       severe,
		ProgramData:  pg,
	}
}

func (e ApiError) Error() string {
	bOut, err := json.Marshal(e)
	if err != nil {
		return BindError(err).Error()
	}
	return string(bOut)
}

func (e ApiError) BodyError() BodyError {
	return BodyError{StatusCode: e.StatusCode, Title: e.Title, Detail: e.Detail}
}

func SetCaller(err error, stackLevel int) ProgramData {
	pc := make([]uintptr, 1)
	runtime.Callers(stackLevel, pc)
	frames := runtime.CallersFrames(pc)
	f, _ := frames.Next()
	return ProgramData{InternalError: err.Error(), FileName: f.File, FuncName: f.Func.Name(), LineNumber: f.Line}
}

func getApiErrorCode() string {
	codeLen := 6
	const validChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	ll := len(validChars)
	b := make([]byte, codeLen)
	rand.Read(b)
	for i := 0; i < codeLen; i++ {
		b[i] = validChars[int(b[i])%ll]
	}
	return string(b)
}

// ***** Error methods *****
func GeneralError(detail string, err error) ApiError {
	return NewApiError(
		http.StatusInternalServerError,
		"Internal Server Error",
		detail,
		false,
		err,
	)
}

func DBError(detail string, err error) ApiError {
	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
		return NewApiError(
			http.StatusBadRequest,
			"No Results Error",
			"",
			false,
			err,
		)
	}
	return NewApiError(
		http.StatusBadRequest,
		"Database Query Error",
		detail,
		true,
		err,
	)
}

func DBEmptyRowError(err error) ApiError {
	return NewApiError(
		http.StatusBadRequest,
		"No Results Error",
		"",
		false,
		err,
	)
}

func BindError(err error) ApiError {
	return NewApiError(
		http.StatusBadRequest,
		"JSON Bind Error",
		"",
		false,
		err,
	)
}

func RouteNotFoundError(detail string, err error) ApiError {
	return NewApiError(
		http.StatusNotFound,
		"Route Not Found",
		detail,
		false,
		err,
	)
}

func ContentTypeError() ApiError {
	return NewApiError(
		http.StatusBadRequest,
		"Invalid Content-Type",
		"Required: application/json",
		false,
		nil,
	)
}

func ParseError(detail string) ApiError {
	return NewApiError(
		http.StatusBadRequest,
		"Parse Error",
		detail,
		false,
		nil,
	)
}

func MissingParamError(name string) ApiError {
	return NewApiError(
		http.StatusBadRequest,
		"Missing Parameter Error",
		fmt.Sprintf("%s is required", name),
		false,
		nil,
	)
}

func ParamError(name string, err error) ApiError {
	return NewApiError(
		http.StatusBadRequest,
		"Invalid Parameter",
		fmt.Sprintf("%s is invalid", name),
		false,
		err,
	)
}

func AuthorizationError(detail string) ApiError {
	return NewApiError(
		http.StatusUnauthorized,
		"Unauthorized",
		detail,
		false,
		nil,
	)
}

func LimiterError(err error) ApiError {
	return NewApiError(
		http.StatusTooManyRequests,
		"Too Many Request",
		"Received too many request in duration",
		false,
		err,
	)
}

func InvalidMethodError(method string, uri string, err error) ApiError {
	return NewApiError(
		http.StatusMethodNotAllowed,
		"Method Not Allowed",
		method+" "+uri,
		false,
		err,
	)
}

func StringLengthError(field string, colLen int) ApiError {
	return NewApiError(
		http.StatusBadRequest,
		"Max Length Exceeded",
		fmt.Sprintf("Field: %s - max length of %d", field, colLen),
		false,
		nil,
	)
}
