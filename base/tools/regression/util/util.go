package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

var DynamicValues map[string]Dynamic

const DynamicConst = "dyn:"

type (
	Dynamic struct {
		Value    string
		CastType string // can be: string | int | float | boolean (string is infered)
	}
)

func init() {
	DynamicValues = make(map[string]Dynamic)
}

func DynamicInputString(input *string) {
	// check if there is a DynamicConst contained in a string
	// this only called for the url, whether it is a string or possibly an int, we will treat it like a string
	// this will replace all occurences
	if strings.Contains(*input, DynamicConst) {
		for key, value := range DynamicValues {
			search := DynamicConst + key
			*input = strings.ReplaceAll(*input, search, value.Value)
		}
	}
}

func DynamicInputByte(input *[]byte) {
	// check if there is a DynamicConst contained in a byte
	// this is called on the body, we might want to substitute for string or one of the other types (int | float | bool)
	// since this will go into a json payload the type "format" is important
	// this will replace all occurences
	if bytes.Contains(*input, []byte(DynamicConst)) {
		for key, value := range DynamicValues {
			search := []byte(DynamicConst + key)
			if value.CastType != "" {
				castType := []byte(value.CastType)
				search = append(search, byte(':'))
				search = append(search, castType...)
			}
			if !(value.CastType == "string" || value.CastType == "") {
				search = []byte(fmt.Sprintf("\"%s\"", search))
			}
			*input = bytes.ReplaceAll(*input, search, []byte(value.Value))
		}
	}
}

func IsDynamicInput(expectedByte, responseByte []byte) bool {
	// checks to see if the DynamicConst is contained in the single element
	// if captures the DynamicValue even if it has seen it before
	if bytes.Contains(expectedByte, []byte(DynamicConst)) {
		parts := bytes.Split(expectedByte, []byte(":")) // => [0] = dyn; [1] = <value>
		dynKey := string(bytes.Trim(parts[1], "\""))
		dynValue := string(bytes.Trim(responseByte, "\""))
		dynCastType := "string"
		if len(parts) == 3 {
			dynCastType = string(bytes.Trim(parts[2], "\""))
		}
		DynamicValues[dynKey] = Dynamic{Value: dynValue, CastType: dynCastType}
		return true
	}
	return false
}

func PrePopulateDynamicValues(dynamicValuePath string) {
	now := time.Now().UTC()
	DynamicValues["DynDateMinusSeven"] = Dynamic{Value: now.AddDate(0, 0, -7).Format(time.RFC3339)}
	DynamicValues["DynDateMinusSix"] = Dynamic{Value: now.AddDate(0, 0, -6).Format(time.RFC3339)}
	DynamicValues["DynDateMinusFive"] = Dynamic{Value: now.AddDate(0, 0, -5).Format(time.RFC3339)}
	DynamicValues["DynDateMinusFour"] = Dynamic{Value: now.AddDate(0, 0, -4).Format(time.RFC3339)}
	DynamicValues["DynDateMinusThree"] = Dynamic{Value: now.AddDate(0, 0, -3).Format(time.RFC3339)}
	DynamicValues["DynDateMinusTwo"] = Dynamic{Value: now.AddDate(0, 0, -2).Format(time.RFC3339)}
	DynamicValues["DynDateMinusOne"] = Dynamic{Value: now.AddDate(0, 0, -1).Format(time.RFC3339)}
	DynamicValues["DynDateNow"] = Dynamic{Value: now.Format(time.RFC3339)}
	DynamicValues["DynDateAddOne"] = Dynamic{Value: now.AddDate(0, 0, 1).Format(time.RFC3339)}
	DynamicValues["DynDateAddTwo"] = Dynamic{Value: now.AddDate(0, 0, 2).Format(time.RFC3339)}
	DynamicValues["DynDateAddThree"] = Dynamic{Value: now.AddDate(0, 0, 3).Format(time.RFC3339)}
	DynamicValues["DynDateAddFour"] = Dynamic{Value: now.AddDate(0, 0, 4).Format(time.RFC3339)}
	DynamicValues["DynDateAddFive"] = Dynamic{Value: now.AddDate(0, 0, 5).Format(time.RFC3339)}
	DynamicValues["DynDateAddSix"] = Dynamic{Value: now.AddDate(0, 0, 6).Format(time.RFC3339)}
	DynamicValues["DynDateAddSeven"] = Dynamic{Value: now.AddDate(0, 0, 7).Format(time.RFC3339)}

	if dynamicValuePath != "" {
		if _, err := os.Stat(dynamicValuePath); !os.IsNotExist(err) {
			// found the file load the file
			content, errRead := os.ReadFile(dynamicValuePath)
			if errRead != nil {
				fmt.Println("Unable to read the values file:", errRead)
				return
			}
			value := map[string]Dynamic{}
			if errJson := json.Unmarshal(content, &value); errJson != nil {
				fmt.Println("Unable to unmarshal values file:", errJson)
				return
			}
			for k, v := range value {
				DynamicValues[k] = v
			}
		}
	}
}

func SaveDynamicValuesToFile(dynamicValuesPath string) {
	if dynamicValuesPath != "" {
		values := map[string]Dynamic{}
		// only get non-date
		for k, v := range DynamicValues {
			if strings.Index(k, "DynDate") != 0 {
				values[k] = v
			}
		}
		content, errJson := json.Marshal(values)
		if errJson != nil {
			fmt.Println("Unable to marshal values:", errJson)
			return
		}
		if err := os.WriteFile(dynamicValuesPath, content, 0644); err != nil {
			fmt.Println("Unable to write values file:", err)
		}
	}
}

// TODO: might not need
func RemoveDynamicValuesFile(dynamicValuesPath string) {
	if dynamicValuesPath != "" {
		os.Remove(dynamicValuesPath)
	}
}

func PrintDynamicValue() {
	fmt.Println("")
	fmt.Println("Dynamic Values:")
	for k, v := range DynamicValues {
		if strings.Index(k, "DynDate") != 0 {
			fmt.Printf("%s: %v\n", k, v)
		}
	}
	fmt.Println("")
}
