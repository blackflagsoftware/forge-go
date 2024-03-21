package config

import (
	"bytes"
	"fmt"
	"os"
	"path"

	// --- replace config import text - do not remove ---
	"github.com/kardianos/osext"
)

var (
	RootDir           = "full-path-here" // this is good for local dev but it used below to read the .env files so setting it for production is vital
	AppName           = "forge-go-base"
	AppVersion        string
	RestPort          string
	GrpcPort          string
	PidPath           string
	Env               string
	LogPath           string
	EnableMetrics     bool
	UseMigration      bool
	MigrationPath     string
	MigrationSkipInit bool
	EnableAuditing    bool
	AuditStorage      string
	AuditFilePath     string
	LogOutput         = os.Stdout
	ExecDir           = ""
	// --- replace config var text - do not remove ---
)

func init() {
	ExecDir, _ = osext.ExecutableFolder()
	loadEnvFiles()
	loadEnvVars()
}

func loadEnvVars() {
	RootDir = GetEnvOrDefault("FORGE_GO_BASE_ROOT_DIR", "")
	AppVersion = GetEnvOrDefault("FORGE_GO_BASE_APP_VERSION", "1.0.0")
	RestPort = GetEnvOrDefault("FORGE_GO_BASE_REST_PORT", "12580")
	GrpcPort = GetEnvOrDefault("FORGE_GO_BASE_GRPC_PORT", "12581")
	PidPath = GetEnvOrDefault("FORGE_GO_BASE_PID_PATH", fmt.Sprintf("/tmp/%s.pid", AppName))
	Env = GetEnvOrDefault("FORGE_GO_BASE_ENV", "dev")
	LogPath = GetEnvOrDefault("FORGE_GO_BASE_LOG_PATH", fmt.Sprintf("/tmp/%s.out", AppName))
	EnableMetrics = GetEnvOrDefaultBool("FORGE_GO_BASE_ENABLE_METRICS", true)
	UseMigration = GetEnvOrDefaultBool("FORGE_GO_BASE_MIGRATION_ENABLED", false)
	MigrationPath = GetEnvOrDefault("FORGE_GO_BASE_MIGRATION_PATH", "")
	MigrationSkipInit = GetEnvOrDefaultBool("FORGE_GO_BASE_MIGRATION_SKIP_INIT", false)
	EnableAuditing = GetEnvOrDefaultBool("FORGE_GO_BASE_ENABLE_AUDITING", false)
	AuditStorage = GetEnvOrDefault("FORGE_GO_BASE_AUDIT_STORAGE", "file") // file or sql
	AuditFilePath = GetEnvOrDefault("FORGE_GO_BASE_AUDIT_FILE_PATH", "./audit")
	// --- replace config init text - do not remove ---
}

func loadEnvFiles() {
	// load any and all .env.* files if present at the root level of the binary
	// the order of presedence goes from local to just plain .env, the later WILL override earlier
	// if the env var is already declared at the console/terminal level, this will override it too
	envFiles := []string{".env.local", ".env.dev", ".env.qa", ".env.prod", ".env"} // if you want another name, just add it to this list in the order of precedence
	for _, envFile := range envFiles {
		envFileFull := path.Join(RootDir, envFile)
		if _, err := os.Stat(envFileFull); !os.IsNotExist(err) {
			// found the file
			loadEnvVarsFromFile(envFileFull)
		}
	}
}

func loadEnvVarsFromFile(fileName string) {
	// line format: name=value
	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Printf("env file: %s could not be read\n", fileName)
		return
	}
	lines := bytes.Split(fileContent, []byte("\n"))
	for _, line := range lines {
		if len(line) > 0 {
			if line[0] != byte('#') {
				// not a comment
				split := bytes.Split(line, []byte("="))
				if len(split) == 2 {
					os.Setenv(string(split[0]), string(split[1]))
				}
			}
		}
	}
}

func GetEnvOrDefault(envVar string, defEnvVar string) (newEnvVar string) {
	if newEnvVar = os.Getenv(envVar); len(newEnvVar) == 0 {
		return defEnvVar
	} else {
		return newEnvVar
	}
}

func GetEnvOrDefaultBool(envVar string, defEnvVar bool) (newEnvVar bool) {
	newEnvVarStr := os.Getenv(envVar)
	if len(newEnvVarStr) == 0 {
		return defEnvVar
	}
	return newEnvVarStr == "true"
}

func GetUniqueNumberForLock() (number int) {
	for i := range AppName {
		number += int(AppName[i])
	}
	return
}
