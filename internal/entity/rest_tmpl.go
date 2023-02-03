package entity

import (
	"fmt"

	con "github.com/blackflagsoftware/forge-go/internal/constant"
)

func (ep *Entity) BuildRestTemplate() {
	// build get/delete url
	getDeleteUrl := ""
	foundOne := false
	for _, c := range ep.Columns {
		if c.PrimaryKey {
			if foundOne {
				getDeleteUrl += fmt.Sprintf("/%s/:%s", c.ColumnName.Lower, c.ColumnName.Lower)
			} else {
				getDeleteUrl = fmt.Sprintf(":%s", c.ColumnName.Lower)
				foundOne = true
			}
		}
	}
	ep.RestGetDeleteUrl = getDeleteUrl
	// build get/delete assign and args
	getDeleteAssign := ""
	setArgs := ""
	foundOne = false
	for _, c := range ep.Columns {
		if c.PrimaryKey {
			if c.GoType == "string" {
				if foundOne {
					getDeleteAssign += "\n"
					setArgs += ", "
				}
				getDeleteAssign += fmt.Sprintf(con.REST_PRIMARY_STR, c.ColumnName.Lower, c.ColumnName.Lower)
				setArgs += fmt.Sprintf("%s: %s", c.ColumnName.Camel, c.ColumnName.Lower)
				foundOne = true
			}
			if c.GoType == "int" {
				if foundOne {
					getDeleteAssign += "\n"
					setArgs += ", "
				}
				getDeleteAssign += fmt.Sprintf(con.REST_PRIMARY_INT, c.ColumnName.Lower, c.ColumnName.Lower, c.ColumnName.Lower, c.ColumnName.Lower)
				setArgs += fmt.Sprintf("%s: int(%s)", c.ColumnName.Camel, c.ColumnName.Lower)
				foundOne = true
				ep.RestStrConv = "\n\t\"strconv\""
			}
		}
	}
	ep.RestGetDeleteAssign = getDeleteAssign
	ep.RestArgSet = setArgs
}
