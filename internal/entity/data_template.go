package entity

import (
	"fmt"
	"strings"

	con "github.com/blackflagsoftware/forge-go/internal/constant"
)

func (ep *Entity) BuildDataTemplate() {
	if ep.DynamicSchema {
		ep.DataTablePrefix = "fmt.Sprintf("
		ep.DataTable = fmt.Sprintf("%s.%s", ep.Schema, ep.Lower)
		ep.DataTablePostfix = fmt.Sprintf(", %s)", ep.DynamicSchemaPostfix)
	} else {
		ep.DataTable = ep.Lower
	}
	SqlGetColumns := ""
	foundOneKey := false
	foundOnePatch := false
	foundOnePost := false
	keys := ""
	patchKeys := ""
	keyCount := 1
	values := ""
	listOrder := ""
	postColumn := ""
	postColumnNames := ""
	patchColumn := ""
	keysCount := 1
	foundSerial := ""
	foundSerialDB := ""
	fileKey := []string{}
	fileGetColumn := []string{}
	filePostIncrInit := []string{}
	filePostIncrCheck := []string{}
	filePostIncr := []string{}
	for i, c := range ep.Columns {
		fileGetColumn = append(fileGetColumn, fmt.Sprintf("%s.%s = %sObj.%s", ep.Abbr, c.ColumnName.Camel, ep.Abbr, c.ColumnName.Camel))
		if c.DBType == "autoincrement" {
			foundSerialDB = c.ColumnName.Lower
			foundSerial = c.ColumnName.Camel
		}
		if c.DBType != "autoincrement" {
			if foundOnePost {
				postColumn += fmt.Sprintf(",\n\t\t\t%s", c.ColumnName.Lower)
				postColumnNames += fmt.Sprintf(",\n\t\t\t:%s", c.ColumnName.Lower)
			} else {
				postColumn += fmt.Sprintf("%s", c.ColumnName.Lower)
				postColumnNames += fmt.Sprintf(":%s", c.ColumnName.Lower)
				foundOnePost = true
			}
		}
		if i == 0 {
			SqlGetColumns += fmt.Sprintf("%s", c.ColumnName.Lower)
		} else {
			SqlGetColumns += fmt.Sprintf(",\n\t\t\t%s", c.ColumnName.Lower)
		}
		if c.PrimaryKey {
			if foundOneKey {
				keys += " and "
				values += ", "
				listOrder += ", "
			}
			foundOneKey = true
			if ep.SQLProvider == con.MYSQL {
				keys += fmt.Sprintf("%s = ?", c.ColumnName.Lower)
			} else {
				if c.DBType == "uuid" {
					keys += fmt.Sprintf("text(%s) = $%d", c.ColumnName.Lower, keyCount)
				} else {
					keys += fmt.Sprintf("%s = $%d", c.ColumnName.Lower, keyCount)
				}
			}
			patchKeys += fmt.Sprintf("%s = :%s", c.ColumnName.Lower, c.ColumnName.Lower)
			keyCount++
			values += fmt.Sprintf("%s.%s", ep.Name.Abbr, c.ColumnName.Camel)
			listOrder += fmt.Sprintf("%s", c.ColumnName.Lower)
			keysCount++
			fileKey = append(fileKey, fmt.Sprintf("%sObj.%s == %s.%s", ep.Abbr, c.ColumnName.Camel, ep.Abbr, c.ColumnName.Camel))
			if c.DBType == "autoincrement" || c.DBType == "int" {
				filePostIncrInit = append(filePostIncrInit, fmt.Sprintf("max%s := 0", c.ColumnName.Camel))
				filePostIncrCheck = append(filePostIncrCheck, fmt.Sprintf("\t\tif %sObj.%s > max%s {\n\t\t\tmax%s = %sObj.%s\n\t\t}", ep.Abbr, c.ColumnName.Camel, c.ColumnName.Camel, c.ColumnName.Camel, ep.Abbr, c.ColumnName.Camel))
				filePostIncr = append(filePostIncr, fmt.Sprintf("\t%s.%s = max%s + 1", ep.Abbr, c.ColumnName.Camel, c.ColumnName.Camel))
			}
		} else {
			if foundOnePatch {
				patchColumn += ",\n\t\t\t"
			}
			patchColumn += fmt.Sprintf("%s = :%s", c.ColumnName.Lower, c.ColumnName.Lower)
			foundOnePatch = true
		}
	}
	ep.SqlGetColumns = strings.TrimRight(SqlGetColumns, "\n")
	ep.SqlTableKeyKeys = keys
	ep.SqlTableKeyValues = values
	ep.SqlTableKeyListOrder = listOrder
	ep.SqlPostColumns = strings.TrimRight(postColumn, "\n")
	ep.SqlPostColumnsNamed = strings.TrimRight(postColumnNames, "\n")
	ep.SqlPatchColumns = strings.TrimRight(patchColumn, "\n")
	ep.SqlPatchWhere = patchKeys
	ep.SqlPostQuery = fmt.Sprintf(con.SQL_POST_QUERY, ep.Abbr)
	if foundSerial != "" {
		ep.SqlPostQuery = fmt.Sprintf(con.SQL_POST_QUERY_MYSQL, ep.Abbr)
		ep.SqlPostLastId = fmt.Sprintf(con.SQL_LAST_ID_MYSQL, ep.Camel, ep.Abbr, foundSerial)
		if ep.SQLProvider == con.POSTGRESQL {
			ep.SqlPostReturning = fmt.Sprintf(" returning %s", foundSerialDB)
			ep.SqlPostQuery = fmt.Sprintf(con.SQL_POST_QUERY_POSTGRES, ep.Abbr)
			ep.SqlPostLastId = fmt.Sprintf(con.SQL_LAST_ID_POSTGRES, ep.Abbr, foundSerial)
		}
	}
	ep.FileKeys = strings.Join(fileKey, " && ")
	ep.FileGetColumns = strings.Join(fileGetColumn, "\n\t\t\t")
	ep.FilePostIncr = fmt.Sprintf("%s\n\tfor _, %sObj := range %ss {\n%s\n\t}\n%s", strings.Join(filePostIncrInit, "\n"), ep.Abbr, ep.Abbr, strings.Join(filePostIncrCheck, "\n"), strings.Join(filePostIncr, "\n"))
}
