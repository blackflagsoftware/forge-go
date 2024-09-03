package task

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"strings"

	con "github.com/blackflagsoftware/forge-go/internal/constant"
	m "github.com/blackflagsoftware/forge-go/internal/model"
)

func buildAPIHooks(p *m.Project) {
	// hook into rest main
	apiFile := fmt.Sprintf("%s/cmd/rest/main.go", p.ProjectFile.FullPath)
	if _, err := os.Stat(apiFile); os.IsNotExist(err) {
		fmt.Printf("%s is missing unable to write in hooks\n", apiFile)
	} else {
		var serverReplace bytes.Buffer
		tServer := template.Must(template.New("server").Parse(con.SERVER_ROUTE))
		errServer := tServer.Execute(&serverReplace, p)
		if errServer != nil {
			fmt.Printf("%s: template error [%s]\n", apiFile, errServer)
		} else {
			cmdServer := fmt.Sprintf(`perl -pi -e 's/\/\/ --- replace server text - do not remove ---/%s/g' %s`, serverReplace.String(), apiFile)
			execServer := exec.Command("bash", "-c", cmdServer)
			errServerCmd := execServer.Run()
			if errServerCmd != nil {
				fmt.Printf("%s: error in replace for server [%s]\n", apiFile, errServerCmd)
			}
		}
		onceReplace := `routeGroup := e.Group("v1") \/\/ change to match your uri prefix`
		cmdOnceServer := fmt.Sprintf(`perl -pi -e 's/\/\/ --- replace server once text - do not remove ---/%s/g' %s`, onceReplace, apiFile)
		execServerOnce := exec.Command("bash", "-c", cmdOnceServer)
		errServerOnceCmd := execServerOnce.Run()
		if errServerOnceCmd != nil {
			fmt.Printf("%s: error in replace for server once [%s]\n", apiFile, errServerOnceCmd)
		}
		var mainReplace bytes.Buffer
		tMain := template.Must(template.New("server").Parse(con.MAIN_COMMON_PATH))
		errServer = tMain.Execute(&mainReplace, p)
		if errServer != nil {
			fmt.Printf("%s: template error [%s]\n", apiFile, errServer)
		} else {
			cmdServer := fmt.Sprintf(`perl -pi -e 's/\/\/ --- replace main header text - do not remove ---/%s/g' %s`, mainReplace.String(), apiFile)
			execServer := exec.Command("bash", "-c", cmdServer)
			errServerCmd := execServer.Run()
			if errServerCmd != nil {
				fmt.Printf("%s: error in replace for main [%s]\n", apiFile, errServerCmd)
			}
		}
		if p.SQLProvider != "" {
			// handling sql add migration
			nonSqlite := con.MIGRATION_NON_SQLITE
			if p.SQLProvider == "Sqlite" {
				nonSqlite = "\n\t\t\tHost: config.SqlitePath,"
			}
			migCall := fmt.Sprintf(con.MIGRATION_CALL, nonSqlite)
			migRestMain := fmt.Sprintf(`perl -pi -e 's/\/\/ --- replace migration once text - do not remove ---/%s/g' %s`, migCall, apiFile)
			execRestMain := exec.Command("bash", "-c", migRestMain)
			errExecRestMain := execRestMain.Run()
			if errExecRestMain != nil {
				fmt.Printf("%s: error in replace migration main [%s]\n", apiFile, errExecRestMain)
			}
			migRestHeader := fmt.Sprintf(`perl -pi -e 's/\/\/ --- replace migration header once text - do not remove ---/%s/g' %s`, fmt.Sprintf(`mig "%s\/tools\/migration\/src"`, p.ProjectPathEncoded), apiFile)
			execRestHeader := exec.Command("bash", "-c", migRestHeader)
			errExecRestHeader := execRestHeader.Run()
			if errExecRestHeader != nil {
				fmt.Printf("%s: error in replace migration main [%s]\n", apiFile, errExecRestMain)
			}
		}
	}
	// header
	var headerReplace bytes.Buffer
	tHeader := template.Must(template.New("header").Parse(con.COMMON_HEADER))
	errHeader := tHeader.Execute(&headerReplace, p)
	if errHeader != nil {
		fmt.Println("Header template error:", errHeader)
		return
	}
	// section
	var sectionReplace bytes.Buffer
	tSection := template.Must(template.New("section").Parse(con.COMMON_SECTION))
	errSection := tSection.Execute(&sectionReplace, p)
	if errSection != nil {
		fmt.Println("Section template error:", errSection)
		return
	}
	// hook into grpc file
	grpcFile := fmt.Sprintf("%s/cmd/grpc/main.go", p.ProjectFile.FullPath)
	if _, err := os.Stat(grpcFile); os.IsNotExist(err) {
		fmt.Printf("%s is missing unable to write in hooks\n", grpcFile)
	} else {
		if !p.CurrentEntity.SkipGrpc {
			var grpcReplace bytes.Buffer
			tGrpc := template.Must(template.New("grpc").Parse(con.GRPC_TEXT))
			errGrpc := tGrpc.Execute(&grpcReplace, p)
			if errGrpc != nil {
				fmt.Printf("%s: template error [%s]\n", grpcFile, errGrpc)
			} else {
				cmdGrpc := fmt.Sprintf(`perl -pi -e 's/\/\/ --- replace grpc text - do not remove ---/%s/g' %s`, grpcReplace.String(), grpcFile)
				execGrpc := exec.Command("bash", "-c", cmdGrpc)
				errGrpcCmd := execGrpc.Run()
				if errGrpcCmd != nil {
					fmt.Printf("%s: error in replace for grpc text [%s]\n", grpcFile, errGrpcCmd)
				}
			}
			var importReplace bytes.Buffer
			tImport := template.Must(template.New("grpc").Parse(con.GRPC_IMPORT))
			errGrpc = tImport.Execute(&importReplace, p)
			if errGrpc != nil {
				fmt.Printf("%s: template error [%s]\n", grpcFile, errGrpc)
			} else {
				cmdGrpc := fmt.Sprintf(`perl -pi -e 's/\/\/ --- replace grpc import - do not remove ---/%s/g' %s`, importReplace.String(), grpcFile)
				execGrpc := exec.Command("bash", "-c", cmdGrpc)
				errGrpcCmd := execGrpc.Run()
				if errGrpcCmd != nil {
					fmt.Printf("%s: error in replace for grpc [%s]\n", grpcFile, errGrpcCmd)
				}
			}
			var importOnceReplace bytes.Buffer
			tOnce := template.Must(template.New("grpc").Parse(con.GRPC_IMPORT_ONCE))
			errGrpc = tOnce.Execute(&importOnceReplace, p)
			if errGrpc != nil {
				fmt.Printf("%s: template error [%s]\n", grpcFile, errGrpc)
			} else {
				cmdGrpc := fmt.Sprintf(`perl -pi -e 's/\/\/ --- replace grpc import once - do not remove ---/%s/g' %s`, importOnceReplace.String(), grpcFile)
				execGrpc := exec.Command("bash", "-c", cmdGrpc)
				errGrpcCmd := execGrpc.Run()
				if errGrpcCmd != nil {
					fmt.Printf("%s: error in replace for grpc [%s]\n", grpcFile, errGrpcCmd)
				}
			}
		}
		if p.SQLProvider != "" {
			// handling sql add migration
			nonSqlite := con.MIGRATION_NON_SQLITE
			if p.SQLProvider == "Sqlite" {
				nonSqlite = "\n\t\t\tHost: config.SqlitePath,"
			}
			migCall := fmt.Sprintf(con.MIGRATION_CALL, nonSqlite)
			migGrpcMain := fmt.Sprintf(`perl -pi -e 's/\/\/ --- replace migration once text - do not remove ---/%s/g' %s`, migCall, grpcFile)
			execGrpcMain := exec.Command("bash", "-c", migGrpcMain)
			errExecGrpcMain := execGrpcMain.Run()
			if errExecGrpcMain != nil {
				fmt.Printf("%s: error in replace migration main [%s]\n", grpcFile, errExecGrpcMain)
			}
			migGrpcHeader := fmt.Sprintf(`perl -pi -e 's/\/\/ --- replace migration header once text - do not remove ---/%s/g' %s`, fmt.Sprintf(`mig "%s\/tools\/migration\/src"`, p.ProjectPathEncoded), grpcFile)
			execGrpcHeader := exec.Command("bash", "-c", migGrpcHeader)
			errExecGrpcHeader := execGrpcHeader.Run()
			if errExecGrpcHeader != nil {
				fmt.Printf("%s: error in replace migration main [%s]\n", grpcFile, errExecGrpcHeader)
			}
			migGrpcHeaderOs := fmt.Sprintf(`perl -pi -e 's/\/\/ --- replace migration header os once text - do not remove ---/%s/g' %s`, con.MIGRATION_GRPC_HEADER_ONCE, grpcFile)
			execGrpcHeaderOs := exec.Command("bash", "-c", migGrpcHeaderOs)
			errExecGrpcHeaderOs := execGrpcHeaderOs.Run()
			if errExecGrpcHeaderOs != nil {
				fmt.Printf("%s: error in replace migration main [%s]\n", grpcFile, errExecGrpcHeaderOs)
			}
		}
	}
	// hook into config.go
	PopulateConfig(&p.ProjectFile)
	// audit
	auditFile := fmt.Sprintf("%s/internal/audit/audit.go", p.ProjectFile.FullPath)
	if _, err := os.Stat(auditFile); os.IsNotExist(err) {
		fmt.Printf("%s is missing unable to write in hooks\n", auditFile)
	} else {
		onceReplace := ""
		if p.SQLProvider != "" {
			onceReplace = fmt.Sprintf(`stor "%s\/internal\/storage"`, p.ProjectPathEncoded)
		}
		auditOnce := fmt.Sprintf(`perl -pi -e 's/\/\/ --- replace storage once text - do not remove ---/%s/g' %s`, onceReplace, auditFile)
		execServerOnce := exec.Command("bash", "-c", auditOnce)
		errServerOnceCmd := execServerOnce.Run()
		if errServerOnceCmd != nil {
			fmt.Printf("%s: error in replace for audit once [%s]\n", auditFile, errServerOnceCmd)
		}
		onceReference := ""
		if p.SQLProvider != "" {
			onceReference = fmt.Sprintf("DB: stor.%sInit(),", p.SQLProvider)
		}
		auditOnceReference := fmt.Sprintf(`perl -pi -e 's/\/\/ --- replace storage reference once text - do not remove ---/%s/g' %s`, onceReference, auditFile)
		execServerOnce = exec.Command("bash", "-c", auditOnceReference)
		errServerOnceCmd = execServerOnce.Run()
		if errServerOnceCmd != nil {
			fmt.Printf("%s: error in replace for audit reference once [%s]\n", auditFile, errServerOnceCmd)
		}
	}
}

func PopulateConfig(projectFile *m.ProjectFile) {
	// only run this once, sets this to true at the end of this function
	if projectFile.LoadedConfig {
		return
	}
	configFile := fmt.Sprintf("%s/config/config.go", projectFile.FullPath)
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		fmt.Printf("%s is missing unable to write in hooks\n", configFile)
	} else {
		// env lines
		configVarLines := []string{}
		configInitLines := []string{}
		switch projectFile.Storage {
		case "s":
			lowerSqlEngine := strings.ToLower(m.SqlTypeToProper(projectFile.SqlStorage))
			configVarLines = append(configVarLines, "StorageSQL = true")
			configVarLines = append(configVarLines, "DBEngine string")
			configInitLines = append(configInitLines, fmt.Sprintf("DBEngine = GetEnvOrDefault(\"{{.Name.EnvVar}}_DB_ENGINE\", \"%s\")", lowerSqlEngine))
			if projectFile.SqlStorage == "p" || projectFile.SqlStorage == "m" {
				configVarLines = append(configVarLines, "DBHost string")
				configInitLines = append(configInitLines, "DBHost = GetEnvOrDefault(\"{{.Name.EnvVar}}_DB_HOST\", \"\")")
				configVarLines = append(configVarLines, "DBDB string")
				configInitLines = append(configInitLines, "DBDB = GetEnvOrDefault(\"{{.Name.EnvVar}}_DB_DB\", \"\")")
				configVarLines = append(configVarLines, "DBUser string")
				configInitLines = append(configInitLines, "DBUser = GetEnvOrDefault(\"{{.Name.EnvVar}}_DB_USER\",\"\")")
				configVarLines = append(configVarLines, "DBPass string")
				configInitLines = append(configInitLines, "DBPass = GetEnvOrDefault(\"{{.Name.EnvVar}}_DB_PASS\", \"\")")
				port := "3306" // default to mysql
				if projectFile.SqlStorage == "p" {
					port = "5432"
				}
				configVarLines = append(configVarLines, "DBPort string")
				configInitLines = append(configInitLines, fmt.Sprintf("DBPort = GetEnvOrDefault(\"{{.Name.EnvVar}}_DB_PORT\", \"%s\")", port))
				configVarLines = append(configVarLines, "AdminDBUser string")
				configInitLines = append(configInitLines, "AdminDBUser = GetEnvOrDefault(\"{{.Name.EnvVar}}_ADMIN_DB_USER\",\"\")")
				configVarLines = append(configVarLines, "AdminDBPass string")
				configInitLines = append(configInitLines, "AdminDBPass = GetEnvOrDefault(\"{{.Name.EnvVar}}_ADMIN_DB_PASS\", \"\")")
			} else {
				configVarLines = append(configVarLines, "SqlitePath string")
				configInitLines = append(configInitLines, "SqlitePath = GetEnvOrDefault(\"{{.Name.EnvVar}}_SQLITE_PATH\",\"\")")
			}
		case "f":
			configVarLines = append(configVarLines, "StorageFile = true")
			configVarLines = append(configVarLines, "StorageFileDir string")
			configInitLines = append(configInitLines, `StorageFileDir = GetEnvOrDefault("{{.Name.EnvVar}}_STORAGE_FILE_DIR", "\/tmp")`)
		case "m":
			configVarLines = append(configVarLines, "StorageMongo = true")
			configVarLines = append(configVarLines, "MongoHost string")
			configInitLines = append(configInitLines, "MongoHost = GetEnvOrDefault(\"{{.Name.EnvVar}}_MONGO_HOST\", \"localhost\")")
			configVarLines = append(configVarLines, "MongoPort string")
			configInitLines = append(configInitLines, "MongoPort = GetEnvOrDefault(\"{{.Name.EnvVar}}_MONGO_PORT\", \"27017\")")
		}
		configVarLines = append(configVarLines, `\/\/ --- replace config var text - do not remove ---`)
		configInitLines = append(configInitLines, `\/\/ --- replace config init text - do not remove ---`)
		configVarLine := strings.Join(configVarLines, "\n\t")
		configInitLine := strings.Join(configInitLines, "\n\t")

		var varReplace bytes.Buffer
		var initReplace bytes.Buffer
		// config var lines
		tConfig := template.Must(template.New("config-var").Parse(configVarLine))
		errConfig := tConfig.Execute(&varReplace, projectFile)
		if errConfig != nil {
			fmt.Printf("%s: template error [%s]\n", configFile, errConfig)
		} else {
			cmdConfig := fmt.Sprintf(`perl -pi -e 's/\/\/ --- replace config var text - do not remove ---/%s/g' %s`, varReplace.String(), configFile)
			execConfig := exec.Command("bash", "-c", cmdConfig)
			errConfigCmd := execConfig.Run()
			if errConfigCmd != nil {
				fmt.Printf("%s: error in replace for config var text [%s]\n", configFile, errConfigCmd)
			}
		}
		// config init lines
		tConfig = template.Must(template.New("config-init").Parse(configInitLine))
		errConfig = tConfig.Execute(&initReplace, projectFile)
		if errConfig != nil {
			fmt.Printf("%s: template error [%s]\n", configFile, errConfig)
		} else {
			cmdConfig := fmt.Sprintf(`perl -pi -e 's/\/\/ --- replace config init text - do not remove ---/%s/g' %s`, initReplace.String(), configFile)
			execConfig := exec.Command("bash", "-c", cmdConfig)
			errConfigCmd := execConfig.Run()
			if errConfigCmd != nil {
				fmt.Printf("%s: error in replace for config init text [%s]\n", configFile, errConfigCmd)
			}
		}
		projectFile.LoadedConfig = true
		projectFile.SaveProjectFile()
	}
}
