package entity

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"strings"

	con "github.com/blackflagsoftware/forge-go/internal/constant"
)

func (ep *Entity) BuildAPIHooks() {
	// hook into rest main
	apiFile := fmt.Sprintf("%s/cmd/rest/main.go", ep.ProjectFile.FullPath)
	if _, err := os.Stat(apiFile); os.IsNotExist(err) {
		fmt.Printf("%s is missing unable to write in hooks\n", apiFile)
	} else {
		var serverReplace bytes.Buffer
		tServer := template.Must(template.New("server").Parse(con.SERVER_ROUTE))
		errServer := tServer.Execute(&serverReplace, ep)
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
		errServer = tMain.Execute(&mainReplace, ep)
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
		if ep.SQLProvider != "" {
			// handling sql add migration
			nonSqlite := con.MIGRATION_NON_SQLITE
			if ep.SQLProvider == "Sqlite" {
				nonSqlite = "\n\t\t\tc.Host = config.SqlitePath"
			}
			migCall := fmt.Sprintf(con.MIGRATION_CALL, nonSqlite, ep.SQLProviderLower)
			migRestMain := fmt.Sprintf(`perl -pi -e 's/\/\/ --- replace migration once text - do not remove ---/%s/g' %s`, migCall, apiFile)
			execRestMain := exec.Command("bash", "-c", migRestMain)
			errExecRestMain := execRestMain.Run()
			if errExecRestMain != nil {
				fmt.Printf("%s: error in replace migration main [%s]\n", apiFile, errExecRestMain)
			}
			migRestHeader := fmt.Sprintf(`perl -pi -e 's/\/\/ --- replace migration header once text - do not remove ---/%s/g' %s`, fmt.Sprintf(`mig "%s\/tools\/migration\/src"`, ep.ProjectPathEncoded), apiFile)
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
	errHeader := tHeader.Execute(&headerReplace, ep)
	if errHeader != nil {
		fmt.Println("Header template error:", errHeader)
		return
	}
	// section
	var sectionReplace bytes.Buffer
	tSection := template.Must(template.New("section").Parse(con.COMMON_SECTION))
	errSection := tSection.Execute(&sectionReplace, ep)
	if errSection != nil {
		fmt.Println("Section template error:", errSection)
		return
	}
	// hook into grpc file
	grpcFile := fmt.Sprintf("%s/cmd/grpc/main.go", ep.ProjectFile.FullPath)
	if _, err := os.Stat(grpcFile); os.IsNotExist(err) {
		fmt.Printf("%s is missing unable to write in hooks\n", grpcFile)
	} else {
		var grpcReplace bytes.Buffer
		tGrpc := template.Must(template.New("grpc").Parse(con.GRPC_TEXT))
		errGrpc := tGrpc.Execute(&grpcReplace, ep)
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
		errGrpc = tImport.Execute(&importReplace, ep)
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
		errGrpc = tOnce.Execute(&importOnceReplace, ep)
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
		if ep.SQLProvider != "" {
			// handling sql add migration
			nonSqlite := con.MIGRATION_NON_SQLITE
			if ep.SQLProvider == "Sqlite" {
				nonSqlite = "\n\t\t\tc.Host = config.SqlitePath"
			}
			migCall := fmt.Sprintf(con.MIGRATION_CALL, nonSqlite, ep.SQLProviderLower)
			migGrpcMain := fmt.Sprintf(`perl -pi -e 's/\/\/ --- replace migration once text - do not remove ---/%s/g' %s`, migCall, grpcFile)
			execGrpcMain := exec.Command("bash", "-c", migGrpcMain)
			errExecGrpcMain := execGrpcMain.Run()
			if errExecGrpcMain != nil {
				fmt.Printf("%s: error in replace migration main [%s]\n", grpcFile, errExecGrpcMain)
			}
			migGrpcHeader := fmt.Sprintf(`perl -pi -e 's/\/\/ --- replace migration header once text - do not remove ---/%s/g' %s`, fmt.Sprintf(`mig "%s\/tools\/migration\/src"`, ep.ProjectPathEncoded), grpcFile)
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
	configFile := fmt.Sprintf("%s/config/config.go", ep.ProjectFile.FullPath)
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		fmt.Printf("%s is missing unable to write in hooks\n", configFile)
	} else {
		// env lines
		configLines := []string{}
		switch ep.ProjectFile.Storage {
		case "s":
			configLines = append(configLines, "StorageSQL = true")
			if ep.ProjectFile.SqlStorage == "p" || ep.ProjectFile.SqlStorage == "m" {
				configLines = append(configLines, "DBHost = GetEnvOrDefault(\"{{.ProjectFile.Name.EnvVar}}_DB_HOST\", \"\")")
				configLines = append(configLines, "DBDB = GetEnvOrDefault(\"{{.ProjectFile.Name.EnvVar}}_DB_DB\", \"\")")
				configLines = append(configLines, "DBUser = GetEnvOrDefault(\"{{.ProjectFile.Name.EnvVar}}_DB_USER\",\"\")")
				configLines = append(configLines, "DBPass = GetEnvOrDefault(\"{{.ProjectFile.Name.EnvVar}}_DB_PASS\", \"\")")
				configLines = append(configLines, "AdminDBUser = GetEnvOrDefault(\"{{.ProjectFile.Name.EnvVar}}_ADMIN_DB_USER\",\"\")")
				configLines = append(configLines, "AdminDBPass = GetEnvOrDefault(\"{{.ProjectFile.Name.EnvVar}}_ADMIN_DB_PASS\", \"\")")
			} else {
				configLines = append(configLines, "SqlitePath = GetEnvOrDefault(\"{{.ProjectFile.Name.EnvVar}}_SQLITE_PATH\",\"\")")
			}
		case "f":
			configLines = append(configLines, "StorageFile = true")
			configLines = append(configLines, "SqlitePath = GetEnvOrDefault(\"{{.ProjectFile.Name.EnvVar}}_SQLITE_PATH\", \"/tmp/{{.Name.Lower}}.db\")")
		case "m":
			configLines = append(configLines, "StorageMongo = true")
			configLines = append(configLines, "MongoHost = GetEnvOrDefault(\"{{.ProjectFile.Name.EnvVar}}_MONGO_HOST\", \"localhost\")")
			configLines = append(configLines, "MongoPort = GetEnvOrDefault(\"{{.ProjectFile.Name.EnvVar}}_MONGO_PORT\", \"27017\")")
		}
		configLine := strings.Join(configLines, "\n\t")

		var configReplace bytes.Buffer
		tConfig := template.Must(template.New("config").Parse(configLine))
		errConfig := tConfig.Execute(&configReplace, ep)
		if errConfig != nil {
			fmt.Printf("%s: template error [%s]\n", configFile, errConfig)
		} else {
			cmdConfig := fmt.Sprintf(`perl -pi -e 's/\/\/ --- replace config text - do not remove ---/%s/g' %s`, configReplace.String(), configFile)
			execConfig := exec.Command("bash", "-c", cmdConfig)
			errConfigCmd := execConfig.Run()
			if errConfigCmd != nil {
				fmt.Printf("%s: error in replace for config text [%s]\n", configFile, errConfigCmd)
			}
		}
	}
}
