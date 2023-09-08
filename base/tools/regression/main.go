package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	mod "github.com/blackflagsoftware/forge-go/base/tools/regression/model"
	"github.com/blackflagsoftware/forge-go/base/tools/regression/src"
)

func main() {
	var testPathWithFile string
	var testPath string
	var environment string
	var dynamicValuesPath string
	var printDebug bool
	flag.StringVar(&testPathWithFile, "testPathFile", "", "path/file of tests")
	flag.StringVar(&testPath, "testPath", "", "optional: path to files of test, if present then testPathWithFile is not needed")
	flag.StringVar(&environment, "environment", "", "optional: environment to run against the test's run_environment")
	flag.StringVar(&dynamicValuesPath, "dynamicValuesPath", "", "optional: load and save path/file for dynamic values")
	flag.BoolVar(&printDebug, "printDebug", false, "optional: print all debugging lines")
	flag.Parse()

	// set defaults if needed
	// environment: what comes over via flags, trumps; set to 'default' if nothing is set via flags or env var
	varEnvironment := os.Getenv("FORGE_GO_BASE_REGRESSION_ENVIRONMENT")
	if environment == "" && varEnvironment != "" {
		environment = varEnvironment
	}
	if environment == "" {
		environment = "dev"
	}

	if testPathWithFile == "" && testPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	// test files: see if a file or path is set
	testFiles := []string{}

	if testPathWithFile != "" {
		if _, err := os.Stat(*&testPathWithFile); os.IsNotExist(err) {
			fmt.Printf("File not found at: %s, nothing to do, ending!", testPathWithFile)
			os.Exit(1)
		}
		testFiles = append(testFiles, testPathWithFile)
	}
	if testPath != "" {
		if _, err := os.Stat(*&testPath); os.IsNotExist(err) {
			fmt.Printf("Path not found at: %s, nothing to do, ending!", testPath)
			os.Exit(1)
		}
		// read all files in the directory looking for .json or .yml or .yaml
		files, err := os.ReadDir(testPath)
		if err != nil {
			fmt.Printf("Error reading dirctory: %s with error: %s", testPath, err)
			os.Exit(1)
		}
		for _, file := range files {
			extName := filepath.Ext(file.Name())
			if !file.IsDir() && extName == ".json" || extName == ".yml" || extName == ".yaml" {
				testFiles = append(testFiles, fmt.Sprintf("%s/%s", filepath.Base(testPath), file.Name()))
			}
		}
		// since there might be dynamic values set from file to another file, make sure something is set if nothing is set with flags but
		// this is only done when multiple files are used
		if dynamicValuesPath == "" {
			dynamicValuesPath = "./dynamic_values"
		}
	}
	for _, testFile := range testFiles {
		args := mod.Args{FileName: testFile, Environment: environment, PrintDebug: printDebug, DynamicValuesPath: dynamicValuesPath}
		src.Process(args)
	}
}
