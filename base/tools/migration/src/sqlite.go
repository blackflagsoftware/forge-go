package src

import (
	"fmt"

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
		return db, fmt.Errorf("Could not connect with connection string: %s", conn)
	}
	db.SetMaxOpenConns(1)
	return db, nil
}

func (s *Sqlite) CheckUser(db *sqlx.DB, c Connection) error {
	return nil
}

func (s *Sqlite) CheckDB(db *sqlx.DB, c Connection) error {
	return nil
}

func (s *Sqlite) CheckTable(db *sqlx.DB) error {
	createSql := "CREATE TABLE IF NOT EXISTS migration (id integer primary key autoincrement, file_name varchar(100) NOT NULL)"
	if _, err := db.Exec(createSql); err != nil {
		return fmt.Errorf("CheckTable[sqlite]: unable to create table; %s", err)
	}
	return nil
}

// the nature of sqlite is not to be distributed, no need to lock/unlock
func (s *Sqlite) LockTable(db *sqlx.DB) bool {
	return true
}

func (s *Sqlite) UnlockTable(db *sqlx.DB) error {
	return nil
}
