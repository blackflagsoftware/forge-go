package project

// import (
// 	"bytes"
// 	"fmt"
// 	"os"
// 	"os/exec"
// 	"strings"

// 	"github.com/blackflagsoftware/forge-go/internal/menu"
// 	m "github.com/blackflagsoftware/forge-go/internal/model"
// )

// func StartForge() {
// 	if os.Getenv("FORGE_PATH") == "" {
// 		fmt.Println("FORGE_PATH is not set... Goodbye")
// 		return
// 	}
// 	if err := checkDependencies(); err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	project := m.Project{}
// 	if !project.ProjectFile.LoadProjectFile() {
// 		// .forge is not present, create new
// 		pwd, err := menu.SetupMenu()
// 		if err != nil {
// 			fmt.Println(err)
// 			return
// 		}
// 		if !project.ProjectFile.CreateProjectFile(pwd) {
// 			fmt.Println("Bye...")
// 			return
// 		}
// 		menu.StorageInitMenu(&project.ProjectFile)
// 		menu.TagFormatInitMenu(&project.ProjectFile)
// 	}
// 	defer project.ProjectFile.SaveProjectFile()
// 	// need to load the settings
// 	project.ProjectFile.LoadProjectFile()

// 	menu.ProjectMenu(&project)
// 	fmt.Println("Goodbye...")
// }

// func checkDependencies() error {
// 	processes := []string{"protoc", "protoc-gen-go", "protoc-gen-go-grpc", "mockgen"}

// 	missingProcess := []string{}
// 	for _, p := range processes {
// 		cmdProcess := fmt.Sprintf("which %s | wc -l", p)
// 		execProcess := exec.Command("bash", "-c", cmdProcess)
// 		out, err := execProcess.Output()
// 		if err != nil {
// 			return fmt.Errorf("Unable to check for dependent processes: %s", err)
// 		}
// 		out = bytes.TrimSpace(out)
// 		if bytes.Equal(out, []byte("0")) {
// 			missingProcess = append(missingProcess, p)
// 		}
// 	}
// 	if len(missingProcess) > 0 {
// 		return fmt.Errorf("Missing dependent process(es): %s; these will need to be installed.", strings.Join(missingProcess, ", "))
// 	}
// 	return nil
// }
