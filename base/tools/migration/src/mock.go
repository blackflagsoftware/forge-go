package src

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type (
	Mock struct{}
)

func (m *Mock) ConnectDB(c Connection, rootDB bool) (*sqlx.DB, error) {
	fmt.Println("Warning: using mock...")
	var db *sqlx.DB
	return db, nil
}

func (m *Mock) CheckUser(db *sqlx.DB, c Connection) error {
	fmt.Println("Warning: using mock...")
	return nil
}

func (m *Mock) CheckDB(db *sqlx.DB, c Connection) error {
	fmt.Println("Warning: using mock...")
	return nil
}

func (m *Mock) CheckTable(db *sqlx.DB) error {
	fmt.Println("Warning: using mock...")
	return nil
}

func (m *Mock) LockTable(db *sqlx.DB) bool {
	fmt.Println("Warning: using mock...")
	return true
}

func (m *Mock) UnlockTable(db *sqlx.DB) error {
	fmt.Println("Warning: using mock...")
	return nil
}
