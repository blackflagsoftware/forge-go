package entity

import (
	"fmt"
	"strings"

	con "github.com/blackflagsoftware/forge-go/internal/constant"
	"github.com/blackflagsoftware/forge-go/internal/util"
)

func (ep *Entity) BuildManagerTemplate() {
	getDeleteRow := ""
	postRow := "\t"
	putRow := ""
	patchRow := ""
	patchSearch := ""
	patchInit := []string{}
	patchTestInit := []string{}
	setArgs := ""
	foundOneKey := false
	getDeleteKeyTestSuccessful := []string{}
	getDeleteKeyTestFailure := []string{}
	for _, c := range ep.Columns {
		if c.PrimaryKey {
			if c.GoType == "string" {
				if foundOneKey {
					setArgs += ", "
					patchSearch += "\n"
					getDeleteRow += "\n"
				}
				getDeleteRow += fmt.Sprintf(con.MANAGER_GET_STRING, ep.Abbr, c.ColumnName.Camel, c.ColumnName.Camel)
				patchSearch += fmt.Sprintf(con.MANAGER_PATCH_SEARCH_STRING, c.ColumnName.Lower, c.ColumnName.Camel, c.ColumnName.Camel, c.ColumnName.Camel, c.ColumnName.Camel)
				getDeleteKeyTestSuccessful = append(getDeleteKeyTestSuccessful, fmt.Sprintf(`%s: "test id"`, c.ColumnName.Camel))
				getDeleteKeyTestFailure = append(getDeleteKeyTestFailure, fmt.Sprintf(`%s: ""`, c.ColumnName.Camel))
				foundOneKey = true
				patchTest := "test value"
				if c.DBType == "uuid" {
					patchTest = "76A21E7C-A155-4472-AEC5-14C84AC82B9A"
				}
				patchTestInit = append(patchTestInit, fmt.Sprintf("%s: \"%s\"", c.ColumnName.Camel, patchTest))
			}
			if c.GoType == "int" {
				if foundOneKey {
					setArgs += ", "
					patchSearch += "\n"
				}
				getDeleteRow += fmt.Sprintf(con.MANAGER_GET_INT, ep.Abbr, c.ColumnName.Camel, c.ColumnName.Camel)
				patchSearch += fmt.Sprintf(con.MANAGER_PATCH_SEARCH_INT, c.ColumnName.Lower, c.ColumnName.Camel, c.ColumnName.Camel, c.ColumnName.Camel, c.ColumnName.Camel, c.ColumnName.Lower, c.ColumnName.Lower)
				getDeleteKeyTestSuccessful = append(getDeleteKeyTestSuccessful, fmt.Sprintf(`%s: 1`, c.ColumnName.Camel))
				getDeleteKeyTestFailure = append(getDeleteKeyTestFailure, fmt.Sprintf(`%s: 0`, c.ColumnName.Camel))
				foundOneKey = true
				patchTestInit = append(patchTestInit, fmt.Sprintf("%s: 1", c.ColumnName.Camel))
			}
			patchInit = append(patchInit, fmt.Sprintf("%s: %sIn.%s", c.ColumnName.Camel, ep.Abbr, c.ColumnName.Camel))
			setArgs += fmt.Sprintf("%s: %s", c.ColumnName.Camel, c.ColumnName.Lower)
		}
	}
	// patchRow = patchSearch + fmt.Sprintf(MANAGER_PATCH_STRUCT_STMT, ep.Abbr, ep.Camel, setArgs) + fmt.Sprintf(MANAGER_PATCH_GET_STMT, ep.Abbr)
	ep.ManagerPatchInitArgs = strings.Join(patchInit, ", ")
	ep.ManagerPatchTestInit = strings.Join(patchTestInit, ", ")
	putTests := []PostPutTest{{Name: "successful", Failure: false}}
	postTests := []PostPutTest{{Name: "successful", Failure: false}}
	InitializeColumnTests()
	sortColumns := []string{}
	addedTime := false
	for _, c := range ep.Columns {
		columnTestStrAdded := false
		// put rows
		if c.PrimaryKey {
			putRow += fmt.Sprintf(con.MANAGER_GET_INT, ep.Abbr, c.ColumnName.Camel, c.ColumnName.Camel)
			AppendColumnTest(c.ColumnName.Camel, c.GoType, true)
			putTests = append(putTests, PostPutTest{Name: fmt.Sprintf("invalid %s", c.ColumnName.LowerCamel), Failure: true, ForColumn: c.ColumnName.Camel})
			columnTestStrAdded = true
		}
		// post rows
		if c.GoType == "string" || c.GoType == "null.String" {
			if c.DBType == "uuid" {
				postRow += fmt.Sprintf(con.MANAGER_POST_UUID, ep.Abbr, c.ColumnName.Camel, c.ColumnName.Camel)
				putRow += fmt.Sprintf(con.MANAGER_POST_UUID, ep.Abbr, c.ColumnName.Camel, c.ColumnName.Camel)
				AppendColumnTest(c.ColumnName.Camel, c.GoType, false)
				postTests = append(postTests, PostPutTest{Name: "invalid UUID", Failure: true, ForColumn: c.ColumnName.Camel})
				putTests = append(putTests, PostPutTest{Name: "invalid UUID", Failure: true, ForColumn: c.ColumnName.Camel})
				columnTestStrAdded = true
			} else {
				if !c.Null {
					postRow += fmt.Sprintf(con.MANAGER_POST_NULL, ep.Abbr, c.ColumnName.Camel, c.ColumnName.Camel)
					putRow += fmt.Sprintf(con.MANAGER_POST_NULL, ep.Abbr, c.ColumnName.Camel, c.ColumnName.Camel)
					AppendColumnTest(c.ColumnName.Camel, c.GoType, false)
					postTests = append(postTests, PostPutTest{Name: fmt.Sprintf("invalid %s", c.ColumnName.LowerCamel), Failure: true, ForColumn: c.ColumnName.Camel})
					putTests = append(putTests, PostPutTest{Name: fmt.Sprintf("invalid %s", c.ColumnName.LowerCamel), Failure: true, ForColumn: c.ColumnName.Camel})
					columnTestStrAdded = true
				}
				if c.Length > 0 {
					postRow += fmt.Sprintf(con.MANAGER_POST_VARCHAR_LEN, ep.Abbr, c.ColumnName.Camel, ep.Abbr, c.ColumnName.Camel, c.Length, c.ColumnName.Camel, c.Length)
					putRow += fmt.Sprintf(con.MANAGER_POST_VARCHAR_LEN, ep.Abbr, c.ColumnName.Camel, ep.Abbr, c.ColumnName.Camel, c.Length, c.ColumnName.Camel, c.Length)
					AppendColumnTest(c.ColumnName.Camel, c.GoType, false)
					postTests = append(postTests, PostPutTest{Name: fmt.Sprintf("length %s", c.ColumnName.LowerCamel), Failure: true, ForColumn: c.ColumnName.Camel, ColumnLength: int(c.Length)})
					putTests = append(putTests, PostPutTest{Name: fmt.Sprintf("length %s", c.ColumnName.LowerCamel), Failure: true, ForColumn: c.ColumnName.Camel, ColumnLength: int(c.Length)})
					columnTestStrAdded = true
				}
			}
		} else {
			if !c.Null && !c.PrimaryKey {
				postRow += fmt.Sprintf(con.MANAGER_POST_NULL, ep.Abbr, c.ColumnName.Camel, c.ColumnName.Camel)
				putRow += fmt.Sprintf(con.MANAGER_POST_NULL, ep.Abbr, c.ColumnName.Camel, c.ColumnName.Camel)
				AppendColumnTest(c.ColumnName.Camel, c.GoType, false)
				postTests = append(postTests, PostPutTest{Name: fmt.Sprintf("invalid %s", c.ColumnName.LowerCamel), Failure: true, ForColumn: c.ColumnName.Camel})
				putTests = append(putTests, PostPutTest{Name: fmt.Sprintf("invalid %s", c.ColumnName.LowerCamel), Failure: true, ForColumn: c.ColumnName.Camel})
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
			AppendColumnTest(c.ColumnName.Camel, c.GoType, false)
		}
		// patch rows
		if !c.PrimaryKey {
			switch c.GoType {
			case "null.String":
				if !c.PrimaryKey {
					patchLenCheck := ""
					if c.Length > 0 {
						patchLenCheck = fmt.Sprintf(con.MANAGER_PATCH_VARCHAR_LEN, ep.Abbr, c.ColumnName.Camel, ep.Abbr, c.ColumnName.Camel, c.Length, c.ColumnName.Camel, c.Length)
					}
					patchRow += fmt.Sprintf(con.MANAGER_PATCH_DEFAULT_ASSIGN, c.ColumnName.Camel, ep.Abbr, c.ColumnName.Camel, patchLenCheck, ep.Abbr, c.ColumnName.Camel, ep.Abbr, c.ColumnName.Camel)
				}
			case "null.Time":
				if !addedTime {
					ep.ManagerImportTest += "\n\t\"time\"\n"
					addedTime = true
				}
				// ColCamel, Abbr, ColCamel, Abbr, ColCamel, ColCamel, Abbr, ColCamel, Abbr, ColCamel
				patchRow += fmt.Sprintf(con.MANAGER_PATCH_TIME_NULL_ASSIGN, c.ColumnName.Camel, ep.Abbr, c.ColumnName.Camel, ep.Abbr, c.ColumnName.Camel, c.ColumnName.Camel, ep.Abbr, c.ColumnName.Camel, ep.Abbr, c.ColumnName.Camel)
				ep.ManagerTime = "\n\t\"time\""
			case "*json.RawMessage":
				patchRow += fmt.Sprintf(con.MANAGER_PATCH_JSON_NULL_ASSIGN, c.ColumnName.Camel, ep.Abbr, c.ColumnName.Camel, ep.Abbr, c.ColumnName.Camel, c.ColumnName.LowerCamel, ep.Abbr, c.ColumnName.Camel, ep.Abbr, c.ColumnName.Camel)
				ep.ManagerImportTest += "\n\t\"encoding/json\"\n"
			default:
				patchRow += fmt.Sprintf(con.MANAGER_PATCH_DEFAULT_ASSIGN, c.ColumnName.Camel, ep.Abbr, c.ColumnName.Camel, "", ep.Abbr, c.ColumnName.Camel, ep.Abbr, c.ColumnName.Camel)
			}
		}
		ep.SortColumns = strings.Join(sortColumns, ", ")
	}
	uuidColumn := ""
	for _, c := range ep.Columns {
		if c.DBType == "UUID" {
			uuidColumn = c.ColumnName.Camel
		}
	}
	if uuidColumn != "" {
		postRow += fmt.Sprintf(`%s.%s = util.GenerateUUID()`, ep.Abbr, uuidColumn)
		ep.ManagerUtilPath = fmt.Sprintf(`"%s/util"`, ep.ProjectPath)
	}
	ep.ManagerGetRow = strings.TrimRight(getDeleteRow, "\n")
	ep.ManagerPostRows = strings.TrimRight(postRow, "\n")
	ep.ManagerPutRows = strings.TrimRight(putRow, "\n")
	ep.ManagerPatchRows = strings.TrimRight(patchRow, "\n")
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
	managerPutTestRow := []string{}

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
			columnStr = append(columnStr, TranslateType(name, column.GoType, columnLength, columnValid))
		}
		managerPostTestRow = append(managerPostTestRow, fmt.Sprintf("{\n\t\t\t\"%s\",\n\t\t\t&%s{%s},\n\t\t\t%t,\n\t\t\t[]*gomock.Call{%s},\n\t\t},", postTest.Name, ep.Camel, strings.Join(columnStr, ", "), postTest.Failure, call))
	}
	for _, putTest := range putTests {
		call := ""
		if !putTest.Failure {
			call = fmt.Sprintf("mockData%s.EXPECT().Update(gomock.Any()).Return(nil).AnyTimes(),\n\t\t\t", ep.Camel)
		}
		columnStr := []string{}
		for name, column := range PutTests {
			columnValid := true
			columnLength := 0
			if putTest.ForColumn == name {
				columnValid = false
				if putTest.ColumnLength > 0 {
					columnLength = putTest.ColumnLength
				}
			}
			columnStr = append(columnStr, TranslateType(name, column.GoType, columnLength, columnValid))
		}
		managerPutTestRow = append(managerPutTestRow, fmt.Sprintf("{\n\t\t\t\"%s\",\n\t\t\t&%s{%s},\n\t\t\t%t,\n\t\t\t[]*gomock.Call{%s},\n\t\t},", putTest.Name, ep.Camel, strings.Join(columnStr, ", "), putTest.Failure, call))
	}
	ep.ManagerPostTestRow = strings.Join(managerPostTestRow, "\n\t\t")
	ep.ManagerPutTestRow = strings.Join(managerPutTestRow, "\n\t\t")
}

func InitializeColumnTests() {
	PostTests = make(map[string]ColumnTest)
	PutTests = make(map[string]ColumnTest)
}

func AppendColumnTest(name, goType string, justPut bool) {
	if _, ok := PutTests[name]; !ok {
		PutTests[name] = ColumnTest{GoType: goType}
	}
	if !justPut {
		if _, ok := PostTests[name]; !ok {
			PostTests[name] = ColumnTest{GoType: goType}
		}
	}
}

func TranslateType(columnName, columnType string, length int, valid bool) string {
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
	case "null.Bool":
		return fmt.Sprintf("%s: null.NewBool(true, %t)", columnName, valid)
	case "*json.RawMessage":
		return fmt.Sprintf("%s: &json.RawMessage{}", columnName)
	default:
		fmt.Println("Missing type in TranslateType:", columnType)
	}
	return ""
}
