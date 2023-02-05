package config

import (
	"fmt"
	"os"
	"time"

	"github.com/client9/reopen"
	"github.com/kardianos/osext"
	log "github.com/sirupsen/logrus"
)

var (
	AppName       = "forge-go-base"
	AppVersion    = getEnvOrDefault("FORGE_GO_BASE_APP_VERSION", "1.0.0")
	RestPort      = getEnvOrDefault("FORGE_GO_BASE_REST_PORT", "12580")
	GrpcPort      = getEnvOrDefault("FORGE_GO_BASE_GRPC_PORT", "12581")
	PidPath       = getEnvOrDefault("FORGE_GO_BASE_PID_PATH", fmt.Sprintf("/tmp/%s.pid", AppName))
	Env           = getEnvOrDefault("FORGE_GO_BASE_ENV", "dev")
	LogPath       = getEnvOrDefault("FORGE_GO_BASE_LOG_PATH", fmt.Sprintf("/tmp/%s.out", AppName))
	UseMigration  = true
	MigrationPath = getEnvOrDefault("FORGE_GO_BASE_MIGRATION_PATH", "")
	LogOutput     *reopen.FileWriter
	ExecDir       = ""
	// --- replace config text - do not remove ---
)

func init() {
	ExecDir, _ = osext.ExecutableFolder()

	InitializeLogging()
}

func getEnvOrDefault(envVar string, defEnvVar string) (newEnvVar string) {
	if newEnvVar = os.Getenv(envVar); len(newEnvVar) == 0 {
		return defEnvVar
	} else {
		return newEnvVar
	}
}

func InitializeLogging() {
	var err error
	if LogOutput == nil {
		LogOutput, err = reopen.NewFileWriter(LogPath)
		if err != nil {
			log.Fatalf("Log output file was not set: %s", err)
		}

		// set up log format
		logFormat := &log.JSONFormatter{}
		logFormat.TimestampFormat = time.RFC3339Nano

		log.SetOutput(LogOutput)
		log.SetFormatter(logFormat)
	}
}
