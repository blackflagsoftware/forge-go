package task

import (
	"fmt"
	"strings"

	m "github.com/blackflagsoftware/forge-go/internal/model"
)

const (
	COLUMN_W_GORM  = "\t\t%s\t%s\t`db:\"%s\" json:\"%s\" gorm:\"column:%s\"`"
	COLUMN_WO_GORM = "\t\t%s\t%s\t`db:\"%s\" json:\"%s\"`"
)

func buildModelTemplate(p *m.Project) {
	cArray := []string{}
	for _, c := range p.CurrentEntity.Columns {
		tagFormat := m.BuildAltName(c.ColumnName.RawName, p.TagFormat)
		colType := c.GoType
		if c.PrimaryKey {
			colType = c.GoTypeNonSql
		}
		cArray = append(cArray, fmt.Sprintf(COLUMN_WO_GORM, c.ColumnName.Camel, colType, c.ColumnName.Lower, tagFormat))
		if p.ProjectFile.UseORM {
			cArray = append(cArray, fmt.Sprintf(COLUMN_W_GORM, c.ColumnName.Camel, colType, c.ColumnName.Lower, tagFormat, c.ColumnName.Lower))
		}
	}
	p.ModelRows = strings.Join(cArray, "\n")
	// TODO: what happens when they change the storage
	switch p.ProjectFile.Storage {
	case "s":
		p.ModelInitStorage = "\tif config.StorageSQL {\n\t\treturn InitSQL()\n\t}"
	case "f":
		p.ModelInitStorage = fmt.Sprintf("\tif config.StorageFile {\n\t\treturn &File%s{}\n\t}", p.Camel)
	case "m":
		p.ModelInitStorage = fmt.Sprintf("\tif config.StorageMongo {\n\t\treturn InitMongo()\n\t}")
	}
	// build import
	importLines := []string{}
	if p.CurrentEntity.HasJsonColumn() {
		importLines = append(importLines, "\"encoding/json\"\n")
	}
	importLines = append(importLines, fmt.Sprintf("\"%s/config\"", p.ProjectPath))
	importLines = append(importLines, fmt.Sprintf("\"%s/internal/util\"", p.ProjectPath))
	if p.CurrentEntity.HasNullColumn() {
		importLines = append(importLines, "\"gopkg.in/guregu/null.v3\"")
	}
	p.ModelImport = strings.Join(importLines, "\n\t")
}
