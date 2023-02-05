package client

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

var ValidMethods = map[string]struct{}{"GET": {}, "POST": {}, "PUT": {}, "PATCH": {}, "DELETE": {}}

func HTTPRequest(request *http.Request) (body []byte, statusCode int, err error) {
	client := &http.Client{}
	resp := &http.Response{}
	resp, err = client.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	statusCode = resp.StatusCode
	body, err = io.ReadAll(resp.Body)
	return
}

func ValidateFormatMethod(method *string) error {
	*method = strings.ToUpper(*method)
	if _, ok := ValidMethods[*method]; !ok {
		return fmt.Errorf("Method not valid: %s", *method)
	}
	return nil
}
