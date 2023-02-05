package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/blackflagsoftware/forge-go/internal/util"
)

/*
To clone:
- ask for new directory
- make new directory (if not exists)
- move the base files
- replace blackflagsotware/forge-go/base -> new folder structure
- replace name with new app name
- copy go.sum/.mod do a 'go mod tidy'
*/

func main() {
	var projectPath string
	flag.StringVar(&projectPath, "projectPath", "", "destination clone project path (e.g. github.com/<name>/<new project name>")
	flag.Parse()
	if projectPath == "" {
		fmt.Println("** Forge Clone **")
		fmt.Println("")
		fmt.Print("Where would you like to clone this (e.g. github.com/<name>/<new project_name> or (e) to exit? ")
		projectPath := util.ParseInput()
		if strings.ToLower(projectPath) == "e" {
			fmt.Println("Good bye")
			return
		}
	}
	sourceDir := filepath.Clean("../../base")
	goPath := os.Getenv("GOPATH") + "/src"
	destinationDir := filepath.Join(goPath, projectPath)

	if _, errExists := os.Stat(destinationDir); !os.IsNotExist(errExists) {
		fmt.Println("Project directory already exists... Good bye")
		return
	}
	if err := util.CopyDirectory(sourceDir, destinationDir); err != nil {
		fmt.Println("Good bye")
		return
	}
	split := strings.Split(projectPath, "/")
	projectName := split[len(split)-1]
	fmt.Printf("Cloning for %s is complete!", projectName)
	fmt.Println("Good bye")
}
