package src

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
)

type (
	Sqlite struct{}
)

func (s *Sqlite) ConnectDB(c Connection, rootDB bool) (*sqlx.DB, error) {
	var db *sqlx.DB
	conn := fmt.Sprintf("%s?cache=shared&mode=wrc", c.Host)
	db, err := sqlx.Open("sqlite3", conn)
	if err != nil {
		fmt.Println("Could not connect with connection string:", conn)
		os.Exit(1)
	}
	db.SetMaxOpenConns(1)
	return db, nil
}

func (s *Sqlite) CheckDB(db *sqlx.DB, dbName string) error {
	return nil
}

func (s *Sqlite) CheckTable(db *sqlx.DB) error {
	return nil
}

// the nature of sqlite is not to be distributed, no need to lock/unlock
func (s *Sqlite) LockTable(db *sqlx.DB) bool {
	return true
}

func (s *Sqlite) UnlockTable(db *sqlx.DB) error {
	return nil
}
