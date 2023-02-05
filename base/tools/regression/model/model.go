package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"time"

	"github.com/Jeffail/gabs/v2"
	cli "github.com/blackflagsoftware/forge-go/base/tools/regression/client"
	util "github.com/blackflagsoftware/forge-go/base/tools/regression/util"
)

type (
	Args struct {
		Content     []byte
		Environment string
	}

	Test struct {
		Name                 string            `json:"name"`
		Active               bool              `json:"active"`
		Host                 string            `json:"host"`
		Path                 string            `json:"path"`
		TestType             string            `json:"test_type"` // rest or grpc
		Method               string            `json:"method"`    // for rest: e.g. GET, POST, PUT, etc
		AuthUser             string            `json:"auth_user"`
		AuthPwd              string            `json:"auth_pwd"`
		RequestBody          interface{}       `json:"request_body"`
		RequestHeaders       map[string]string `json:"request_header"`
		ExpectedResponseBody interface{}       `json:"expected_response_body"`
		ExpectedStatus       int               `json:"expected_response_status"`
		ActualResponseBody   interface{}       `json:"acutal_response_body"`
		ActualStatus         string            `json:"acutal_status"`
		Messages             []string          `json:"messages"`
		Status               string            `json:"status"`
		WaitTime             int               `json:"wait_time"`
	}
)

func (t *Test) RunRest() {
	if !t.Active {
		t.Status = "SKIPPED"
		return
	}
	t.Status = "SUCCEEDED"
	// check to replace all dynamic values
	util.DynamicInputString(&t.Path)
	// build url
	urlCall := url.URL{}
	urlCall.Host = t.Host
	urlCall.Path = t.Path
	// transform request body
	// only make it a new reader if there is something there
	bodyByte, err := json.Marshal(t.RequestBody)
	if err != nil {
		t.AppendMessage("Error: unable to make request body into bytes")
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
	req, err := http.NewRequest(t.Method, urlCall.String(), body)
	if err != nil {
		t.AppendMessage(fmt.Sprintf("Error: unable to making request: %s", err))
	} else {
		responseBody, responseStatus, err := cli.HTTPRequest(req)
		if err != nil {
			t.AppendMessage(fmt.Sprintf("Error: http call failed - %s", err))
		}
		if responseStatus != t.ExpectedStatus {
			t.AppendMessage(fmt.Sprintf("Status => want: %d; got: %d", t.ExpectedStatus, responseStatus))
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
	responseContainer, err := gabs.ParseJSON(responseBody)
	if err != nil {
		t.AppendMessage("Error: unable to covert response body to generic container")
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
			t.AppendMessage("Error: unable to marshal expected body, not a map")
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
			t.AppendMessage("Error: not an array")
			return
		}
		expBodyByte, err = json.Marshal(expectedBodyArray)
		if err != nil {
			t.AppendMessage("Error: unable to marshal expected body, not an array")
			return
		}
	}
	expectedContainer, err := gabs.ParseJSON(expBodyByte)
	if err != nil {
		t.AppendMessage("Error: unable to covert expected body to generic container")
		return
	}
	// now that we have a container of containers thanks to gabs
	// determine if it is a map or an array
	expectedMap := expectedContainer.ChildrenMap()
	if len(expectedMap) != 0 {
		for key, value := range expectedMap {
			if !responseContainer.Exists(key) {
				t.AppendMessage("Error: key not found")
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
			path = append(path, key)
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
				t.AppendMessage(fmt.Sprintf("Mismatch => want: %s; got: %s", expectedElementBytes, responseElementBytes))
			}
		}
	}
}

func (t *Test) AppendMessage(msg string) {
	// yes, it is a simple one-liner but... much less typing when calling this over and over
	t.Messages = append(t.Messages, msg)
}
