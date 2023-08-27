package entity

import (
	"strings"

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
		RestStrConv          string
		RestGetDeleteUrl     string
		RestGetDeleteAssign  string
		RestArgSet           string
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
		SQLProvider          string // optional if using SQL as a storage, either Psql, MySql or Sqlite; this interfaces with sqlx
		SQLProviderLower     string // optional is using SQL, lowercase of above
		GrpcTranslateIn      string
		GrpcTranslateOut     string
		GrpcImport           string
		GrpcArgsInit         string
		DefaultColumn        string
		SortColumns          string
		n.Name
		ModelTemplateItems
		ManagerTemplateItems
		pf.ProjectFile
		c.ColumnExistence
	}

	ModelTemplateItems struct {
		ModelImport      string
		ModelRows        string
		ModelInitStorage string
	}

	ManagerTemplateItems struct {
		ManagerImport        string
		ManagerGetRows       string
		ManagerPostRows      string
		ManagerPatchInitArgs string
		ManagerPatchRows     string
		ManagerTestImport    string
		ManagerGetTestRow    string
		ManagerPostTestRow   string
		ManagerPatchTestInit string
		ManagerDeleteTestRow string
		ManagerAuditKey      string
	}

	PostPutTest struct {
		Name         string
		ForColumn    string
		ColumnLength int
		Failure      bool
	}

	ColumnTest struct {
		GoType string
		DBType string
	}
)

func (e *Entity) HasNullColumn() bool {
	for _, c := range e.Columns {
		if strings.Contains(c.GoType, "null") {
			return true
		}
	}
	return false
}

func (e *Entity) HasTimeColumn() bool {
	for _, c := range e.Columns {
		if c.GoType == "time.Time" || c.ColumnName.Lower == "created_at" || c.ColumnName.Lower == "updated_at" {
			return true
		}
	}
	return false
}

func (e *Entity) HasNullTimeColumn() bool {
	for _, c := range e.Columns {
		if c.GoType == "null.Time" {
			return true
		}
	}
	return false
}

func (e *Entity) HasJsonColumn() bool {
	for _, c := range e.Columns {
		if c.GoType == "*json.RawMessage" {
			return true
		}
	}
	return false
}

func (e *Entity) HasPrimaryUUIDColumn() bool {
	for _, c := range e.Columns {
		if c.DBType == "uuid" && c.PrimaryKey {
			return true
		}
	}
	return false
}
