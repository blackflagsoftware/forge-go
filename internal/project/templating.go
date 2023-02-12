package project

import (
	"fmt"
	"os"
	"os/exec"
	"text/template"

	c "github.com/blackflagsoftware/forge-go/internal/constant"
	e "github.com/blackflagsoftware/forge-go/internal/entity"
)

var (
	templatePath = fmt.Sprintf("%s/templates", os.Getenv("FORGE_PATH"))
	tmplFiles    = []string{"model", "rest", "manager", "grpc"} // TODO: sql is optional depending on which data storages they want to template, dynamically build this
)

func (project *Project) StartTemplating() {
	sqlProvider := buildStorage(*project)
	buildMigration(*project)

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
		project.Entities[i].SQLProvider = sqlProvider
		project.Entities[i].ProjectFile = project.ProjectFile // TODO: can we get away from this

		// build the templates
		buildTemplateParts(&project.Entities[i])
		processTemplateFiles(*project, &project.Entities[i], savePath)
	}
	updateModFiles(project.ProjectFile.AppName)
}

// send back SQLProvider
func buildStorage(project Project) (sqlProvider string) {
	storagePath := fmt.Sprintf("%s/internal/storage", project.ProjectFile.FullPath)
	if errMakeAll := os.MkdirAll(storagePath, os.ModeDir|os.ModePerm); errMakeAll != nil {
		fmt.Println("New storage folder was not able to be made", errMakeAll)
		return
	}
	storageVars := StorageVars{ProjectPath: project.ProjectFile.ProjectPath}
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
	case "m":
		storageFiles = append(storageFiles, "mongo")
	}
	sqlProvider = storageVars.SQLProvider

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
		err = t.Execute(file, storageVars)
		if err != nil {
			fmt.Println("Execution of template:", err)
		}
	}
	return
}

// TODO: write to migration to have all varieties of code an set a env var to determine which sql engine to support
func buildMigration(project Project) {
	if project.ProjectFile.Storage == "s" {
		migrationVars := MigrationVars{ProjectPath: project.ProjectFile.ProjectPath, ProjectFile: project.ProjectFile}
		if project.ProjectFile.UseORM {
			tmplFiles = append(tmplFiles, "gorm")
		} else {
			tmplFiles = append(tmplFiles, "sql")
		}
		switch project.ProjectFile.SqlStorage {
		case "m":
			migrationVars.MigrationVerify = c.MIGRATION_VERIFY_MYSQL
			migrationVars.MigrationConnection = c.MIGRATION_CONNECTION_MYSQL
			migrationVars.MigrationHeader = c.MIGRATION_VERIFY_HEADER_MYSQL
		case "p":
			migrationVars.MigrationVerify = c.MIGRATION_VERIFY_POSTGRES
			migrationVars.MigrationConnection = c.MIGRATION_CONNECTION_POSTGRES
			migrationVars.MigrationHeader = c.MIGRATION_VERIFY_HEADER_POSTGRES
		case "s":
			migrationVars.MigrationVerify = c.MIGRATION_VERIFY_SQLITE
			migrationVars.MigrationHeader = c.MIGRATION_VERIFY_HEADER_SQLITE
		}
		// TODO: if they change the sql engine, should we re-do this code
		// save migration
		migPath := fmt.Sprintf("%s/tools/migration", project.ProjectFile.FullPath)
		if _, err := os.Stat(migPath); os.IsNotExist(err) {
			errMk := os.MkdirAll(migPath, 0755)
			if errMk != nil {
				fmt.Printf("Unable to make tools/migration path: %s", errMk)
				return
			}
			// migration/main.go
			tmplPath := fmt.Sprintf("%s/tools/migration_main.tmpl", templatePath)
			t, errParse := template.ParseFiles(tmplPath)
			if errParse != nil {
				fmt.Printf("Template migration/main could not parse file: %s; %s", tmplPath, errParse)
				return
			}
			newFileName := fmt.Sprintf("%s/main.go", migPath)
			file, err := os.Create(newFileName)
			if err != nil {
				fmt.Println("File: migration main was not able to be created", err)
				return
			}
			err = t.Execute(file, migrationVars)
			if err != nil {
				fmt.Println("Execution of template:", err)
			}
			// migration/src/main.go
			tmplPath = fmt.Sprintf("%s/tools/migration.tmpl", templatePath)
			t, errParse = template.ParseFiles(tmplPath)
			if errParse != nil {
				fmt.Printf("Template src/migration could not parse file: %s; %s", tmplPath, errParse)
				return
			}
			migPath = migPath + "/src"
			errMk = os.MkdirAll(migPath, 0755)
			if errMk != nil {
				fmt.Printf("Unable to make tools/migration/src path: %s", errMk)
				return
			}
			newFileName = fmt.Sprintf("%s/migration.go", migPath)
			file, err = os.Create(newFileName)
			if err != nil {
				fmt.Println("File: migration src was not able to be created", err)
				return
			}
			err = t.Execute(file, migrationVars)
			if err != nil {
				fmt.Println("Execution of template:", err)
			}
		}
	}
}

func buildTemplateParts(ep *e.Entity) {
	ep.BuildModelTemplate()
	ep.BuildRestTemplate()
	ep.BuildManagerTemplate()
	ep.BuildDataTemplate()
	ep.BuildGrpc()
	ep.BuildAPIHooks()
}

func processTemplateFiles(project Project, ep *e.Entity, savePath string) {
	blankInsert := ""
	if project.UseBlank {
		blankInsert = "_blank"
	}
	for _, tmpl := range tmplFiles {
		tmplPath := fmt.Sprintf("%s/%s/%s%s.tmpl", templatePath, tmpl, tmpl, blankInsert)
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
		err = t.Execute(file, ep)
		if err != nil {
			fmt.Println("Execution of template:", err)
		}
		// process _test file
		tmplPath = fmt.Sprintf("%s/%s/%s_test.tmpl", templatePath, tmpl, tmpl)
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
			err = t.Execute(file, ep)
			if err != nil {
				fmt.Println("Execution of template:", err)
			}
		}
	}
}

func updateModFiles(projectName string) {
	// this assumes we are in the root folder
	commands := []*exec.Cmd{
		exec.Command("protoc", "--go_out=./pkg/proto", "--go-grpc_out=./pkg/proto", fmt.Sprintf("./pkg/proto/%s.proto", projectName)),
		exec.Command("go", "get", "-u", "all"),
		exec.Command("go", "mod", "tidy"),
		exec.Command("go", "fmt", "./..."),
		exec.Command("go", "generate", "./..."),
		exec.Command("go", "get", "-u", "all"),
		exec.Command("go", "mod", "tidy"),
	}
	for _, command := range commands {
		output, err := command.CombinedOutput()
		if err != nil {
			fmt.Printf("command: %s failed\n", command.String())
			fmt.Printf("output: %s\n", output)
		}
	}
}
