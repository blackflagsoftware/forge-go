package entity

import (
	"fmt"
	"strings"

	n "github.com/blackflagsoftware/forge-go/internal/name"
)

const (
	COLUMN_W_GORM  = "\t\t%s\t%s\t`db:\"%s\" json:\"%s\" gorm:\"column:%s\"`"
	COLUMN_WO_GORM = "\t\t%s\t%s\t`db:\"%s\" json:\"%s\"`"
)

func (ep *Entity) BuildModelTemplate() {
	cArray := []string{}
	for _, c := range ep.Columns {
		tagFormat := n.BuildAltName(c.ColumnName.RawName, ep.TagFormat)
		colType := c.GoType
		if c.PrimaryKey {
			colType = c.GoTypeNonSql
		}
		cArray = append(cArray, fmt.Sprintf(COLUMN_WO_GORM, c.ColumnName.Camel, colType, c.ColumnName.Lower, tagFormat))
		if ep.ProjectFile.UseORM {
			cArray = append(cArray, fmt.Sprintf(COLUMN_W_GORM, c.ColumnName.Camel, colType, c.ColumnName.Lower, tagFormat, c.ColumnName.Lower))
		}
	}
	ep.ModelRows = strings.Join(cArray, "\n")
	// TODO: what happens when they change the storage
	switch ep.ProjectFile.Storage {
	case "s":
		ep.ModelInitStorage = "\tif config.StorageSQL {\n\t\treturn InitSQL()\n\t}"
	case "f":
		ep.ModelInitStorage = fmt.Sprintf("\tif config.StorageFile {\n\t\treturn &File%s{}\n\t}", ep.Camel)
	case "m":
		ep.ModelInitStorage = fmt.Sprintf("\tif config.StorageMongo {\n\t\treturn InitMongo()\n\t}")
	}
	// build import
	importLines := []string{}
	if ep.HasJsonColumn() {
		importLines = append(importLines, "\"encoding/json\"\n")
	}
	importLines = append(importLines, fmt.Sprintf("\"%s/config\"", ep.ProjectPath))
	importLines = append(importLines, fmt.Sprintf("\"%s/internal/util\"", ep.ProjectPath))
	if ep.HasNullColumn() {
		importLines = append(importLines, "\"gopkg.in/guregu/null.v3\"")
	}
	ep.ModelImport = strings.Join(importLines, "\n\t")
}
