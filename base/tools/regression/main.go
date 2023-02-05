package main

import (
	"flag"
	"fmt"
	"os"

	mod "github.com/blackflagsoftware/forge-go/base/tools/regression/model"
	reg "github.com/blackflagsoftware/forge-go/base/tools/regression/regression"
)

// gather args and run the tests
func main() {
	file := flag.String("testFile", "", "path the the test file")
	if file == nil || *file == "" {
		flag.Usage()
		os.Exit(1)
	}
	if _, err := os.Stat(*file); os.IsNotExist(err) {
		fmt.Printf("%s does not exists", *file)
		os.Exit(1)
	}
	content, err := os.ReadFile(*file)
	if err != nil {
		fmt.Printf("Unable to open file: %s; %s", *file, err)
		os.Exit(1)
	}
	args := mod.Args{Content: content}
	reg.Process(args)
}
