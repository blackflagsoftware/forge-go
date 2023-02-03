package column

import (
	n "github.com/blackflagsoftware/forge-go/internal/name"
)

type (
	Column struct {
		ColumnName   n.Name
		DBType       string
		GoType       string
		GoTypeNonSql string
		Null         bool
		DefaultValue string
		Length       int64
		PrimaryKey   bool
	}

	ColumnExistence struct {
		HaveNullColumns bool
		TimeColumn      bool
	}
)
