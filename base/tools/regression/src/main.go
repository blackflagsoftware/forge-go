package src

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	mod "github.com/blackflagsoftware/forge-go/base/tools/regression/model"
	util "github.com/blackflagsoftware/forge-go/base/tools/regression/util"
	"gopkg.in/yaml.v2"
)

type (
	Result struct {
		DateTime  time.Time  `json:"dateTime"`
		Total     int        `json:"total"`
		Successed int        `json:"succeeded"`
		Failure   int        `json:"failed"`
		Skipped   int        `json:"skipped"`
		Tests     []mod.Test `json:"tests"`
	}
)

func Process(args mod.Args) {
	content, err := os.ReadFile(args.FileName)
	if err != nil {
		fmt.Printf("Unable to open file: %s; %s", args.FileName, err)
		return
	}
	var tests []mod.Test
	// decide which unmarshal to use
	switch filepath.Ext(args.FileName) {
	case ".json":
		if err := json.Unmarshal(content, &tests); err != nil {
			fmt.Println("Error going from content to list of tests [json]:", err)
			return
		}
	case ".yml", ".yaml":
		if err := yaml.Unmarshal(content, &tests); err != nil {
			fmt.Println("Error going from content to list of tests [yaml]:", err)
			return
		}
	}
	RunTests(tests, args)

	fmt.Println("*** End regression testing ***")
}

func RunTests(tests []mod.Test, args mod.Args) {
	util.PrePopulateDynamicValues(args.DynamicValuesPath)
	fileName := filepath.Base(args.FileName)
	fmt.Printf("\n*** Start of Testing: %s ***\n\n", fileName)
	for i := range tests {
		// active and in the correct environment
		tests[i].Status = "SUCCEEDED"
		if !tests[i].Active {
			tests[i].Status = "SKIPPED"
			tests[i].Messages = append(tests[i].Messages, "test is not active, skipping")
			tests[i].PrintResults(fileName)
			continue
		}
		// if "run_environments" is empty, set at least one to "dev"
		if len(tests[i].RunEnvironments) == 0 {
			tests[i].RunEnvironments = append(tests[i].RunEnvironments, "dev")
		}
		// now check if the test's environment call be ran against the environment incoming environment
		found := false
		for _, e := range tests[i].RunEnvironments {
			if e == args.Environment {
				found = true
				break
			}
		}
		if !found {
			tests[i].Status = "SKIPPED"
			tests[i].Messages = append(tests[i].Messages, fmt.Sprintf("test's run_enviroment does not match current environment: %s, skipping", args.Environment))
			tests[i].PrintResults(fileName)
			continue
		}
		switch strings.ToLower(tests[i].TestType) {
		case "rest":
			tests[i].RunRest()
			tests[i].PrintResults(fileName)
		case "grpc":
			RunGrpc(&tests[i])
		default:
			fmt.Println("Invalid test type:", tests[i].TestType)
		}
	}
	util.SaveDynamicValuesToFile(args.DynamicValuesPath)
	result := printSummary(tests)
	saveResults(result)
	if args.PrintDebug {
		util.PrintDynamicValue()
	}
}

func RunGrpc(test *mod.Test) {

}

func printSummary(tests []mod.Test) Result {
	now := time.Now()
	result := Result{DateTime: now, Tests: tests}
	for _, test := range result.Tests {
		switch test.Status {
		case "SUCCEEDED":
			result.Successed++
		case "SKIPPED":
			result.Skipped++
		case "FAILED":
			result.Failure++
		}
	}
	result.Total = len(result.Tests)
	totalLine := fmt.Sprintf("|  Total: %d | Successed: %d | Failed: %d | Skipped: %d |", result.Total, result.Successed, result.Failure, result.Skipped)
	line := ""
	for i := 0; i < len(totalLine)-2; i++ {
		line += "-"
	}
	box := fmt.Sprintf("+%s+", line)
	fmt.Printf("\n%s\n", box)
	fmt.Println(totalLine)
	fmt.Printf("%s\n\n", box)
	return result
}

func saveResults(result Result) {
	err := os.MkdirAll("./regression_results", 0744)
	if err != nil {
		fmt.Printf("Unable to create results directory structure: %s\n", err)
		return
	}
	fileName := fmt.Sprintf("./regression_results/reg_%s.json", result.DateTime.Format("20060102150405"))
	fileContent, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		fmt.Printf("Unable to encode results: %s\n", err)
		return
	}
	os.WriteFile(fileName, fileContent, 0644)
	fmt.Printf("Saved file to: %s\n\n", fileName)
}
