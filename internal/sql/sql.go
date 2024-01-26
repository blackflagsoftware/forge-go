package sql

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	m "github.com/blackflagsoftware/forge-go/internal/model"
	"github.com/blackflagsoftware/forge-go/internal/util"
)

func ParseSqlLines(lines []string) (sqlEntity SqlEntity, err error) {
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
		err = fmt.Errorf("tableParse: sql parse invalid, 'create table' not in correct format")
		return
	}
	tableNameIdx := reg.SubexpIndex("table_name")
	columnPartIdx := reg.SubexpIndex("column_part")
	tableName = matches[tableNameIdx]
	columnPart = matches[columnPartIdx]
	return
}

func columnsParse(columnPart string) (columns []m.Column, err error) {
	// parse primary key () line
	columnPart, keys := parsePrimaryKey(columnPart)
	columnPart = replaceNumericPart(columnPart)
	reg := regexp.MustCompile(`(?P<field_name>[a-zA-Z_]+) (?P<field_type>[a-zA-Z0-9\(\)_]+)(?P<the_rest>.+)?`)
	// split columns and go through each column
	lines := strings.Split(columnPart, ",")
	for i := range lines {
		line := strings.ReplaceAll(strings.TrimSpace(strings.ToLower(lines[i])), "`", "")
		if strings.Index(line, "primary") == 0 {
			continue
		}
		matches := reg.FindStringSubmatch(line)
		if len(matches) < 2 {
			err = fmt.Errorf("columnsParse: sql parse invalid, '<column name> <column type>...' not in correct format, skipping column")
			continue
		}
		fieldNameIdx := reg.SubexpIndex("field_name")
		fieldTypeIdx := reg.SubexpIndex("field_type")
		theRestIdx := reg.SubexpIndex("the_rest")
		col := m.Column{}
		rawName := matches[fieldNameIdx]
		col.ColumnName.BuildName(rawName, []string{})
		col.DBType = matches[fieldTypeIdx]
		theRest := strings.TrimSpace(matches[theRestIdx])
		theRestParse(&col, theRest) // TODO: this could have PRIMARY KEY
		setGoType(&col)
		columns = append(columns, col)
	}
	markPrimary(&columns, keys)
	return
}

func replaceNumericPart(columnPart string) string {
	// kind of a hack but need to get rid of the comma so the splitting of column lines goes smoothly
	reg := regexp.MustCompile(`\((\d+), *?(\d+)\)`)
	columnPart = reg.ReplaceAllString(columnPart, `(${1}_${2})`)
	return columnPart
}

func parsePrimaryKey(columnPart string) (newColumnPart string, keyNames []string) {
	newColumnPart = columnPart
	reg := regexp.MustCompile(`(?i)primary key ?\((?P<keys>.+)\)`)

	matches := reg.FindStringSubmatch(columnPart)
	if len(matches) < 1 {
		return
	}
	keysIdx := reg.SubexpIndex("keys")
	keys := matches[keysIdx]
	split := strings.Split(keys, ",")
	for i := range split {
		keyNames = append(keyNames, strings.TrimSpace(split[i]))
	}
	// replace the ',' in (key1, key2) if needed
	regReplace := regexp.MustCompile(`(?i)primary key ?\((.+), *?(.+)\)`)
	newColumnPart = regReplace.ReplaceAllString(columnPart, `primary key (${1}_${2})`)
	return
}

func markPrimary(columns *[]m.Column, keys []string) {
	for _, key := range keys {
		for i := range *columns {
			if key == (*columns)[i].ColumnName.RawName {
				(*columns)[i].PrimaryKey = true
			}
		}
	}
}

func theRestParse(col *m.Column, theRest string) {
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

func setGoType(col *m.Column) {
	varcharMatch := regexp.MustCompile(`^varchar[\(]\d+[\)]`)
	varyingMatch := regexp.MustCompile(`^varying[\(]\d+[\)]`)
	charMatch := regexp.MustCompile(`^char[\(]\d+[\)]`)
	binaryMatch := regexp.MustCompile(`^binary[\(]\d+[\)]`)
	varbinaryMatch := regexp.MustCompile(`^binary[\(]\d+[\)]`)
	switch {
	case col.DBType == "tinyint":
		col.GoType = "null.Bool"
		col.GoTypeNonSql = "bool"
	case strings.Contains(col.DBType, "int"):
		col.GoType = "null.Int"
		col.GoTypeNonSql = "int"
	case strings.Contains(col.DBType, "numeric"):
		col.DBType = strings.Replace(col.DBType, "_", ",", 1) // replace the comma as it should, see replaceNumericPart()
		col.GoType = "null.Float"
		col.GoTypeNonSql = "float64"
	case strings.Contains(col.DBType, "dec"):
		col.DBType = strings.Replace(col.DBType, "_", ",", 1) // replace the comma as it should, see replaceNumericPart()
		col.GoType = "null.Float"
		col.GoTypeNonSql = "float64"
	case strings.Contains(col.DBType, "float"):
		col.GoType = "null.Float"
		col.GoTypeNonSql = "float64"
	case strings.Contains(col.DBType, "double"):
		col.GoType = "null.Float"
		col.GoTypeNonSql = "float64"
	case strings.Contains(col.DBType, "real"):
		col.GoType = "null.Float"
		col.GoTypeNonSql = "float64"
	case col.DBType == "money":
		col.GoType = "null.Float"
		col.GoTypeNonSql = "float64"
	case strings.Contains(col.DBType, "text"):
		col.GoType = "null.String"
		col.GoTypeNonSql = "string"
	case col.DBType == "json":
		col.GoType = "*json.RawMessage"
		col.GoTypeNonSql = "[]byte"
	case strings.Contains(col.DBType, "blob"):
		col.GoType = "null.Byte"
		col.GoTypeNonSql = "[]byte"
	case strings.Contains(col.DBType, "time"):
		col.GoType = "null.Time"
		col.GoTypeNonSql = "time.Time"
	case strings.Contains(col.DBType, "date"):
		col.GoType = "null.Time"
		col.GoTypeNonSql = "time.Time"
	case col.DBType == "uuid":
		col.GoType = "string"
		col.GoTypeNonSql = "string"
		if !col.PrimaryKey {
			col.GoType = "null.String"
		}
	case col.DBType == "autoincrement":
		col.GoType = "int"
		col.GoTypeNonSql = "int"
		col.Null = false
	case col.DBType == "serial":
		col.DBType = "autoincrement"
		col.GoType = "int"
		col.GoTypeNonSql = "int"
		col.Null = false
	case strings.Contains(col.DBType, "bool"):
		col.DBType = "boolean"
		col.GoType = "null.Bool"
		col.GoTypeNonSql = "bool"
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
	default:
		fmt.Printf("column type: %s is invalid, setting to string with length of 100\n", col.DBType)
		col.Length = 100
		col.DBType = "varchar"
		col.GoType = "null.String"
		col.GoTypeNonSql = "string"
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

func determineColExists(columns []m.Column) (colExist m.ColumnExistence) {
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
