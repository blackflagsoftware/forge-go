package sql

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/tabwriter"

	c "github.com/blackflagsoftware/forge-go/internal/column"
	"github.com/blackflagsoftware/forge-go/internal/util"
)

func ParseSqlLines(lines []string) (sqlEntity SqlEntity) {
	line := util.FormatSql(lines)
	tableName, columnPart, err := tableNameParse(line)
	if err != nil {
		fmt.Println(err)
		return
	}
	sqlEntity.Name = tableName
	columns, err := columnsParse(columnPart)
	if err != nil {
		fmt.Println(err)
		return
	}
	sqlEntity.Columns = columns
	sqlEntity.ColExistence = determineColExists(columns)
	return
}

func tableNameParse(line string) (tableName, columnPart string, err error) {
	reg := regexp.MustCompile(`(?i)create table (if )?(not )?(exists )?(?P<table_name>\w+) \((?P<column_part>.+)\);*`)
	matches := reg.FindStringSubmatch(line)
	if len(matches) < 2 {
		err = fmt.Errorf("tableParse: sql parse invalid")
		return
	}
	tableNameIdx := reg.SubexpIndex("table_name")
	columnPartIdx := reg.SubexpIndex("column_part")
	tableName = matches[tableNameIdx]
	columnPart = matches[columnPartIdx]
	return
}

func columnsParse(columnPart string) (columns []c.Column, err error) {
	reg := regexp.MustCompile(`(?P<field_name>[a-zA-Z_]+) (?P<field_type>[a-zA-Z0-9\(\),]+)(?P<the_rest>.+)?`)
	keys := []string{}
	// split columns and go through each column
	lines := strings.Split(columnPart, ",")
	for i := range lines {
		line := strings.ReplaceAll(strings.TrimSpace(strings.ToLower(lines[i])), "`", "")
		if strings.Index(line, "primary") == 0 {
			// parse primary key () line
			keys = parsePrimaryKey(line)
			continue
		}
		matches := reg.FindStringSubmatch(line)
		if len(matches) < 2 {
			err = fmt.Errorf("columnsParse: sql parse invalid")
			return
		}
		fieldNameIdx := reg.SubexpIndex("field_name")
		fieldTypeIdx := reg.SubexpIndex("field_type")
		theRestIdx := reg.SubexpIndex("the_rest")
		col := c.Column{}
		rawName := matches[fieldNameIdx]
		col.ColumnName.BuildName(rawName, []string{})
		col.DBType = matches[fieldTypeIdx]
		theRest := strings.TrimSpace(matches[theRestIdx])
		theRestParse(&col, theRest)
		setGoType(&col)
		columns = append(columns, col)
	}
	markPrimary(&columns, keys)
	return
}

func parsePrimaryKey(line string) (keyNames []string) {
	reg := regexp.MustCompile(`primary key ?\((?P<keys>.+)\)`)

	matches := reg.FindStringSubmatch(line)
	if len(matches) < 1 {
		fmt.Println("parsePrimaryKey: sql parse invalid")
		return
	}
	keysIdx := reg.SubexpIndex("keys")
	keys := matches[keysIdx]
	split := strings.Split(keys, ",")
	for i := range split {
		keyNames = append(keyNames, strings.TrimSpace(split[i]))
	}
	return
}

func markPrimary(columns *[]c.Column, keys []string) {
	for _, key := range keys {
		for i := range *columns {
			if key == (*columns)[i].ColumnName.RawName {
				(*columns)[i].PrimaryKey = true
			}
		}
	}
}

func theRestParse(col *c.Column, theRest string) {
	// parsing for primary, auto_increment, etc
	if strings.Contains(theRest, "auto_increment") || strings.Contains(theRest, "autoincrement") {
		col.DBType = "autoincrement"
	}
	if strings.Contains(theRest, "primary") {
		col.PrimaryKey = true
	}
	if !(strings.Contains(theRest, "not null") && strings.Contains(theRest, "null")) {
		col.Null = true
	}
	defaultIdx := strings.Index(theRest, "default")
	if defaultIdx != -1 {
		defaultValue := theRest[defaultIdx+8:]
		spaceIdx := strings.Index(defaultValue, " ")
		if spaceIdx == -1 {
			spaceIdx = len(defaultValue)
		}
		defaultCheck := defaultValue[:spaceIdx]
		if !(defaultCheck == "not" || defaultCheck == "null") {
			col.DefaultValue = defaultCheck
		}
	}
}

func setGoType(col *c.Column) {
	varcharMatch := regexp.MustCompile(`^varchar[\(]\d+[\)]`)
	varyingMatch := regexp.MustCompile(`^varying[\(]\d+[\)]`)
	charMatch := regexp.MustCompile(`^char[\(]\d+[\)]`)
	binaryMatch := regexp.MustCompile(`^binary[\(]\d+[\)]`)
	varbinaryMatch := regexp.MustCompile(`^binary[\(]\d+[\)]`)
	switch col.DBType {
	case "tinyint":
		col.GoType = "null.Bool"
		col.GoTypeNonSql = "bool"
	case "int":
		col.GoType = "null.Int"
		col.GoTypeNonSql = "int"
	case "numeric":
		col.GoType = "null.Float"
		col.GoTypeNonSql = "float64"
	case "decimal", "dec":
		col.GoType = "null.Float"
		col.GoTypeNonSql = "float64"
	case "float":
		col.GoType = "null.Float"
		col.GoTypeNonSql = "float64"
	case "double":
		col.GoType = "null.Float"
		col.GoTypeNonSql = "float64"
	case "real":
		col.GoType = "null.Float"
		col.GoTypeNonSql = "float64"
	case "money":
		col.GoType = "null.Float"
		col.GoTypeNonSql = "float64"
	case "text", "tinytext", "mediumtext", "longtext":
		col.GoType = "null.String"
		col.GoTypeNonSql = "string"
	case "json":
		col.GoType = "*json.RawMessage"
		col.GoTypeNonSql = "[]byte"
	case "tinyblob", "blob", "mediumblob", "longblob":
		col.GoType = "null.Byte"
		col.GoTypeNonSql = "[]byte"
	case "time", "date", "datetime", "timestamp":
		col.GoType = "null.Time"
		col.GoTypeNonSql = "time.Time"
	case "uuid":
		col.GoType = "string"
		col.GoTypeNonSql = "string"
	case "autoincrement":
		col.GoType = "int"
		col.GoTypeNonSql = "int"
		col.Null = false
	case "serial":
		col.DBType = "autoincrement"
		col.GoType = "int"
		col.GoTypeNonSql = "int"
		col.Null = false
	case "boolean", "bool":
		col.DBType = "boolean"
		col.GoType = "null.Bool"
		col.GoTypeNonSql = "bool"
	}
	switch {
	case varcharMatch.MatchString(col.DBType):
		length, err := splitChar(col.DBType)
		if err != nil {
			fmt.Printf("Syntax error in getting length from varchar field: %s\n", col.DBType)
		}
		col.Length = length
		col.DBType = "varchar"
		col.GoType = "null.String"
		col.GoTypeNonSql = "string"
	case varyingMatch.MatchString(col.DBType):
		length, err := splitChar(col.DBType)
		if err != nil {
			fmt.Printf("Syntax error in getting length from varying field: %s\n", col.DBType)
		}
		col.Length = length
		col.DBType = "varying"
		col.GoType = "null.String"
		col.GoTypeNonSql = "string"
	case charMatch.MatchString(col.DBType):
		length, err := splitChar(col.DBType)
		if err != nil {
			fmt.Printf("Syntax error in getting length from char field: %s\n", col.DBType)
		}
		col.Length = length
		col.DBType = "char"
		col.GoType = "null.String"
		col.GoTypeNonSql = "string"
	case binaryMatch.MatchString(col.DBType):
		length, err := splitChar(col.DBType)
		if err != nil {
			fmt.Printf("Syntax error in getting length from binary field: %s\n", col.DBType)
		}
		col.Length = length
		col.DBType = "binary"
		col.GoType = "null.Byte"
		col.GoTypeNonSql = "byte"
	case varbinaryMatch.MatchString(col.DBType):
		length, err := splitChar(col.DBType)
		if err != nil {
			fmt.Printf("Syntax error in getting length from varbinary field: %s\n", col.DBType)
		}
		col.Length = length
		col.DBType = "varbinary"
		col.GoType = "null.Byte"
		col.GoTypeNonSql = "[]byte"
	}
}

func splitChar(strChar string) (length int64, err error) {
	paranOpenIdx := strings.Index(strChar, "(")
	paranCloseIdx := strings.Index(strChar, ")")
	if paranOpenIdx == -1 || paranCloseIdx == -1 {
		err = fmt.Errorf("Parse error for char column")
		return
	}
	lengthStr := strChar[paranOpenIdx+1 : paranCloseIdx]
	length, err = strconv.ParseInt(lengthStr, 10, 64)
	return
}

func determineColExists(columns []c.Column) (colExist c.ColumnExistence) {
	for _, col := range columns {
		if col.GoType == "null.Time" {
			colExist.TimeColumn = true
		}
		if col.Null {
			colExist.HaveNullColumns = true
		}
	}
	return
}

func PrintSqlColumns(cols []c.Column) {
	fmt.Println("--- Saved Columns ---")
	tab := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	fmt.Fprintln(tab, "|Name\t|Type\t|Default Value\t|Null\t|Primary Key")
	fmt.Fprintln(tab, "+----\t+----\t+-------------\t+----\t+-----------")
	for _, col := range cols {
		fmt.Fprintf(tab, " %s\t %s\t %s\t %t\t %t\n", col.ColumnName.Camel, col.DBType, col.DefaultValue, col.Null, col.PrimaryKey)
	}
	tab.Flush()
	fmt.Println("")
}
