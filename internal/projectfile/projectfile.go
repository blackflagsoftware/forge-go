package projectfile

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
)

func (p *ProjectFile) CreateProjectFile(pwd string) bool {
	p.Message = "This is used for forge program for convenience"
	p.FullPath = pwd
	p.VersionPath = "v1"
	p.AppName = path.Base(pwd)
	name := p.Name.BuildName(p.AppName, p.KnownAliases)
	p.KnownAliases = append(p.KnownAliases, name)
	p.SubDir = "internal/" + p.VersionPath

	// create projectpath and subpackage
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		fmt.Println("***** GOPATH is not set, some of your paths will not be correct! *****")
	} else {
		idx := len(fmt.Sprintf("%s/src/", goPath))
		p.ProjectPath = p.FullPath[idx:]
	}
	subPath := strings.Split(p.SubDir, "/")
	p.SubPackage = subPath[len(subPath)-1]
	// encode paths
	p.ProjectPathEncoded = strings.Replace(p.ProjectPath, `/`, `\/`, -1)
	p.SubDirEncoded = strings.Replace(p.SubDir, `/`, `\/`, -1)
	return true
}

func (p *ProjectFile) LoadProjectFile() bool {
	if _, errStat := os.Stat("./.forge"); os.IsNotExist(errStat) {
		return false
	}
	bContent, errRead := os.ReadFile("./.forge")
	if errRead != nil {
		fmt.Printf("Error reading from .forge: %s", errRead)
		return true
	}
	errUnmarshal := json.Unmarshal(bContent, &p)
	if errUnmarshal != nil {
		fmt.Printf("Error extracting data from .forge: %s", errUnmarshal)
		return true
	}
	return true
}

func (p *ProjectFile) SaveProjectFile() {
	bContent, errMarshal := json.MarshalIndent(p, "", "    ")
	if errMarshal != nil {
		fmt.Println("Saving project file: unable to save -", errMarshal)
		return
	}
	errSave := os.WriteFile("./.forge", bContent, 0644)
	if errSave != nil {
		fmt.Println("Saving project file: unable to save -", errSave)
	}
}

func StorageTypeToProper(abbr string) string {
	switch abbr {
	case "s":
		return "SQL"
	case "f":
		return "File"
	case "m":
		return "MongoDB"
	}
	return abbr
}

func SqlTypeToProper(abbr string) string {
	switch abbr {
	case "p":
		return "Postgres"
	case "m":
		return "MySQL"
	case "s":
		return "Sqlite"
	}
	return abbr
}

func TagFormatToProper(tag string) string {
	switch tag {
	case "snakeCase":
		return "Snake Case"
	case "kebabCase":
		return "Kebab Case"
	case "camelCase":
		return "Camel Case"
	case "pascalCase":
		return "Pascal Case"
	case "upperCase":
		return "Upper Case"
	case "lowerCase":
		return "Lower Case"
	}
	return tag
}
