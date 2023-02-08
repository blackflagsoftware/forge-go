package project

import (
	"fmt"
	"os"

	pf "github.com/blackflagsoftware/forge-go/internal/projectfile"
)

func StartForge() {
	if os.Getenv("FORGE_PATH") == "" {
		fmt.Println("FORGE_PATH is not set... Goodbye")
		return
	}
	project := Project{}
	if !project.ProjectFile.LoadProjectFile() {
		// .forge is not present, create new
		pwd, err := pf.SetupMenu()
		if err != nil {
			fmt.Println(err)
			return
		}
		if !project.ProjectFile.CreateProjectFile(pwd) {
			fmt.Println("Bye...")
			return
		}
		project.ProjectFile.StorageMenu()
		project.ProjectFile.TagFormatMenu()
	}
	defer project.ProjectFile.SaveProjectFile()
	// need to load the settings
	project.ProjectFile.LoadProjectFile()

	project.ProjectMenu()
	fmt.Println("Goodbye...")
}
