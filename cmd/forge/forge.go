package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	p "github.com/blackflagsoftware/forge-go/internal/project"
	"github.com/blackflagsoftware/forge-go/internal/util"
)

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		// run forge
		p.StartForge()
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
		if len(args) == 2 {
			util.ParseInput()
			os.Chdir(args[1])
			fmt.Println("I'm in", args[1])
			p.StartForge()
		}
	}

	// usage
	fmt.Printf("Usage:\n\tforge (no args): run forge process, make new set of endpoints, etc\n\tforge clone: interactive clone option to create new project\n\tforge clone <directory>: clone and create new project at <directory> location, directory can be full path or path past GOPATH/src\n")
}

func cloneMe(directory string) {
	util.ClearScreen()
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		fmt.Println("GOPATH is not set, please set this env. var")
		return
	}
	os.Chdir("tools/clone")
	output, _ := exec.Command("go", "run", "clone.go", "-projectPath", directory).CombinedOutput()
	fmt.Printf("%s\n", output)
}

func getProjectPath(goPath, directory string) string {
	basePath := goPath + "/src/"
	basePathLen := len(basePath)
	if strings.Contains(directory, basePath) {
		return string(directory[basePathLen:])
	}
	return directory
}
