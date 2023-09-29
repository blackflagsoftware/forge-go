package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	ae "github.com/blackflagsoftware/forge-go/base/internal/api_error"
)

// a simplified http client to call external http servers, it is designed to fail if the status code is not in the 200's
// inputs:
//		method: standard http method, i.e.: POST, GET, PATCH, etc
//		url: the full http url
//		bodyIn: map or struct, this function will json.Marshal it for the body payload, pass in 'nil' if not needed
//		bodyOut: map or *struct, this function will json.Unmarshal the response.Body to this variable
//		headers: map[string]string of headers to add to the request
func HTTPRequest(method, url string, bodyIn, bodyOut interface{}, header map[string]string) (err error) {
	var bodyInReader *bytes.Reader
	if bodyIn != nil {
		// take the data structure and create a body (io.Reader)
		bodyInBytes, err := json.Marshal(bodyIn)
		if err != nil {
			return ae.GeneralError("", err) // TODO: better error
		}
		bodyInReader = bytes.NewReader(bodyInBytes)
	}
	request, err := http.NewRequest(strings.ToUpper(method), url, bodyInReader)
	if err != nil {
		return ae.GeneralError("", err) // TODO: better error
	}
	// add headers
	for k, v := range header {
		request.Header.Add(k, v)
	}
	client := &http.Client{}
	resp := &http.Response{}
	resp, err = client.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return ae.GeneralError("", fmt.Errorf("Bad status code: %s", resp.Status))
	}
	body, errRead := io.ReadAll(resp.Body)
	if errRead != nil {
		return ae.GeneralError("", errRead)
	}
	// try to push response body into bodyOut
	if errJson := json.Unmarshal(body, bodyOut); errJson != nil {
		return ae.GeneralError("", errJson)
	}
	return
}
