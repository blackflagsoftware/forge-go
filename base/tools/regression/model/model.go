package model

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/Jeffail/gabs/v2"
	cli "github.com/blackflagsoftware/forge-go/base/tools/regression/client"
	util "github.com/blackflagsoftware/forge-go/base/tools/regression/util"
)

type (
	Args struct {
		FileName          string
		Environment       string
		PrintDebug        bool
		DynamicValuesPath string
	}

	Test struct {
		Name                 string            `json:"name" yaml:"name"`
		Status               string            `json:"status" yaml:"status"`
		Active               bool              `json:"active" yaml:"active"`
		Host                 string            `json:"host" yaml:"host"`
		Path                 string            `json:"path" yaml:"path"`
		TestType             string            `json:"test_type" yaml:"test_type"`               // rest or grpc
		RunEnvironments      []string          `json:"run_environments" yaml:"run_environments"` // array of environments this test can run against
		Method               string            `json:"method" yaml:"method"`                     // for rest: e.g. GET, POST, PUT, etc
		AuthUser             string            `json:"auth_user" yaml:"auth_user"`
		AuthPwd              string            `json:"auth_pwd" yaml:"auth_pwd"`
		RequestBody          interface{}       `json:"request_body" yaml:"request_body"`
		RequestHeaders       map[string]string `json:"request_headers" yaml:"request_headers"`
		RequestTimeout       int               `json:"request_timeout" yaml:"request_timeout"`
		ExpectedResponseBody interface{}       `json:"expected_response_body" yaml:"expected_response_body"`
		ExpectedStatus       int               `json:"expected_response_status" yaml:"expected_response_status"`
		ActualResponseBody   interface{}       `json:"actual_response_body" yaml:"actual_response_body"`
		ActualStatus         int               `json:"actual_status" yaml:"actual_status"`
		Messages             []string          `json:"messages" yaml:"messages"`
		WaitTime             int               `json:"wait_time" yaml:"wait_time"`
	}
)

func (t *Test) RunRest() {
	// if request timeout is 0 (empty) then set default to 30
	if t.RequestTimeout == 0 {
		t.RequestTimeout = 30
	}
	// check to replace all dynamic values
	util.DynamicInputString(&t.Path)
	// build url
	urlCall, err := url.PathUnescape(fmt.Sprintf("%s/%s", t.Host, t.Path))
	if err != nil {
		t.AppendMessage(fmt.Sprintf("unable to make url: %s", err))
		return
	}
	// transform request body
	// only make it a new reader if there is something there
	bodyByte, err := json.Marshal(t.RequestBody)
	if err != nil {
		t.AppendMessage("unable to make request body into bytes")
		return
	}
	util.DynamicInputByte(&bodyByte)
	var body io.Reader
	if len(bodyByte) > 0 {
		body = bytes.NewReader(bodyByte)
	}
	if err := cli.ValidateFormatMethod(&t.Method); err != nil {
		fmt.Println(err.Error())
		return
	}
	req, err := http.NewRequest(t.Method, urlCall, body)
	if err != nil {
		t.AppendMessage(fmt.Sprintf("unable to making request: %s", err))
	} else {
		ctx, cancel := context.WithTimeout(req.Context(), time.Duration(t.RequestTimeout*int(time.Second)))
		defer cancel()
		req = req.WithContext(ctx)
		// set up basic auth
		if t.AuthUser != "" && t.AuthPwd != "" {
			req.SetBasicAuth(t.AuthUser, t.AuthPwd)
		}
		for key, value := range t.RequestHeaders {
			req.Header.Add(key, value)
		}
		responseBody, responseStatus, err := cli.HTTPRequest(req)
		if err != nil {
			t.AppendMessage(fmt.Sprintf("http call failed - %s", err))
		}
		t.ActualStatus = responseStatus
		if responseStatus != t.ExpectedStatus {
			t.Status = "FAILED"
			t.AppendMessage(fmt.Sprintf("Status - want: %d; got: %d", t.ExpectedStatus, responseStatus))
			return
		}
		t.BodyCompare(responseBody)
	}
	if len(t.Messages) > 0 {
		t.Status = "FAILED"
	}
	if t.WaitTime > 0 {
		time.Sleep(time.Duration(time.Second * time.Duration(t.WaitTime)))
	}
}

func (t *Test) BodyCompare(responseBody []byte) {
	if (t.ExpectedResponseBody == nil || t.ExpectedResponseBody == "") && len(responseBody) != 0 {
		t.AppendMessage(fmt.Sprintf("expected request body is expecting null (or empty), actual response body has data: %s", responseBody))
		return
	}
	if t.ExpectedResponseBody == nil || t.ExpectedResponseBody == "" {
		return
	}
	responseContainer, err := gabs.ParseJSON(responseBody)
	if err != nil {
		t.AppendMessage("unable to covert response body to generic container")
		t.ActualResponseBody = string(responseBody)
		return
	}
	t.ActualResponseBody = responseContainer

	// since we are comparing what is in expected, we need to find the correct json base [] vs {}
	// the expected body will determine which values we will compare agains in the actual request body
	var expBodyByte []byte
	expectedBodyMap, ok := t.ExpectedResponseBody.(map[string]interface{})
	if ok {
		expBodyByte, err = json.Marshal(expectedBodyMap)
		if err != nil {
			t.AppendMessage("unable to marshal expected body, not a map")
			return
		}
	} else {
		expectedBodyArray := []interface{}{}
		switch reflect.TypeOf(t.ExpectedResponseBody).Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(t.ExpectedResponseBody)
			for i := 0; i < s.Len(); i++ {
				expectedBodyArray = append(expectedBodyArray, s.Index(i).Interface())
			}
		default:
			t.AppendMessage("not an array")
			return
		}
		expBodyByte, err = json.Marshal(expectedBodyArray)
		if err != nil {
			t.AppendMessage("unable to marshal expected body, not an array")
			return
		}
	}
	expectedContainer, err := gabs.ParseJSON(expBodyByte)
	if err != nil {
		t.AppendMessage("unable to covert expected body to generic container")
		return
	}
	// now that we have a container of containers thanks to gabs
	// determine if it is a map or an array
	expectedMap := expectedContainer.ChildrenMap()
	if len(expectedMap) != 0 {
		for key, value := range expectedMap {
			if !responseContainer.Exists(key) {
				t.AppendMessage(fmt.Sprintf("key: %s; not found", key))
			}
			path := []string{key}
			t.BodyCompareRecursive(path, value, responseContainer)
		}
	} else {
		expectedArray := expectedContainer.Children()
		if len(expectedArray) != 0 {
			for i, child := range expectedArray {
				path := []string{strconv.Itoa(i)}
				t.BodyCompareRecursive(path, child, responseContainer)
			}
		}
	}
}

func (t *Test) BodyCompareRecursive(path []string, expectedContainer, responseContainer *gabs.Container) {
	expectedMap := expectedContainer.ChildrenMap()
	if len(expectedMap) != 0 {
		for key, value := range expectedMap {
			path := append(path, key)
			t.BodyCompareRecursive(path, value, responseContainer)
		}
	} else {
		expectedArray := expectedContainer.Children()
		if len(expectedArray) != 0 {
			for i, child := range expectedArray {
				path := []string{strconv.Itoa(i)}
				t.BodyCompareRecursive(path, child, responseContainer)
			}
		} else {
			// not a map or array, a single element, let's check the value
			expectedElementBytes := expectedContainer.Bytes()
			responseElementBytes := responseContainer.Search(path...).Bytes()
			if util.IsDynamicInput(expectedElementBytes, responseElementBytes) {
				// the expected has a dynamic value, no need to compare
				return
			}
			if bytes.Compare(expectedElementBytes, responseElementBytes) != 0 {
				t.AppendMessage(fmt.Sprintf("mismatch [%s] => want: %s; got: %s", strings.Join(path, "/"), expectedElementBytes, responseElementBytes))
			}
		}
	}
}

func (t *Test) AppendMessage(msg string) {
	// yes, it is a simple one-liner but... much less typing when calling this over and over
	t.Messages = append(t.Messages, msg)
}

func (t Test) PrintResults(fileName string) {
	fmt.Printf("* Test: %s [%s]\n", t.Name, fileName)
	fmt.Printf("\tStatus: %s\n", t.Status)
	if t.Status == "FAILED" || t.Status == "SKIPPED" {
		fmt.Println("\tErrors:")
	}
	for _, errorMessage := range t.Messages {
		fmt.Printf("\t\t%s\n", errorMessage)
	}
}
