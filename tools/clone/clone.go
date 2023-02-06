package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
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
	// determine if project directory already exists, copy directory if not
	if _, errExists := os.Stat(destinationDir); !os.IsNotExist(errExists) {
		fmt.Println("Project directory already exists... Good bye")
		return
	}
	if err := util.CopyDirectory(sourceDir, destinationDir); err != nil {
		fmt.Println("Good bye")
		return
	}
	// replace import items
	if err := SearchAndReplace(destinationDir, "github.com/blackflagsoftware/forge-go/base", projectPath); err != nil {
		fmt.Println("Good bye")
		return
	}
	split := strings.Split(projectPath, "/")
	projectName := split[len(split)-1]
	// replace normal syntax name
	if err := SearchAndReplace(destinationDir, "forge-go-base", projectName); err != nil {
		fmt.Println("Good bye")
		return
	}
	upperName := strings.ReplaceAll(strings.ToUpper(projectName), "-", "_")
	// replace upper syntax name
	if err := SearchAndReplace(destinationDir, "FORGE_GO_BASE", upperName); err != nil {
		fmt.Println("Good bye")
		return
	}
	// get mod files and copy them... do go mod tidy
	modSrc := sourceDir + "/../go.mod"
	modDest := destinationDir + "/go.mod"
	util.CopyFile(modSrc, modDest)
	sumSrc := sourceDir + "/../go.sum"
	sumDest := destinationDir + "/go.sum"
	util.CopyFile(sumSrc, sumDest)
	if err := os.Chdir(destinationDir); err != nil {
		fmt.Println("Unable to change to new folder:", err)
		fmt.Println("Good bye")
		return
	}
	fileContent, err := os.ReadFile(modDest)
	if err != nil {
		fmt.Printf("Unable to read file content %s: %s\n", modDest, err)
		fmt.Println("Good bye")
	}
	fileContent = bytes.ReplaceAll(fileContent, []byte("github.com/blackflagsoftware/forge-go"), []byte(projectPath))
	err = os.WriteFile(modDest, fileContent, 0644)
	if err != nil {
		fmt.Printf("Unable to write file content %s: %s\n", modDest, err)
	}
	cmd := exec.Command("go", "mod", "tidy")
	if _, err := cmd.Output(); err != nil {
		fmt.Println("'go mod tidy' had errors:", err)
	}

	fmt.Printf("Cloning for '%s' is complete!\n", projectName)
	fmt.Println("Good bye")
}

func SearchAndReplace(src, find, replace string) error {
	dirEntry, err := os.ReadDir(src)
	if err != nil {
		fmt.Printf("Unable to read directory %s: %s\n", src, err)
		return err
	}
	for _, entry := range dirEntry {
		newSrc := filepath.Join(src, entry.Name())
		if entry.IsDir() {
			if err := SearchAndReplace(newSrc, find, replace); err != nil {
				return err
			}
			continue
		}
		if filepath.Ext(newSrc) == ".go" {
			fileContent, err := os.ReadFile(newSrc)
			if err != nil {
				fmt.Printf("Unable to read file content %s: %s", newSrc, err)
				continue
			}
			fileContent = bytes.ReplaceAll(fileContent, []byte(find), []byte(replace))
			err = os.WriteFile(newSrc, fileContent, 0644)
			if err != nil {
				fmt.Printf("Unable to write file content %s: %s", newSrc, err)
			}
		}
	}
	return nil
}
