package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	t "github.com/blackflagsoftware/forge-go/internal/task"
	"github.com/blackflagsoftware/forge-go/internal/util"
)

func main() {
	flag.Parse()
	args := flag.Args()

	if os.Getenv("FORGE_PATH") == "" {
		fmt.Println("FORGE_PATH is not set, please set before going on... goodbye")
		return
	}
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		fmt.Println("GOPATH is not set, please set this env. var")
		return
	}
	if len(args) == 0 {
		// run forge
		if err := startForge(); err != nil {
			fmt.Println(err)
		}
		return
	}
	if args[0] == "clone" {
		// run clone
		directory := ""
		if len(args) > 1 {
			directory = args[1]
		}
		cloneMe(directory)
		return
	}
	if args[0] == "debug" {
		// this is used to debug a project that has been cloned
		// usage: forge debug <full path to project dir>
		if len(args) == 2 {
			os.Chdir(args[1])
			fmt.Printf("Debugging %s: attach debug process now, press 'enter' to continue", args[1])
			util.ParseInput()
			if err := startForge(); err != nil {
				fmt.Println(err)
			}
			return
		}
	}
	// usage
	fmt.Printf("Usage:\n\tforge (no args): run forge process, make new set of endpoints, etc\n\tforge clone: interactive clone option to create new project\n\tforge clone <directory>: clone and create new project at <directory> location, directory can be full path or path past GOPATH/src\n")
}

// clone the project
func cloneMe(directory string) {
	util.ClearScreen()
	err := os.Chdir(os.Getenv("FORGE_PATH") + "/tools/clone")
	if err != nil {
		fmt.Println(err)
	}
	output, _ := exec.Command("go", "run", "clone.go", "-projectPath", directory).CombinedOutput()
	fmt.Printf("%s\n", output)
}

// start the forge process in the project's directory
// check for dependency programs that forge uses
func startForge() error {
	processes := []string{"protoc", "protoc-gen-go", "protoc-gen-go-grpc", "mockgen"}

	missingProcess := []string{}
	for _, p := range processes {
		cmdProcess := fmt.Sprintf("which %s | wc -l", p)
		execProcess := exec.Command("bash", "-c", cmdProcess)
		out, err := execProcess.Output()
		if err != nil {
			return fmt.Errorf("Unable to check for dependent processes: %s", err)
		}
		out = bytes.TrimSpace(out)
		if bytes.Equal(out, []byte("0")) {
			missingProcess = append(missingProcess, p)
		}
	}
	if len(missingProcess) > 0 {
		return fmt.Errorf("Missing dependent process(es): %s; these will need to be installed.", strings.Join(missingProcess, ", "))
	}
	return t.StartProject()
}

// TODO: remove? - not used but has a test
func getProjectPath(goPath, directory string) string {
	basePath := goPath + "/src/"
	basePathLen := len(basePath)
	if strings.Contains(directory, basePath) {
		return string(directory[basePathLen:])
	}
	return directory
}
