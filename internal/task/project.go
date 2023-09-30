package task

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/blackflagsoftware/forge-go/internal/menu"
	m "github.com/blackflagsoftware/forge-go/internal/model"
)

func StartProject() error {
	project := m.Project{}
	if !loadProjectFile(&project.ProjectFile) {
		// .forge is not present, create new
		pwd, err := menu.SetupMenu()
		if err != nil {
			return err
		}
		if !initProjectFile(&project.ProjectFile, pwd) {
			return fmt.Errorf("Unable to create project file")
		}
		buildProjectName(&project)
		menu.StorageInitMenu(&project.ProjectFile)
		menu.TagFormatInitMenu(&project.ProjectFile)
		SaveProjectFile(project.ProjectFile)
	}

	menu.ProjectMenu(&project)
	return nil
}

func loadProjectFile(p *m.ProjectFile) bool {
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

func initProjectFile(p *m.ProjectFile, pwd string) bool {
	p.Message = "This is used for forge program for convenience"
	p.FullPath = pwd
	p.VersionPath = "v1"
	p.AppName = path.Base(pwd)
	p.Name.BuildName(p.AppName, []string{})
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

func buildProjectName(p *m.Project) {
	projectName := m.Name{}
	projectName.BuildName(p.ProjectFile.RawName, []string{})
	p.ProjectNameAbbr = projectName.Abbr
	p.ProjectNameAllLower = projectName.AllLower
	p.ProjectNameCamel = projectName.Camel
	p.ProjectNameLower = projectName.Lower
	p.ProjectNameLowerCamel = projectName.LowerCamel
	p.ProjectNameUpper = projectName.Upper
}

func SaveProjectFile(p m.ProjectFile) {
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
