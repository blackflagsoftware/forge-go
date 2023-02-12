package entity

import (
	c "github.com/blackflagsoftware/forge-go/internal/column"
	n "github.com/blackflagsoftware/forge-go/internal/name"
	pf "github.com/blackflagsoftware/forge-go/internal/projectfile"
)

var (
	PostTests map[string]ColumnTest
	PutTests  map[string]ColumnTest
)

type (
	Entity struct {
		Columns              []c.Column
		ModelIncludeNull     string
		ModelRows            string
		ModelIncludeJson     string
		RestStrConv          string
		RestGetDeleteUrl     string
		RestGetDeleteAssign  string
		RestArgSet           string
		ManagerTime          string
		ManagerGetRow        string
		ManagerPostRows      string
		ManagerPutRows       string
		ManagerPatchInitArgs string
		ManagerPatchTestInit string
		ManagerPatchRows     string
		ManagerGetTestRow    string
		ManagerPostTestRow   string
		ManagerPutTestRow    string
		ManagerDeleteTestRow string
		ManagerUtilPath      string
		ManagerImportTest    string
		DataTable            string
		DataTablePrefix      string
		DataTablePostfix     string
		SqlGetColumns        string
		SqlTableKeyKeys      string
		SqlTableKeyValues    string
		SqlTableKeyListOrder string
		SqlPostColumns       string
		SqlPostColumnsNamed  string
		SqlPostReturning     string
		SqlPostQuery         string
		SqlPostLastId        string
		SqlPatchColumns      string
		SqlPatchWhere        string
		SqlPatchWhereValues  string
		FileKeys             string
		FileGetColumns       string
		FilePostIncr         string
		SqlLines             []string
		InitStorage          string // holds the formatted lines for InitStorage for model
		SQLProvider          string // optional if using SQL as a storage, either Psql, MySql or Sqlite; this interfaces with sqlx
		GrpcTranslateIn      string
		GrpcTranslateOut     string
		GrpcImport           string
		GrpcArgsInit         string
		DefaultColumn        string
		SortColumns          string
		n.Name
		pf.ProjectFile
		c.ColumnExistence
	}

	PostPutTest struct {
		Name         string
		ForColumn    string
		ColumnLength int
		Failure      bool
	}

	ColumnTest struct {
		GoType string
	}
)
