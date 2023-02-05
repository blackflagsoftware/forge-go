package util

import (
	"bytes"
	"strings"
)

var DynamicValues map[string]string

const DynamicConst = "dyn:"

func DynamicInputString(input *string) {
	// check for if the DynamicConst is contained in the string
	if strings.Contains(*input, DynamicConst) {
		for key, value := range DynamicValues {
			search := DynamicConst + key
			*input = strings.ReplaceAll(*input, search, value)
		}
	}
}

func DynamicInputByte(input *[]byte) {
	// check for if the DynamicConst is contained in the string
	if bytes.Contains(*input, []byte(DynamicConst)) {
		for key, value := range DynamicValues {
			search := []byte(DynamicConst + key)
			*input = bytes.ReplaceAll(*input, search, []byte(value))
		}
	}
}

func IsDynamicInput(expectedByte, responseByte []byte) bool {
	if bytes.Contains(expectedByte, []byte(DynamicConst)) {
		parts := bytes.Split(responseByte, []byte(":")) // => [0] = dyn; [1] = <value>
		dynKey := string(bytes.Trim(parts[1], "\""))
		dynValue := string(bytes.Trim(responseByte, "\""))
		DynamicValues[dynKey] = dynValue
		return true
	}
	return false
}
