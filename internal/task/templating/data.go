package task

import (
	"fmt"
	"strings"

	con "github.com/blackflagsoftware/forge-go/internal/constant"
	m "github.com/blackflagsoftware/forge-go/internal/model"
)

func buildStorageTemplate(p *m.Project) {
	if p.ProjectFile.DynamicSchema {
		p.StorageTablePrefix = "fmt.Sprintf("
		p.StorageTable = fmt.Sprintf("%s.%s", p.Schema, p.CurrentEntity.Lower)
		p.StorageTablePostfix = fmt.Sprintf(", %s)", p.DynamicSchemaPostfix)
	} else {
		p.StorageTable = p.CurrentEntity.Lower
	}
	SqlGetColumns := ""
	foundOneKey := false
	foundOnePatch := false
	foundOnePost := false
	keys := ""
	patchKeys := ""
	keyCount := 1
	values := ""
	// listOrder := ""
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
	countNeeded := false
	countColLower := ""
	countColCamel := ""
	countFuncNeeded := ""
	for i, c := range p.CurrentEntity.Columns {
		fileGetColumn = append(fileGetColumn, fmt.Sprintf("%s.%s = %sObj.%s", p.CurrentEntity.Abbr, c.ColumnName.Camel, p.CurrentEntity.Abbr, c.ColumnName.Camel))
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
				// listOrder += ", "
			}
			foundOneKey = true
			if !countNeeded {
				// only set for int key that are type mongo or sqlite3
				// this will provide a "auto increment" type functionality
				if c.DBType == "int" {
					includeCount := false
					if p.Storage == "m" {
						countFuncNeeded = "MONGO"
						includeCount = true
					}
					if p.Storage == "s" && p.SQLProvider == con.SQLITE3 {
						countFuncNeeded = "SQL"
						includeCount = true
					}
					if includeCount {
						countNeeded = true
						countColLower = c.ColumnName.Lower
						countColCamel = c.ColumnName.Camel
					}
				}
			}
			if p.SQLProvider == con.MYSQL {
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
			if c.DBType == "uuid" {
				if p.StorageImport == "" {
					p.StorageImport = "\"strings\""
				}
				values += fmt.Sprintf("strings.ToLower(%s.%s)", p.CurrentEntity.Abbr, c.ColumnName.Camel)
			} else {
				values += fmt.Sprintf("%s.%s", p.CurrentEntity.Abbr, c.ColumnName.Camel)
			}
			// listOrder += fmt.Sprintf("%s", c.ColumnName.Lower)
			keysCount++
			fileKey = append(fileKey, fmt.Sprintf("%sObj.%s == %s.%s", p.CurrentEntity.Abbr, c.ColumnName.Camel, p.CurrentEntity.Abbr, c.ColumnName.Camel))
			if c.DBType == "autoincrement" || c.DBType == "int" {
				filePostIncrInit = append(filePostIncrInit, fmt.Sprintf("max%s := 0", c.ColumnName.Camel))
				filePostIncrCheck = append(filePostIncrCheck, fmt.Sprintf("\t\tif %sObj.%s > max%s {\n\t\t\tmax%s = %sObj.%s\n\t\t}", p.CurrentEntity.Abbr, c.ColumnName.Camel, c.ColumnName.Camel, c.ColumnName.Camel, p.CurrentEntity.Abbr, c.ColumnName.Camel))
				filePostIncr = append(filePostIncr, fmt.Sprintf("\t%s.%s = max%s + 1", p.CurrentEntity.Abbr, c.ColumnName.Camel, c.ColumnName.Camel))
			}
		} else {
			if foundOnePatch {
				patchColumn += ",\n\t\t\t"
			}
			patchColumn += fmt.Sprintf("%s = :%s", c.ColumnName.Lower, c.ColumnName.Lower)
			foundOnePatch = true
		}
	}
	p.StorageGetColumns = strings.TrimRight(SqlGetColumns, "\n")
	p.StorageTableKeyKeys = keys
	p.StorageTableKeyValues = values
	// p.StorageTableKeyListOrder = listOrder
	p.StoragePostColumns = strings.TrimRight(postColumn, "\n")
	p.StoragePostColumnsNamed = strings.TrimRight(postColumnNames, "\n")
	p.StoragePatchColumns = strings.TrimRight(patchColumn, "\n")
	p.StoragePatchWhere = patchKeys
	p.StoragePostQuery = fmt.Sprintf(con.SQL_POST_QUERY, p.CurrentEntity.Abbr)
	if foundSerial != "" {
		p.StoragePostQuery = fmt.Sprintf(con.SQL_POST_QUERY_MYSQL, p.CurrentEntity.Abbr)
		p.StoragePostLastId = fmt.Sprintf(con.SQL_LAST_ID_MYSQL, p.Camel, p.CurrentEntity.Abbr, foundSerial)
		if p.SQLProvider == con.POSTGRESQL {
			p.StoragePostReturning = fmt.Sprintf(" returning %s", foundSerialDB)
			p.StoragePostQuery = fmt.Sprintf(con.SQL_POST_QUERY_POSTGRES, p.CurrentEntity.Abbr)
			p.StoragePostLastId = fmt.Sprintf(con.SQL_LAST_ID_POSTGRES, p.CurrentEntity.Abbr, foundSerial)
		}
	}
	if countNeeded {
		p.StorageCountCall = fmt.Sprintf(con.STORAGE_COUNT_CALL, p.CurrentEntity.Abbr, countColCamel)
		p.StorageCountFunc = fmt.Sprintf(con.STORAGE_COUNT_FUNC_SQL, p.CurrentEntity.Camel, countColLower, p.CurrentEntity.Camel)
		if countFuncNeeded == "MONGO" {
			p.StorageCountFunc = fmt.Sprintf(con.STORAGE_COUNT_FUNC_MONGO, p.CurrentEntity.Camel, p.CurrentEntity.Lower, p.CurrentEntity.Abbr, countColLower, p.CurrentEntity.Camel, p.CurrentEntity.Abbr, p.CurrentEntity.Camel, p.CurrentEntity.Abbr, p.CurrentEntity.Camel, p.CurrentEntity.Abbr, countColCamel)
			p.StorageCountImport = "\n\t\"go.mongodb.org/mongo-driver/mongo/options\""
		}
	}
	p.FileKeys = strings.Join(fileKey, " && ")
	p.FileGetColumns = strings.Join(fileGetColumn, "\n\t\t\t")
	p.FilePostIncr = fmt.Sprintf("%s\n\tfor _, %sObj := range %ss {\n%s\n\t}\n%s", strings.Join(filePostIncrInit, "\n"), p.CurrentEntity.Abbr, p.CurrentEntity.Abbr, strings.Join(filePostIncrCheck, "\n"), strings.Join(filePostIncr, "\n"))
}
