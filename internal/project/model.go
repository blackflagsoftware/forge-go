package project

import (
	e "github.com/blackflagsoftware/forge-go/internal/entity"
	pf "github.com/blackflagsoftware/forge-go/internal/projectfile"
)

type (
	Project struct {
		UseBlank    bool
		Entities    []e.Entity
		ProjectFile pf.ProjectFile
	}

	StorageVars struct {
		ProjectPath           string
		SQLProvider           string // optional if using SQL as a storage, either Psql, MySql or Sqlite; this interfaces with sqlx
		SQLProviderLower      string // optional if using SQL as a storage, either psql, mysql or sqlite; this interfaces with gorm
		SQLProviderConnection string // holds the connection string for gorm of the other sql types
	}

	MigrationVars struct {
		ProjectPath         string
		MigrationVerify     string
		MigrationConnection string
		MigrationHeader     string
	}
)
