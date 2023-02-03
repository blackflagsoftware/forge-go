package sql

import (
	c "github.com/blackflagsoftware/forge-go/internal/column"
)

type (
	SqlEntity struct {
		Name         string
		Columns      []c.Column
		ColExistence c.ColumnExistence
	}
)
