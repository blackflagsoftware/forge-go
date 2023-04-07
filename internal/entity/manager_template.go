package entity

import (
	"fmt"
	"strings"

	"github.com/blackflagsoftware/forge-go/internal/util"
)

// const down at the bottom

func (ep *Entity) BuildManagerTemplate() {
	ep.buildGetDelete()
	ep.buildPost()
	ep.buildPatch()
	ep.buildTest()
	// build import
	importLines := []string{}
	if ep.HasJsonColumn() {
		importLines = append(importLines, "\"encoding/json\"\n")
	}
	if ep.HasTimeColumn() {
		importLines = append(importLines, "\"time\"\n")
	}
	importLines = append(importLines, fmt.Sprintf("ae \"%s/internal/api_error\"", ep.ProjectPath))
	if ep.HasPrimaryUUIDColumn() {
		importLines = append(importLines, fmt.Sprintf("\"%s/internal/util\"", ep.ProjectPath))
	}
	// if ep.HasNullColumn() {
	// 	importLines = append(importLines, "\"gopkg.in/guregu/null.v3\"")
	// }
	ep.ManagerImport = strings.Join(importLines, "\n\t")
	// test import
	importLines = []string{}
	importLines = append(importLines, "\"fmt\"")
	importLines = append(importLines, "\"testing\"")
	if ep.HasTimeColumn() {
		importLines = append(importLines, "\"time\"")
	}
	importLines = append(importLines, "\n\"github.com/golang/mock/gomock\"")
	importLines = append(importLines, "\"github.com/stretchr/testify/assert\"")
	if ep.HasNullColumn() {
		importLines = append(importLines, "\"gopkg.in/guregu/null.v3\"")
	}
	ep.ManagerTestImport = strings.Join(importLines, "\n\t")
}

func (ep *Entity) buildGetDelete() {
	rows := []string{}
	for _, c := range ep.Columns {
		if c.PrimaryKey {
			switch c.GoType {
			case "string":
				rows = append(rows, fmt.Sprintf(GET_DELETE_STRING, ep.Abbr, c.ColumnName.Camel, c.ColumnName.Camel))
			case "int":
				rows = append(rows, fmt.Sprintf(GET_DELETE_INT, ep.Abbr, c.ColumnName.Camel, c.ColumnName.Camel))
			}
		}
	}
	ep.ManagerGetRows = strings.Join(rows, "\n\t")
}

func (ep *Entity) buildPost() {
	rows := []string{}
	uuidColumn := ""
	setCreatedAt := false
	for _, c := range ep.Columns {
		if c.DBType == "uuid" && c.PrimaryKey {
			uuidColumn = c.ColumnName.Camel
		}
		if c.PrimaryKey {
			continue
		}
		if c.GoType == "string" {
			rows = append(rows, fmt.Sprintf(POST_STRING, ep.Abbr, c.ColumnName.Camel, c.ColumnName.Camel))
			if c.Length > 0 {
				rows = append(rows, fmt.Sprintf(POST_STRING_LEN, ep.Abbr, c.ColumnName.Camel, c.Length, c.ColumnName.Camel, c.Length))
			}
		} else if c.GoType == "null.String" {
			if !c.Null {
				rows = append(rows, fmt.Sprintf(POST_NULL, ep.Abbr, c.ColumnName.Camel, c.ColumnName.Camel))
			}
			if c.Length > 0 {
				rows = append(rows, fmt.Sprintf(POST_NULL_LEN, ep.Abbr, c.ColumnName.Camel, ep.Abbr, c.ColumnName.Camel, c.Length, c.ColumnName.Camel, c.Length))
			}
		} else {
			if !c.Null {
				rows = append(rows, fmt.Sprintf(POST_NULL, ep.Abbr, c.ColumnName.Camel, c.ColumnName.Camel))
			}
		}
		if c.ColumnName.Lower == "created_at" {
			setCreatedAt = true
		}
	}
	if setCreatedAt {
		rows = append(rows, fmt.Sprintf("%s.CreatedAt.Scan(time.Now().UTC())", ep.Abbr))
	}
	if uuidColumn != "" {
		rows = append(rows, fmt.Sprintf("%s.%s = util.GenerateUUID()", ep.Abbr, uuidColumn))
	}
	ep.ManagerPostRows = strings.Join(rows, "\n\t")
}

func (ep *Entity) buildPatch() {
	rows := []string{}
	patchInit := []string{}
	setUpdatedAt := false
	for _, c := range ep.Columns {
		if c.PrimaryKey {
			patchInit = append(patchInit, fmt.Sprintf("%s: %sIn.%s", c.ColumnName.Camel, ep.Abbr, c.ColumnName.Camel))
			continue
		}
		switch c.GoType {
		case "string":
			// mostly used for uuid
			rows = append(rows, fmt.Sprintf(PATCH_STRING_ASSIGN, c.ColumnName.Camel, ep.Abbr, c.ColumnName.Camel, ep.Abbr, c.ColumnName.Camel, ep.Abbr, c.ColumnName.Camel))
		case "null.String":
			patchLenCheck := ""
			if c.Length > 0 {
				patchLenCheck = fmt.Sprintf(PATCH_VARCHAR_LEN, ep.Abbr, c.ColumnName.Camel, ep.Abbr, c.ColumnName.Camel, c.Length, c.ColumnName.Camel, c.Length)
			}
			rows = append(rows, fmt.Sprintf(PATCH_DEFAULT_ASSIGN, c.ColumnName.Camel, ep.Abbr, c.ColumnName.Camel, patchLenCheck, ep.Abbr, c.ColumnName.Camel, ep.Abbr, c.ColumnName.Camel))
		case "null.Time":
			if !(c.ColumnName.Lower == "created_at" || c.ColumnName.Lower == "updated_at") {
				rows = append(rows, fmt.Sprintf(PATCH_TIME_NULL_ASSIGN, c.ColumnName.Camel, ep.Abbr, c.ColumnName.Camel, ep.Abbr, c.ColumnName.Camel, c.ColumnName.Camel, ep.Abbr, c.ColumnName.Camel, ep.Abbr, c.ColumnName.Camel))
			}
		case "*json.RawMessage":
			rows = append(rows, fmt.Sprintf(PATCH_JSON_NULL_ASSIGN, c.ColumnName.Camel, ep.Abbr, c.ColumnName.Camel, ep.Abbr, c.ColumnName.Camel, c.ColumnName.LowerCamel, ep.Abbr, c.ColumnName.Camel, ep.Abbr, c.ColumnName.Camel))
		default:
			rows = append(rows, fmt.Sprintf(PATCH_DEFAULT_ASSIGN, c.ColumnName.Camel, ep.Abbr, c.ColumnName.Camel, "", ep.Abbr, c.ColumnName.Camel, ep.Abbr, c.ColumnName.Camel))
		}
		if c.ColumnName.Lower == "updated_at" {
			setUpdatedAt = true
		}
	}
	if setUpdatedAt {
		rows = append(rows, fmt.Sprintf("\n\t%s.UpdatedAt.Scan(time.Now().UTC())", ep.Abbr))
	}
	ep.ManagerPatchRows = strings.Join(rows, "\n\t")
	ep.ManagerPatchInitArgs = strings.Join(patchInit, ", ")
}

func (ep *Entity) buildTest() {
	patchTestInit := []string{}
	getDeleteKeyTestSuccessful := []string{}
	getDeleteKeyTestFailure := []string{}
	for _, c := range ep.Columns {
		if c.PrimaryKey {
			if c.GoType == "string" {
				getDeleteKeyTestSuccessful = append(getDeleteKeyTestSuccessful, fmt.Sprintf(`%s: "test id"`, c.ColumnName.Camel))
				getDeleteKeyTestFailure = append(getDeleteKeyTestFailure, fmt.Sprintf(`%s: ""`, c.ColumnName.Camel))
				patchTest := "test value"
				if c.DBType == "uuid" {
					patchTest = "76A21E7C-A155-4472-AEC5-14C84AC82B9A"
				}
				patchTestInit = append(patchTestInit, fmt.Sprintf("%s: \"%s\"", c.ColumnName.Camel, patchTest))
			}
			if c.GoType == "int" {
				getDeleteKeyTestSuccessful = append(getDeleteKeyTestSuccessful, fmt.Sprintf(`%s: 1`, c.ColumnName.Camel))
				getDeleteKeyTestFailure = append(getDeleteKeyTestFailure, fmt.Sprintf(`%s: 0`, c.ColumnName.Camel))
				patchTestInit = append(patchTestInit, fmt.Sprintf("%s: 1", c.ColumnName.Camel))
			}
		}
	}
	ep.ManagerPatchTestInit = strings.Join(patchTestInit, ", ")
	postTests := []PostPutTest{{Name: "successful", Failure: false}}
	InitializeColumnTests()
	sortColumns := []string{}
	for _, c := range ep.Columns {
		columnTestStrAdded := false
		if c.GoType == "string" || c.GoType == "null.String" {
			if c.DBType == "uuid" {
				if !c.PrimaryKey {
					AppendColumnTest(c.ColumnName.Camel, c.GoType, c.DBType, false)
					postTests = append(postTests, PostPutTest{Name: fmt.Sprintf("invalid %s", c.ColumnName.LowerCamel), Failure: true, ForColumn: c.ColumnName.Camel})
					columnTestStrAdded = true
				}
			} else {
				if !c.Null {
					AppendColumnTest(c.ColumnName.Camel, c.GoType, c.DBType, false)
					postTests = append(postTests, PostPutTest{Name: fmt.Sprintf("invalid %s", c.ColumnName.LowerCamel), Failure: true, ForColumn: c.ColumnName.Camel})
					columnTestStrAdded = true
				}
				if c.Length > 0 {
					AppendColumnTest(c.ColumnName.Camel, c.GoType, c.DBType, false)
					postTests = append(postTests, PostPutTest{Name: fmt.Sprintf("length %s", c.ColumnName.LowerCamel), Failure: true, ForColumn: c.ColumnName.Camel, ColumnLength: int(c.Length)})
					columnTestStrAdded = true
				}
			}
		} else {
			if !c.Null && !c.PrimaryKey {
				AppendColumnTest(c.ColumnName.Camel, c.GoType, c.DBType, false)
				postTests = append(postTests, PostPutTest{Name: fmt.Sprintf("invalid %s", c.ColumnName.LowerCamel), Failure: true, ForColumn: c.ColumnName.Camel})
				columnTestStrAdded = true
			}
		}
		if !c.PrimaryKey {
			if ep.DefaultColumn == "" {
				ep.DefaultColumn = c.ColumnName.RawName
			}
			sortColumns = append(sortColumns, fmt.Sprintf("\"%s\": \"%s\"", c.ColumnName.RawName, c.ColumnName.RawName))
		}
		if !columnTestStrAdded && c.Null && !c.PrimaryKey {
			// add column to all the other tests with good data
			AppendColumnTest(c.ColumnName.Camel, c.GoType, c.DBType, false)
		}
		// patch rows
		ep.SortColumns = strings.Join(sortColumns, ", ")
	}
	managerGetSuccessfulRow := ""
	managerGetFailureRow := ""
	managerDeleteSuccessfulRow := ""
	managerDeleteFailureRow := ""
	if len(getDeleteKeyTestSuccessful) > 0 {
		managerGetSuccessfulRow = fmt.Sprintf("{\n\t\t\t\"successful\",\n\t\t\t&%s{%s},\n\t\t\tfalse,\n\t\t\t[]*gomock.Call{\n\t\t\t\tmockData%s.EXPECT().Read(gomock.Any()).Return(nil),\n\t\t\t},\n\t\t},", ep.Camel, strings.Join(getDeleteKeyTestSuccessful, ", "), ep.Camel)
		managerDeleteSuccessfulRow = fmt.Sprintf("{\n\t\t\t\"successful\",\n\t\t\t&%s{%s},\n\t\t\tfalse,\n\t\t\t[]*gomock.Call{\n\t\t\t\tmockData%s.EXPECT().Delete(gomock.Any()).Return(nil),\n\t\t\t},\n\t\t},", ep.Camel, strings.Join(getDeleteKeyTestSuccessful, ", "), ep.Camel)
	}
	if len(getDeleteKeyTestFailure) > 0 {
		managerGetFailureRow = fmt.Sprintf("{\n\t\t\t\"invalid id\",\n\t\t\t&%s{%s},\n\t\t\ttrue,\n\t\t\t[]*gomock.Call{},\n\t\t},", ep.Camel, strings.Join(getDeleteKeyTestFailure, ", "))
		managerDeleteFailureRow = fmt.Sprintf("{\n\t\t\t\"invalid id\",\n\t\t\t&%s{%s},\n\t\t\ttrue,\n\t\t\t[]*gomock.Call{},\n\t\t},", ep.Camel, strings.Join(getDeleteKeyTestFailure, ", "))
	}
	ep.ManagerGetTestRow = fmt.Sprintf("%s\n\t\t%s", managerGetSuccessfulRow, managerGetFailureRow)
	ep.ManagerDeleteTestRow = fmt.Sprintf("%s\n\t\t%s", managerDeleteSuccessfulRow, managerDeleteFailureRow)
	managerPostTestRow := []string{}

	for _, postTest := range postTests {
		call := ""
		if !postTest.Failure {
			call = fmt.Sprintf("mockData%s.EXPECT().Create(gomock.Any()).Return(nil).AnyTimes(),\n\t\t\t", ep.Camel)
		}
		columnStr := []string{}
		for name, column := range PostTests {
			columnValid := true
			columnLength := 0
			if postTest.ForColumn == name {
				columnValid = false
				if postTest.ColumnLength > 0 {
					columnLength = postTest.ColumnLength
				}
			}
			columnStr = append(columnStr, TranslateType(name, column.GoType, column.DBType, columnLength, columnValid))
		}
		managerPostTestRow = append(managerPostTestRow, fmt.Sprintf("{\n\t\t\t\"%s\",\n\t\t\t&%s{%s},\n\t\t\t%t,\n\t\t\t[]*gomock.Call{%s},\n\t\t},", postTest.Name, ep.Camel, strings.Join(columnStr, ", "), postTest.Failure, call))
	}
	ep.ManagerPostTestRow = strings.Join(managerPostTestRow, "\n\t\t")
}

func InitializeColumnTests() {
	PostTests = make(map[string]ColumnTest)
	PutTests = make(map[string]ColumnTest)
}

func AppendColumnTest(name, goType, dbType string, justPut bool) {
	if !justPut {
		if _, ok := PostTests[name]; !ok {
			PostTests[name] = ColumnTest{GoType: goType, DBType: dbType}
		}
	}
}

func TranslateType(columnName, columnType, dbType string, length int, valid bool) string {
	switch columnType {
	case "null.String":
		if length > 0 {
			return fmt.Sprintf("%s: null.NewString(\"%s\", true)", columnName, util.BuildRandomString(length))
		}
		return fmt.Sprintf("%s: null.NewString(\"a\", %t)", columnName, valid)
	case "string":
		if length > 0 {
			return fmt.Sprintf("%s: \"%s\"", columnName, util.BuildRandomString(length))
		}
		if dbType == "uuid" && !valid {
			return fmt.Sprintf("%s: \"\"", columnName)
		}
		return fmt.Sprintf("%s: \"a\"", columnName)
	case "int":
		value := 1
		if !valid {
			value = 0
		}
		return fmt.Sprintf("%s: %d", columnName, value)
	case "null.Int":
		return fmt.Sprintf("%s: null.NewInt(1, %t)", columnName, valid)
	case "null.Float":
		return fmt.Sprintf("%s: null.NewFloat(1.0, %t)", columnName, valid)
	case "null.Time":
		return fmt.Sprintf("%s: null.NewTime(time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC), %t)", columnName, valid)
	case "time.Time":
		return fmt.Sprintf("%s: time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)", columnName)
	case "null.Bool":
		return fmt.Sprintf("%s: null.NewBool(true, %t)", columnName, valid)
	case "*json.RawMessage":
		return fmt.Sprintf("%s: &json.RawMessage{}", columnName)
	default:
		fmt.Println("Missing type in TranslateType:", columnType)
	}
	return ""
}

const (
	GET_DELETE_INT = `if %s.%s < 1 {
		return ae.MissingParamError("%s")
	}`  // Abbr, Camel, Camel

	GET_DELETE_STRING = `if %s.%s == "" {
		return ae.MissingParamError("%s")
	}`  // Abbr, Camel, Camel

	POST_STRING = `if %s.%s == "" {
		return ae.MissingParamError("%s")
	}`  // Abbr, Lower, Camel
	POST_NULL = `if !%s.%s.Valid {
		return ae.MissingParamError("%s")
	}`  // Abbr, Lower, Camel
	POST_NULL_LEN = `if %s.%s.Valid && len(%s.%s.ValueOrZero()) > %d {
		return ae.StringLengthError("%s", %d)
	}`  // Abbr, ColumnCamel, Abbr, ColumnCamel, ColumnLength, ColumnCamel, ColumnLength
	POST_STRING_LEN = `if len(%s.%s) > %d {
		return ae.StringLengthError("%s", %d)
	}`  // Abbr, ColumnCamel, ColumnLength, ColumnCamel, ColumnLength

	PATCH_STRING_ASSIGN = `// %s
	if %sIn.%s != "" {
		%s.%s = %sIn.%s
	}`  // ColCamel, Abbr, ColCamel, Abbr, ColCamel, Abbr, ColCamel
	PATCH_DEFAULT_ASSIGN = `// %s
	if %sIn.%s.Valid {%s
		%s.%s = %sIn.%s
	}`  // ColCamel, Abbr, ColCamel, StringLenCheck, Abbr, ColCamel, Abbr. ColCamel
	PATCH_JSON_NULL_ASSIGN = `// %s
	if %sIn.%s != nil {
		if !util.ValidJson(*%sIn.%s) {
			return ae.ParseError("Invalid JSON syntax for %s")
		}
		%s.%s = %sIn.%s
	}`  // ColCamel, Abbr, ColCamel, Abbr, ColCamel, ColLowerCamel, Abbr, ColCamel, Abbr, ColCamel
	PATCH_TIME_NULL_ASSIGN = `// %s
	if %sIn.%s.Valid {
		_, errParse := time.Parse(time.RFC3339, %s.%s.Time.String())
		if errParse != nil {
			return ae.ParseError("%s: unable to parse time")
		}
		%s.%s = %sIn.%s
	}`  // ColCamel, Abbr, ColCamel, Abbr, ColCamel, ColCamel, Abbr, ColCamel, Abbr, ColCamel
	PATCH_VARCHAR_LEN = `
		if %sIn.%s.Valid && len(%sIn.%s.ValueOrZero()) > %d {
			return ae.StringLengthError("%s", %d)
	}`  // Abbr, ColumnCamel, Abbr, ColumnCamel, ColumnLength, ColumnCamel, ColumnLength
)
