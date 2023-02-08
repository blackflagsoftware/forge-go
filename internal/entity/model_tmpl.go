package entity

import (
	"fmt"
	"strings"

	con "github.com/blackflagsoftware/forge-go/internal/constant"
	n "github.com/blackflagsoftware/forge-go/internal/name"
)

func (ep *Entity) BuildModelTemplate() {
	cArray := []string{}
	for _, c := range ep.Columns {
		tagFormat := n.BuildAltName(c.ColumnName.RawName, ep.TagFormat)
		cArray = append(cArray, fmt.Sprintf(con.MODEL_COLUMN_WO_GORM, c.ColumnName.Camel, c.GoType, c.ColumnName.Lower, tagFormat))
		if ep.ProjectFile.UseORM {
			cArray = append(cArray, fmt.Sprintf(con.MODEL_COLUMN_W_GORM, c.ColumnName.Camel, c.GoType, c.ColumnName.Lower, tagFormat, c.ColumnName.Lower))
		}
		if c.GoType == "*json.RawMessage" {
			ep.ModelIncludeJson = "\n\t\"encoding/json\"\n\n"
		}
	}
	ep.ModelRows = strings.Join(cArray, "\n")
	if ep.HaveNullColumns {
		ep.ModelIncludeNull = con.MODEL_INCLUDE_NULL
	}
	// TODO: what happens when they change the storage
	switch ep.ProjectFile.Storage {
	case "s":
		ep.InitStorage = "\tif config.StorageSQL {\n\t\treturn InitSQL()\n\t}"
	case "f":
		ep.InitStorage = fmt.Sprintf("\tif config.StorageFile {\n\t\treturn &File%s{}\n\t}", ep.Camel)
	case "m":
		ep.InitStorage = fmt.Sprintf("\tif config.StorageMongo {\n\t\treturn InitMongo()\n\t}")
	}
}
