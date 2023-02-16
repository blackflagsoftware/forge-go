package src

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

// All the migration work is done here
type (
	EngineAdapter interface {
		ConnectDB(Connection, bool) (*sqlx.DB, error) // the second argument is to tell the function to use admin user/pwd and connect to root DB
		CheckDB(*sqlx.DB, string) error
		CheckTable(*sqlx.DB) error
		LockTable(*sqlx.DB) bool
		UnlockTable(*sqlx.DB) error
	}

	EngineType string

	Connection struct {
		Engine         EngineType
		Host           string
		DB             string
		User           string
		Pwd            string
		AdminUser      string
		AdminPwd       string
		MigrationPath  string
		SkipInitialize bool
	}

	Migration struct {
		Id       int    `db:"id"`
		FileName string `db:"file_name"`
	}
)

const (
	POSTGRES EngineType = "postgres"
	MYSQL    EngineType = "mysql"
	SQLITE   EngineType = "sqlite3"
)

func StartMigration(c Connection) error {
	fmt.Println("Running: StartMigration")
	ea := DetermineEnginerAdapter(c.Engine)
	if err := InitializeDB(ea, c); err != nil {
		fmt.Println("StartMigration: InitializeDB error", err)
		return err
	}
	// brings back the db connector that will be used in the rest of the process
	db, err := InitializeMigrationTable(ea, c)
	if err != nil {
		fmt.Println("StartMigration: InitializeMigrationTable error:", err)
		return err
	}
	defer db.Close()
	defer ea.UnlockTable(db)
	// lock table
	for {
		if ea.LockTable(db) {
			break
		}
		fmt.Println("lock is on: waiting...")
		time.Sleep(500 * time.Millisecond)
	}
	migrationsMap, err := BuildMigrationMap(db)
	if err != nil {
		fmt.Println("StartMigration: BuildMigrationMap error:", err)
		return err
	}
	return ProcessFiles(db, migrationsMap, c.MigrationPath)
}

func DetermineEnginerAdapter(et EngineType) EngineAdapter {
	if et == POSTGRES {
		return &Postgres{}
	}
	if et == MYSQL {
		return &Mysql{}
	}
	if et == SQLITE {
		return &Sqlite{}
	}
	return &Mock{}
}

func InitializeDB(ea EngineAdapter, c Connection) error {
	// user admin* if need, but make sure the database is made
	// use interface to make sure to get the correct sql
	if !c.SkipInitialize {
		fmt.Println("Running: InitializDB")
		db, err := ea.ConnectDB(c, true)
		if err != nil {
			return err
		}
		defer db.Close()
		if err := ea.CheckDB(db, c.DB); err != nil {
			return err
		}
	}
	return nil
}

func InitializeMigrationTable(ea EngineAdapter, c Connection) (*sqlx.DB, error) {
	// make sure the table is created and the lock control record is created
	fmt.Println("Running: InitializMigrationTable")
	db, err := ea.ConnectDB(c, false)
	if err != nil {
		return db, err
	}
	if !c.SkipInitialize {
		err = ea.CheckTable(db)
	}
	return db, err
}

func BuildMigrationMap(db *sqlx.DB) (migrationsMap map[string]struct{}, err error) {
	// build map of all known files ran via the DB
	migrationsMap = make(map[string]struct{})
	fmt.Println("Running: BuildMigrtionMap")
	migrations := []Migration{}
	// get the known migrations per service
	sqlMigration := "SELECT id, file_name FROM migration ORDER BY id"
	err = db.Select(&migrations, sqlMigration)
	if err != nil {
		err = fmt.Errorf("BuildMigrationMap: unable to select fron migration; %s", err)
		return
	}
	// loop through and make these into a map; for easy lookup later
	for _, migration := range migrations {
		migrationsMap[migration.FileName] = struct{}{}
	}
	return
}

func ProcessFiles(db *sqlx.DB, migrationsMap map[string]struct{}, migrationPath string) error {
	// read the directory
	fmt.Println("Running: ProcessFiles")
	files, err := os.ReadDir(migrationPath)
	if err != nil {
		return fmt.Errorf("ProcessingFiles: unable to read directory; %s", err)
	}
	for _, f := range files {
		// make sure the file is valid format
		if ValidFile(f) {
			// compare with migrationsMap if there is anything new, run the script/exec
			nameParts := strings.Split(f.Name(), ".")
			justFileName := nameParts[0]
			if _, ok := migrationsMap[justFileName]; !ok {
				// have not been seen before, let's continue
				filePath := filepath.Join(migrationPath, f.Name())
				if err := ProcessMigration(db, filePath, justFileName); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func ValidFile(file fs.DirEntry) bool {
	if file.IsDir() {
		// skip directories
		// fmt.Printf("skipping directory: %s\n", file.Name())
		return false
	}
	if matched, err := regexp.Match(`^[0-9]{14}[a-zA-Z-]*\.(sql|bin)$`, []byte(file.Name())); !matched || err != nil {
		// skip non normalized names or error in regex
		fmt.Printf("skipping non normalized name: %s\n", file.Name())
		if err != nil {
			fmt.Printf("\twith error: %s\n", err)
		}
		return false
	}
	return true
}

func ProcessMigration(db *sqlx.DB, migrationPath, fileName string) error {
	fmt.Printf("Running: ProcessMigration for [%s]\n", migrationPath)
	ext := path.Ext(migrationPath)
	switch ext {
	case ".sql":
		fileContent, err := os.ReadFile(migrationPath)
		if err != nil {
			return fmt.Errorf("ProcessMigration: unable to read file: %s; %s\n", migrationPath, err)
		}
		_, errExec := db.Exec(string(fileContent))
		if errExec != nil {
			return fmt.Errorf("ProcessMigration: unable to run file query: %s; %s\n", migrationPath, errExec)
		}
	case ".bin":
		cmd := exec.Command(migrationPath)
		stdOut, errCmd := cmd.Output()
		if errCmd != nil {
			return fmt.Errorf("ProcessMigration: unable to run file binary: %s; %s\n\tOutput: %s", migrationPath, errCmd, stdOut)
		}
	default:
		fmt.Println("ProcessMigration: not the correct extention:", ext)
		return nil
	}
	migration := Migration{FileName: fileName}
	sqlInsert := "INSERT INTO migration (file_name) VALUES (:file_name)"
	_, errInsert := db.NamedExec(sqlInsert, migration)
	if errInsert != nil {
		return fmt.Errorf("ProcessMigration: unable to insert into migration; %s", errInsert)
	}
	return nil
}
