package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	p "github.com/blackflagsoftware/forge-go/internal/project"
	"github.com/blackflagsoftware/forge-go/internal/util"
	"github.com/kardianos/osext"
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
	projectDirectory := getProjectPath(goPath, directory)

	execDir, _ := osext.ExecutableFolder()
	scriptDir := execDir + "/../../scripts"
	// TODO: this un-tars a project, would not work on Windows, can we not use tar but have the whole project and rename all the relevant imports?
	cmd := "./clone.sh " + projectDirectory

	os.Chdir(scriptDir)
	output, _ := exec.Command("sh", "-c", cmd).CombinedOutput()
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
