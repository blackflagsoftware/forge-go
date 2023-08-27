package config

import (
	"fmt"
	"os"

	"github.com/kardianos/osext"
)

var (
	AppName           = "forge-go-base"
	AppVersion        = GetEnvOrDefault("FORGE_GO_BASE_APP_VERSION", "1.0.0")
	RestPort          = GetEnvOrDefault("FORGE_GO_BASE_REST_PORT", "12580")
	GrpcPort          = GetEnvOrDefault("FORGE_GO_BASE_GRPC_PORT", "12581")
	PidPath           = GetEnvOrDefault("FORGE_GO_BASE_PID_PATH", fmt.Sprintf("/tmp/%s.pid", AppName))
	Env               = GetEnvOrDefault("FORGE_GO_BASE_ENV", "dev")
	LogPath           = GetEnvOrDefault("FORGE_GO_BASE_LOG_PATH", fmt.Sprintf("/tmp/%s.out", AppName))
	EnableMetrics     = GetEnvOrDefaultBool("FORGE_GO_BASE_ENABLE_METRICS", true)
	UseMigration      = GetEnvOrDefaultBool("FORGE_GO_BASE_MIGRATION_ENABLED", true)
	MigrationPath     = GetEnvOrDefault("FORGE_GO_BASE_MIGRATION_PATH", "")
	MigrationSkipInit = GetEnvOrDefaultBool("FORGE_GO_BASE_MIGRATION_SKIP_INIT", false)
	EnableAuditing    = GetEnvOrDefaultBool("FORGE_GO_BASE_ENABLE_AUDITING", false)
	AuditStorage      = GetEnvOrDefault("FORGE_GO_BASE_AUDIT_STORAGE", "file") // file or sql
	AuditFilePath     = GetEnvOrDefault("FORGE_GO_BASE_AUDIT_FILE_PATH", "./audit")
	LogOutput         = os.Stdout
	ExecDir           = ""
	// --- replace config text - do not remove ---
)

func init() {
	ExecDir, _ = osext.ExecutableFolder()
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
