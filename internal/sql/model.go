package sql

import (
	m "github.com/blackflagsoftware/forge-go/internal/model"
)

type (
	SqlEntity struct {
		Name         string
		Columns      []m.Column
		ColExistence m.ColumnExistence
	}
)
