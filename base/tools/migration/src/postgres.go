package src

import (
	"fmt"

	"github.com/blackflagsoftware/forge-go/base/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type (
	Postgres struct{}
)

func (p *Postgres) ConnectDB(c Connection, rootDB bool) (*sqlx.DB, error) {
	var db *sqlx.DB
	dbName := c.DB
	user := c.User
	pwd := c.Pwd
	if rootDB {
		dbName = "postgres"
		if c.AdminUser != "" {
			user = c.AdminUser
		}
		if c.AdminPwd != "" {
			pwd = c.AdminPwd
		}
	}
	conn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable", user, pwd, dbName, c.Host)
	if pwd == "" {
		conn = fmt.Sprintf("user=%s dbname=%s host=%s sslmode=disable", user, dbName, c.Host)
	}
	db, errOpen := sqlx.Open("postgres", conn)
	if errOpen != nil {
		return db, fmt.Errorf("ConnectDB[postgres]: unable to open DB %s****; %s", conn[:6], errOpen)
	}
	return db, nil
}

func (p *Postgres) CheckDB(db *sqlx.DB, dbName string) error {
	checkSql := "SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE lower(datname) = lower($1))"
	exists := false
	err := db.Get(&exists, checkSql, dbName)
	if err != nil {
		return fmt.Errorf("CheckDB[postgres]: unable to check for existing database; %s", err)
	}
	if !exists {
		createSql := fmt.Sprintf("CREATE DATABASE %s", dbName)
		if _, err := db.Exec(createSql); err != nil {
			return fmt.Errorf("CheckDB[postgres]: unable to create database; %s", err)
		}
	}
	return nil
}

func (p *Postgres) CheckTable(db *sqlx.DB) error {
	checkSql := "SELECT EXISTS(SELECT table_name FROM information_schema.tables WHERE table_name = 'migration')"
	exists := false
	err := db.Get(&exists, checkSql)
	if err != nil {
		return fmt.Errorf("CheckTable[postgres]: unable to check for existing table; %s", err)
	}
	if !exists {
		createSql := "CREATE TABLE migration (id serial, file_name varchar(100) NOT NULL)"
		if _, err := db.Exec(createSql); err != nil {
			return fmt.Errorf("CheckTable[postgres]: unable to create table; %s", err)
		}
	}
	return nil
}

func (p *Postgres) LockTable(db *sqlx.DB) bool {
	number := config.GetUniqueNumberForLock()
	succeed := false
	lockSql := "SELECT pg_try_advisory_lock($1)"
	if err := db.Get(&succeed, lockSql, number); err != nil {
		fmt.Printf("LockTable[postgres]: unable to lock resource; %s", err)
		return false
	}
	return succeed
}

func (p *Postgres) UnlockTable(db *sqlx.DB) error {
	number := config.GetUniqueNumberForLock()
	succeed := false
	lockSql := "SELECT pg_advisory_unlock($1)"
	if err := db.Get(&succeed, lockSql, number); err != nil {
		return fmt.Errorf("UnlockTable[postgres]: unable to unlock; %s", err)
	}
	if !succeed {
		return fmt.Errorf("UnlockTable[postgres]: unable to unlock with no error")
	}
	return nil
}
