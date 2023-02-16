package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/blackflagsoftware/forge-go/base/config"
	mig "github.com/blackflagsoftware/forge-go/base/tools/migration/src"
)

func main() {
	migration, file, c := InitializeConnection()
	err := os.MkdirAll(c.MigrationPath, 0744)
	if err != nil {
		fmt.Printf("Unable to make scripts/migrations directory structure: %s with error: %s\n", c.MigrationPath, err)
	}

	if migration {
		if c.Host == "" || c.User == "" {
			fmt.Println("Missing host or user values")
			os.Exit(1)
		}
		if err := mig.StartMigration(c); err != nil {
			fmt.Println(err)
		}
		os.Exit(0)
	}

	if file != "" {
		processFile(file, c.MigrationPath)
		os.Exit(0)
	}

	interactive(c.MigrationPath)
}

func InitializeConnection() (bool, string, mig.Connection) {
	var migration bool
	var file string
	var host string
	var db string
	var user string
	var pwd string
	var adminUser string
	var adminPwd string
	var path string
	var engine string
	var skipInit bool

	flag.BoolVar(&migration, "m", false, "Run migration for the project")
	flag.StringVar(&file, "f", "", "File to normalize name; name of the file (only) found in <project root>/scripts/migrations")
	flag.StringVar(&host, "h", "", "Host name")
	flag.StringVar(&db, "d", "", "DB name")
	flag.StringVar(&user, "u", "", "Username")
	flag.StringVar(&pwd, "p", "", "Password")
	flag.StringVar(&adminUser, "au", "", "Admin Username")
	flag.StringVar(&adminPwd, "ap", "", "Admin Password")
	flag.StringVar(&path, "t", "", "Scripts path; typically the path to <project root>/scripts/migrations")
	flag.StringVar(&engine, "e", "", "DB engine [postgres | mysql | sqlite3]")
	flag.BoolVar(&skipInit, "s", false, "Skip db/table initialization process")

	flag.Parse()

	if path == "" {
		path = filepath.Clean(config.ExecDir + "../../../scripts/migrations")
	}
	dbHost := config.GetEnvOrDefault("FORGE_GO_BASE_DB_HOST", "")
	if host != "" {
		dbHost = host
	}
	dbDB := config.GetEnvOrDefault("FORGE_GO_BASE_DB_DB", "")
	if db != "" {
		dbDB = db
	}
	dbUser := config.GetEnvOrDefault("FORGE_GO_BASE_DB_USER", "")
	if user != "" {
		dbUser = user
	}
	dbPass := config.GetEnvOrDefault("FORGE_GO_BASE_DB_PASS", "")
	if pwd != "" {
		dbPass = pwd
	}
	adminDBUser := config.GetEnvOrDefault("FORGE_GO_BASE_ADMIN_DB_USER", "")
	if adminUser != "" {
		adminDBUser = user
	}
	adminDBPass := config.GetEnvOrDefault("FORGE_GO_BASE_ADMIN_DB_PASS", "")
	if adminPwd != "" {
		adminDBPass = pwd
	}
	dbEngine := config.GetEnvOrDefault("FORGE_GO_BASE_MIGRATION_DB_ENGINE", "")
	if engine != "" {
		dbEngine = engine
	}
	return migration, file, mig.Connection{
		Host:           dbHost,
		DB:             dbDB,
		User:           dbUser,
		Pwd:            dbPass,
		AdminUser:      adminDBUser,
		AdminPwd:       adminDBPass,
		MigrationPath:  path,
		SkipInitialize: skipInit,
		Engine:         mig.EngineType(dbEngine),
	}
}

func processFile(fileName, migrationDir string) {
	filePath := fmt.Sprintf("%s/%s", migrationDir, fileName)
	now := time.Now().Format("20060102150405")
	ext := ".sql"
	if path.Ext(fileName) == ".sql" {
		ext = ""
	}
	renameFileName := fmt.Sprintf("%s/%s-%s%s", migrationDir, now, normalizeName(fileName), ext)
	err := os.Rename(filePath, renameFileName)
	if err != nil {
		fmt.Printf("Error in renaming file: %s\n", err)
	}
}

func interactive(migrationDir string) {
	for {
		fmt.Printf("Paste or enter your sql code below; type 'exit' when done\n\n")
		lines := []string{}
		reader := bufio.NewReader(os.Stdin)
		for {
			line := parseInput(reader)
			if line == "exit" {
				break
			}
			lines = append(lines, line)
		}

		fmt.Print("Enter a description of the migration (will cut it off at 85 characters): ")
		name := parseInput(reader)
		now := time.Now().Format("20060102150405")
		fileName := fmt.Sprintf("%s/%s-%s.sql", migrationDir, now, normalizeName(name))
		err := os.WriteFile(fileName, []byte(strings.Join(lines, "\n")), 0644)
		if err != nil {
			fmt.Printf("Unable to save file: %s\n", err)
		}
		fmt.Println("")
		fmt.Print("Add another (y/n): ")
		another := parseInput(reader)
		another = strings.ToLower(another)
		if another != "y" {
			break
		}
	}
}

func parseInput(reader *bufio.Reader) string {
	s, _ := reader.ReadString('\n')
	s = strings.TrimSpace(s)
	return s
}

// Let's make sure all names are normalized and that the length does not exceed 85 due to the column max size and adding the date in front of it
func normalizeName(fileName string) (normalizedName string) {
	normalizedName = strings.ReplaceAll(fileName, " ", "-")
	normalizedName = strings.ReplaceAll(normalizedName, "_", "-")
	if len(normalizedName) > 85 {
		normalizedName = normalizedName[:85]
	}
	return
}
