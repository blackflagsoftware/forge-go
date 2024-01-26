package task

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"

	c "github.com/blackflagsoftware/forge-go/internal/constant"
	m "github.com/blackflagsoftware/forge-go/internal/model"
)

var (
	templatePath = fmt.Sprintf("%s/templates", os.Getenv("FORGE_PATH"))
	tmplFiles    = []string{"model", "rest", "manager", "grpc"} // TODO: sql is optional depending on which data storages they want to template, dynamically build this
)

func StartTemplating(project *m.Project) {
	BuildStorage(project)

	for i := range project.Entities {
		savePath := fmt.Sprintf("%s/%s/%s", project.ProjectFile.FullPath, project.ProjectFile.SubDir, project.Entities[i].AllLower)
		if _, err := os.Stat(savePath); !os.IsNotExist(err) {
			fmt.Printf("Object: %s name already exists, skipping!\n", project.Entities[i].AllLower)
			continue
		}
		if errMakeAll := os.MkdirAll(savePath, os.ModeDir|os.ModePerm); errMakeAll != nil {
			fmt.Printf("Object: %s path was not able to be made: %s\n", project.Entities[i].AllLower, errMakeAll)
			continue
		}
		project.CurrentEntity = project.Entities[i]

		// build the templates
		buildModelTemplate(project)
		buildRestTemplate(project)
		buildManagerTemplate(project)
		buildStorageTemplate(project)
		buildGrpc(project)
		buildAPIHooks(project)
		processTemplateFiles(*project, savePath)
		buildScriptFiles(*project, project.Entities[i])

		project.CurrentEntity = m.Entity{} // blank it out
	}
	UpdateModFiles(project.ProjectFile.AppName)
	// in case you a entity is marked as 'blank'
	project.UseBlank = false
	project.ProjectFile.SaveProjectFile()
}

func BuildStorage(project *m.Project) {
	storagePath := fmt.Sprintf("%s/internal/storage", project.ProjectFile.FullPath)
	if errMakeAll := os.MkdirAll(storagePath, os.ModeDir|os.ModePerm); errMakeAll != nil {
		fmt.Println("New storage folder was not able to be made", errMakeAll)
		return
	}
	storageVars := m.StorageVars{}
	storageFiles := []string{}
	if project.ProjectFile.UseORM {
		storageFiles = append(storageFiles, "gorm")
	}
	switch project.ProjectFile.Storage {
	case "s":
		if !project.ProjectFile.UseORM {
			tmplFiles = append(tmplFiles, "sql")
		}
		switch project.ProjectFile.SqlStorage {
		case "p":
			storageVars.SQLProvider = c.POSTGRESQL
			storageVars.SQLProviderLower = c.POSTGRESQLLOWER
			storageVars.SQLProviderConnection = fmt.Sprintf("%sConnection", c.POSTGRESQL)
			storageFiles = append(storageFiles, "psql")
		case "m":
			storageVars.SQLProvider = c.MYSQL
			storageVars.SQLProviderLower = c.MYSQLLOWER
			storageVars.SQLProviderConnection = fmt.Sprintf("%sConnection", c.MYSQL)
			storageFiles = append(storageFiles, "mysql")
		case "s":
			storageVars.SQLProvider = c.SQLITE3
			storageVars.SQLProviderLower = c.SQLITE3LOWER
			storageVars.SQLProviderConnection = fmt.Sprintf("%sConnection", c.SQLITE3)
			storageFiles = append(storageFiles, "sqlite")
		}
	case "f":
		storageFiles = append(storageFiles, "file")
		tmplFiles = append(tmplFiles, "file")
	case "m":
		storageFiles = append(storageFiles, "mongo")
		tmplFiles = append(tmplFiles, "mongo")
	}
	project.StorageVars = storageVars

	// save and template storage files
	for _, tmpl := range storageFiles {
		newFileName := fmt.Sprintf("%s/%s.go", storagePath, tmpl)
		// don't over-write if already there
		if _, err := os.Stat(newFileName); !os.IsNotExist(err) {
			// fmt.Printf("already exists: %s, skipping\n", newFileName)
			continue
		}
		file, err := os.Create(newFileName)
		if err != nil {
			fmt.Println("File:", tmpl, "was not able to be created", err)
			continue
		}
		tmplPath := fmt.Sprintf("%s/storage/%s.tmpl", templatePath, tmpl)
		t, errParse := template.ParseFiles(tmplPath)
		if errParse != nil {
			fmt.Printf("Template storage could not parse file: %s; %s", tmplPath, errParse)
			continue
		}
		err = t.Execute(file, project)
		if err != nil {
			fmt.Println("Execution of template:", err)
		}
	}
	return
}

func processTemplateFiles(project m.Project, savePath string) {
	blankInsert := ""
	if project.UseBlank {
		blankInsert = "_blank"
	}
	for _, tmpl := range tmplFiles {
		tmplPath := fmt.Sprintf("%s/%s/%s%s.tmpl", templatePath, tmpl, tmpl, blankInsert)
		if project.CurrentEntity.ModuleName != "" {
			tmplPath = fmt.Sprintf("%s/modules/%s/templates/%s.tmpl", os.Getenv("FORGE_PATH"), project.CurrentEntity.ModuleName, tmpl)
		}
		t, errParse := template.ParseFiles(tmplPath)
		if errParse != nil {
			fmt.Printf("Template could not parse file: %s; %s", tmplPath, errParse)
			fmt.Println("Exiting...")
			return
		}
		newFileName := fmt.Sprintf("%s/%s.go", savePath, tmpl)
		//fmt.Println("New file", newFileName, "creating...")
		file, err := os.Create(newFileName)
		if err != nil {
			fmt.Println("File:", tmpl, "was not able to be created", err)
			fmt.Println("Exiting...")
			return
		}
		err = t.Execute(file, project)
		if err != nil {
			fmt.Println("Execution of template:", err)
		}
		// process _test file
		if !project.UseBlank {
			tmplPath = fmt.Sprintf("%s/%s/%s_test.tmpl", templatePath, tmpl, tmpl)
			if project.CurrentEntity.ModuleName != "" {
				tmplPath = fmt.Sprintf("%s/modules/%s/templates/%s_test.tmpl", os.Getenv("FORGE_PATH"), project.CurrentEntity.ModuleName, tmpl)
			}
			if _, err := os.Stat(tmplPath); !os.IsNotExist(err) {
				t, errParse := template.ParseFiles(tmplPath)
				if errParse != nil {
					fmt.Printf("Template could not parse file: %s; %s", tmplPath, errParse)
					fmt.Println("Exiting...")
					return
				}
				newFileName := fmt.Sprintf("%s/%s_test.go", savePath, tmpl)
				file, err := os.Create(newFileName)
				if err != nil {
					fmt.Println("File:", tmpl, "was not able to be created", err)
					fmt.Println("Exiting...")
					return
				}
				err = t.Execute(file, project)
				if err != nil {
					fmt.Println("Execution of template:", err)
				}
			}
		}
	}
}

func UpdateModFiles(projectName string) {
	// this assumes we are in the root folder
	commands := []*exec.Cmd{
		exec.Command("protoc", "--go_out=./pkg/proto", "--go-grpc_out=./pkg/proto", fmt.Sprintf("./pkg/proto/%s.proto", projectName)),
		// exec.Command("go", "get", "-u", "all"),
		exec.Command("go", "mod", "tidy"),
		exec.Command("go", "fmt", "./..."),
		exec.Command("go", "generate", "./..."),
		// exec.Command("go", "get", "-u", "all"),
		exec.Command("go", "mod", "tidy"),
	}
	for _, command := range commands {
		output, err := command.CombinedOutput()
		if err != nil {
			fmt.Printf("command: %s failed\n", command.String())
			fmt.Printf("output: %s\n", output)
		}
	}
	fileName := "pkg/proto/reset_proto.sh"
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		// push out the helper file, has not been created
		content := fmt.Sprintf(`#! /bin/bash
cd ../.. && protoc --go_out=./pkg/proto --go-grpc_out=./pkg/proto ./pkg/proto/%s.proto`, projectName)
		err := os.WriteFile(fileName, []byte(content), 0766)
		if err != nil {
			fmt.Println("Error creating reset_proto.sh")
		}
	}
}

func buildScriptFiles(p m.Project, entity m.Entity) {
	scriptDir := fmt.Sprintf("%s/scripts/migrations", p.ProjectFile.FullPath)
	if err := os.MkdirAll(scriptDir, os.ModePerm); err != nil {
		fmt.Println("Creating scripts/migrations dir", err)
		return
	}
	if entity.ModuleName != "" {
		return
	}
	if len(entity.SqlLines) == 0 {
		return
	}
	now := time.Now().Format("20060102150405")
	fileName := fmt.Sprintf("%s/%s-create-table-%s.sql", scriptDir, now, normalizeName(entity.RawName))
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("Unable to create %s.sql file: %s\n", fileName, err)
		return
	}
	fileScript := []string{}
	for i := range entity.SqlLines {
		if !(i == 0 || i == len(entity.SqlLines)-1) {
			// put a tab at the front expect the first and last row
			fileScript = append(fileScript, fmt.Sprintf("\t%s", entity.SqlLines[i]))
			continue
		}
		fileScript = append(fileScript, entity.SqlLines[i])
	}
	f.WriteString(strings.Join(fileScript, "\n"))
	f.Close()
}

func normalizeName(fileName string) (normalizedName string) {
	normalizedName = strings.ReplaceAll(fileName, " ", "-")
	normalizedName = strings.ReplaceAll(normalizedName, "_", "-")
	if len(normalizedName) > 85 {
		normalizedName = normalizedName[:85]
	}
	return
}
