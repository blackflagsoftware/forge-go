package regression

import (
	"encoding/json"
	"fmt"
	"strings"

	mod "github.com/blackflagsoftware/forge-go/base/tools/regression/model"
)

func Process(args mod.Args) {
	var tests []mod.Test
	if err := json.Unmarshal(args.Content, &tests); err != nil {
		fmt.Println("Error going from content to list of tests:", err)
		return
	}
	// process test
	RunTests(tests)

	// print out test results
	fmt.Println("done")
}

func RunTests(tests []mod.Test) {
	for i := range tests {
		switch strings.ToLower(tests[i].TestType) {
		case "rest":
			tests[i].RunRest()
		case "grpc":
			RunGrpc(&tests[i])
		default:
			fmt.Println("Invalid test type:", tests[i].TestType)
		}
	}
}

func RunGrpc(test *mod.Test) {

}
