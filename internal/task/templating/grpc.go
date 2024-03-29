package task

import (
	"fmt"
	"os"
	"strings"

	m "github.com/blackflagsoftware/forge-go/internal/model"
)

func buildGrpc(p *m.Project) {
	if p.CurrentEntity.MultipleKeys {
		p.GrpcImport = "\n\t\"gopkg.in/guregu/null.v3\""
	}
	if p.CurrentEntity.ModuleName != "" {
		// let the module fill in the grpc
		return
	}
	protoFile := fmt.Sprintf("%s/pkg/proto/%s.proto", p.ProjectFile.FullPath, p.AppName)
	if _, err := os.Stat(protoFile); os.IsNotExist(err) {
		fmt.Printf("%s: file not found\n", protoFile)
		return
	}
	lines := []string{}
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("message %s {", p.CurrentEntity.Camel))
	translateInLines := []string{}
	translateOutLines := []string{}
	argsInit := []string{}
	keys := []string{}
	for i, column := range p.CurrentEntity.Columns {
		if column.PrimaryKey {
			switch column.DBType {
			case "int", "autoincrement", "integer":
				columnName := fmt.Sprintf("int(in.%s)", column.ColumnName.Camel)
				if p.CurrentEntity.MultipleKeys {
					columnName = fmt.Sprintf("null.IntFrom(in.%s)", column.ColumnName.Camel)
				}
				argsInit = append(argsInit, fmt.Sprintf("%s: %s", column.ColumnName.Camel, columnName))
			default:
				columnName := fmt.Sprintf("in.%s", column.ColumnName.Camel)
				if p.CurrentEntity.MultipleKeys {
					columnName = fmt.Sprintf("null.StringFrom(in.%s)", column.ColumnName.Camel)
				}
				argsInit = append(argsInit, fmt.Sprintf("%s: %s", column.ColumnName.Camel, columnName))
			}
			idx := len(keys) + 1
			keys = append(keys, fmt.Sprintf("\t%s %s = %d;", translateGoToProtoType(column.GoTypeNonSql), column.ColumnName.Camel, idx))
		}
		idx := i + 1 // start the count at 1
		typeValue := "string"
		var inLine, outLine string
		columnType := column.GoType
		if column.PrimaryKey {
			columnType = column.GoTypeNonSql
		}
		switch columnType {
		case "float64", "null.Float":
			typeValue = "double"
			outLine = fmt.Sprintf("\tproto%s.%s = float64(%s.%s)", p.CurrentEntity.Camel, column.ColumnName.Camel, p.CurrentEntity.Abbr, column.ColumnName.Camel)
			inLine = fmt.Sprintf("\t%s.%s = in.%s", p.CurrentEntity.Abbr, column.ColumnName.Camel, column.ColumnName.Camel)
			if columnType == "null.Float" {
				outLine = fmt.Sprintf("\tproto%s.%s = %s.%s.Float64", p.CurrentEntity.Camel, column.ColumnName.Camel, p.CurrentEntity.Abbr, column.ColumnName.Camel)
				inLine = fmt.Sprintf("\t%s.%s.Scan(in.%s)", p.CurrentEntity.Abbr, column.ColumnName.Camel, column.ColumnName.Camel)
			}
		case "float32":
			typeValue = "float"
			outLine = fmt.Sprintf("\tproto%s.%s = float32(%s.%s)", p.CurrentEntity.Camel, column.ColumnName.Camel, p.CurrentEntity.Abbr, column.ColumnName.Camel)
			inLine = fmt.Sprintf("\t%s.%s = in.%s", p.CurrentEntity.Abbr, column.ColumnName.Camel, column.ColumnName.Camel)
		case "int32":
			typeValue = "int32"
			outLine = fmt.Sprintf("\tproto%s.%s = int32(%s.%s)", p.CurrentEntity.Camel, column.ColumnName.Camel, p.CurrentEntity.Abbr, column.ColumnName.Camel)
			inLine = fmt.Sprintf("\t%s.%s = in.%s", p.CurrentEntity.Abbr, column.ColumnName.Camel, column.ColumnName.Camel)
			if p.CurrentEntity.MultipleKeys {
				outLine = fmt.Sprintf("\tproto%s.%s = int32(%s.%s.Int64)", p.CurrentEntity.Camel, column.ColumnName.Camel, p.CurrentEntity.Abbr, column.ColumnName.Camel)
				inLine = fmt.Sprintf("\t%s.%s.Scan(in.%s)", p.CurrentEntity.Abbr, column.ColumnName.Camel, column.ColumnName.Camel)
			}
		case "int64", "int", "null.Int":
			typeValue = "int64"
			outLine = fmt.Sprintf("\tproto%s.%s = int64(%s.%s)", p.CurrentEntity.Camel, column.ColumnName.Camel, p.CurrentEntity.Abbr, column.ColumnName.Camel)
			inLine = fmt.Sprintf("\t%s.%s = int(in.%s)", p.CurrentEntity.Abbr, column.ColumnName.Camel, column.ColumnName.Camel)
			if columnType == "null.Int" {
				outLine = fmt.Sprintf("\tproto%s.%s = %s.%s.Int64", p.CurrentEntity.Camel, column.ColumnName.Camel, p.CurrentEntity.Abbr, column.ColumnName.Camel)
				inLine = fmt.Sprintf("\t%s.%s.Scan(in.%s)", p.CurrentEntity.Abbr, column.ColumnName.Camel, column.ColumnName.Camel)
			}
			if p.CurrentEntity.MultipleKeys {
				outLine = fmt.Sprintf("\tproto%s.%s = %s.%s.Int64", p.CurrentEntity.Camel, column.ColumnName.Camel, p.CurrentEntity.Abbr, column.ColumnName.Camel)
				inLine = fmt.Sprintf("\t%s.%s.Scan(in.%s)", p.CurrentEntity.Abbr, column.ColumnName.Camel, column.ColumnName.Camel)
			}
		case "uint32":
			typeValue = "uint32"
			outLine = fmt.Sprintf("\tproto%s.%s = uint32(%s.%s)", p.CurrentEntity.Camel, column.ColumnName.Camel, p.CurrentEntity.Abbr, column.ColumnName.Camel)
			inLine = fmt.Sprintf("\t%s.%s = in.%s", p.CurrentEntity.Abbr, column.ColumnName.Camel, column.ColumnName.Camel)
		case "uint64":
			typeValue = "uint64"
			outLine = fmt.Sprintf("\tproto%s.%s = uint64(%s.%s)", p.CurrentEntity.Camel, column.ColumnName.Camel, p.CurrentEntity.Abbr, column.ColumnName.Camel)
			inLine = fmt.Sprintf("\t%s.%s = in.%s", p.CurrentEntity.Abbr, column.ColumnName.Camel, column.ColumnName.Camel)
		case "bool", "null.Bool":
			typeValue = "bool"
			outLine = fmt.Sprintf("\tproto%s.%s = %s.%s", p.CurrentEntity.Camel, column.ColumnName.Camel, p.CurrentEntity.Abbr, column.ColumnName.Camel)
			inLine = fmt.Sprintf("\t%s.%s = in.%s", p.CurrentEntity.Abbr, column.ColumnName.Camel, column.ColumnName.Camel)
			if columnType == "null.Bool" {
				outLine = fmt.Sprintf("\tproto%s.%s = %s.%s.Bool", p.CurrentEntity.Camel, column.ColumnName.Camel, p.CurrentEntity.Abbr, column.ColumnName.Camel)
				inLine = fmt.Sprintf("\t%s.%s.Scan(in.%s)", p.CurrentEntity.Abbr, column.ColumnName.Camel, column.ColumnName.Camel)
			}
		case "[]byte":
			typeValue = "bytes"
			outLine = fmt.Sprintf("\tproto%s.%s = %s.%s", p.CurrentEntity.Camel, column.ColumnName.Camel, p.CurrentEntity.Abbr, column.ColumnName.Camel)
			inLine = fmt.Sprintf("\t%s.%s = in.%s", p.CurrentEntity.Abbr, column.ColumnName.Camel, column.ColumnName.Camel)
		case "*json.RawMessage":
			typeValue = "bytes"
			outLine = fmt.Sprintf("\tproto%s.%s, _ = %s.%s.MarshalJSON()", p.CurrentEntity.Camel, column.ColumnName.Camel, p.CurrentEntity.Abbr, column.ColumnName.Camel)
			inLine = fmt.Sprintf("\t%s.%s.UnmarshalJSON(in.%s)", p.CurrentEntity.Abbr, column.ColumnName.Camel, column.ColumnName.Camel)
		default:
			typeValue = "string"
			outLine = fmt.Sprintf("\tproto%s.%s = %s.%s", p.CurrentEntity.Camel, column.ColumnName.Camel, p.CurrentEntity.Abbr, column.ColumnName.Camel)
			inLine = fmt.Sprintf("\t%s.%s = in.%s", p.CurrentEntity.Abbr, column.ColumnName.Camel, column.ColumnName.Camel)
			if columnType == "null.String" {
				outLine = fmt.Sprintf("\tproto%s.%s = %s.%s.String", p.CurrentEntity.Camel, column.ColumnName.Camel, p.CurrentEntity.Abbr, column.ColumnName.Camel)
				inLine = fmt.Sprintf("\t%s.%s.Scan(in.%s)", p.CurrentEntity.Abbr, column.ColumnName.Camel, column.ColumnName.Camel)
			}
			if columnType == "null.Time" {
				outLine = fmt.Sprintf("\tproto%s.%s = %s.%s.Time.Format(time.RFC3339)", p.CurrentEntity.Camel, column.ColumnName.Camel, p.CurrentEntity.Abbr, column.ColumnName.Camel)
			}
			if columnType == "null.Time" {
				inLine = fmt.Sprintf("\t%s.%s.Scan(in.%s)", p.CurrentEntity.Abbr, column.ColumnName.Camel, column.ColumnName.Camel)
			}
			if p.CurrentEntity.MultipleKeys {
				outLine = fmt.Sprintf("\tproto%s.%s = %s.%s.String", p.CurrentEntity.Camel, column.ColumnName.Camel, p.CurrentEntity.Abbr, column.ColumnName.Camel)
				inLine = fmt.Sprintf("\t%s.%s.Scan(in.%s)", p.CurrentEntity.Abbr, column.ColumnName.Camel, column.ColumnName.Camel)
			}
		}
		lines = append(lines, fmt.Sprintf("\t%s %s = %d;", typeValue, column.ColumnName.Camel, idx))
		translateInLines = append(translateInLines, inLine)
		translateOutLines = append(translateOutLines, outLine)
	}
	lines = append(lines, "}")
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("message %sResponse {", p.CurrentEntity.Camel))
	lines = append(lines, fmt.Sprintf("\t%s %s = 1;", p.CurrentEntity.Camel, p.CurrentEntity.Camel))
	lines = append(lines, fmt.Sprintf("\tResult result = 2;"))
	lines = append(lines, "}")
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("message %sRepeatResponse {", p.CurrentEntity.Camel))
	lines = append(lines, fmt.Sprintf("\trepeated %s %s = 1;", p.CurrentEntity.Camel, p.CurrentEntity.Camel))
	lines = append(lines, fmt.Sprintf("\tResult result = 2;"))
	lines = append(lines, "}")
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("service %sService {", p.CurrentEntity.Camel))
	lines = append(lines, fmt.Sprintf("\trpc Get%s(%sIDIn) returns (%sResponse);", p.CurrentEntity.Camel, p.CurrentEntity.Camel, p.CurrentEntity.Camel))
	lines = append(lines, fmt.Sprintf("\trpc Search%s(%s) returns (%sRepeatResponse);", p.CurrentEntity.Camel, p.CurrentEntity.Camel, p.CurrentEntity.Camel))
	lines = append(lines, fmt.Sprintf("\trpc Create%s(%s) returns (%sResponse);", p.CurrentEntity.Camel, p.CurrentEntity.Camel, p.CurrentEntity.Camel))
	lines = append(lines, fmt.Sprintf("\trpc Update%s(%s) returns (Result);", p.CurrentEntity.Camel, p.CurrentEntity.Camel))
	lines = append(lines, fmt.Sprintf("\trpc Delete%s(%sIDIn) returns (Result);", p.CurrentEntity.Camel, p.CurrentEntity.Camel))
	lines = append(lines, "}")
	lines = append(lines, "")
	// add the IDIn lines
	lines = append(lines, fmt.Sprintf("message %sIDIn {", p.CurrentEntity.Camel))
	lines = append(lines, keys...)
	lines = append(lines, "}")
	lines = append(lines, "")

	file, err := os.OpenFile(protoFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("%s: unable to open file with error: %s\n", protoFile, err)
		return
	}
	defer file.Close()
	_, err = file.WriteString(strings.Join(lines, "\n"))
	if err != nil {
		fmt.Printf("%s: unable to write to file with error: %s\n", protoFile, err)
	}
	p.GrpcArgsInit = strings.Join(argsInit, ", ")
	// save off translateIn/Out
	p.GrpcTranslateIn = strings.Join(translateInLines, "\n\t")
	p.GrpcTranslateOut = strings.Join(translateOutLines, "\n\t")
}

func translateGoToProtoType(goType string) string {
	switch goType {
	case "float32", "float64":
		return "double"
	case "int":
		return "int64"
	case "uint32", "uint64":
		return "fixed"
	case "[]byte":
		return "byte"
	default:
		return goType
	}
}
