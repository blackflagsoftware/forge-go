package entity

import (
	"fmt"
	"os"
	"strings"
)

func (ep *Entity) BuildGrpc() {
	protoFile := fmt.Sprintf("%s/pkg/proto/%s.proto", ep.ProjectFile.FullPath, ep.AppName)
	if _, err := os.Stat(protoFile); os.IsNotExist(err) {
		fmt.Printf("%s: file not found\n", protoFile)
		return
	}
	lines := []string{}
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("message %s {", ep.Camel))
	translateInLines := []string{}
	translateOutLines := []string{}
	argsInit := []string{}
	for i, column := range ep.Columns {
		if column.PrimaryKey {
			switch column.DBType {
			case "int":
				argsInit = append(argsInit, fmt.Sprintf("%s: int(in.%s)", column.ColumnName.Camel, column.ColumnName.Camel))
			default:
				argsInit = append(argsInit, fmt.Sprintf("%s: in.%s", column.ColumnName.Camel, column.ColumnName.Camel))
			}
		}
		idx := i + 1 // start the count at 1
		typeValue := "string"
		var inLine, outLine string
		switch column.GoType {
		case "float64", "null.Float":
			typeValue = "double"
			outLine = fmt.Sprintf("\tproto%s.%s = float64(%s.%s)", ep.Camel, column.ColumnName.Camel, ep.Abbr, column.ColumnName.Camel)
			inLine = fmt.Sprintf("\t%s.%s = in.%s", ep.Abbr, column.ColumnName.Camel, column.ColumnName.Camel)
			if column.GoType == "null.Float" {
				outLine = fmt.Sprintf("\tproto%s.%s = %s.%s.Float64", ep.Camel, column.ColumnName.Camel, ep.Abbr, column.ColumnName.Camel)
				inLine = fmt.Sprintf("\t%s.%s.Scan(in.%s)", ep.Abbr, column.ColumnName.Camel, column.ColumnName.Camel)
			}
		case "float32":
			typeValue = "float"
			outLine = fmt.Sprintf("\tproto%s.%s = float32(%s.%s)", ep.Camel, column.ColumnName.Camel, ep.Abbr, column.ColumnName.Camel)
			inLine = fmt.Sprintf("\t%s.%s = in.%s", ep.Abbr, column.ColumnName.Camel, column.ColumnName.Camel)
		case "int32":
			typeValue = "int32"
			outLine = fmt.Sprintf("\tproto%s.%s = int32(%s.%s)", ep.Camel, column.ColumnName.Camel, ep.Abbr, column.ColumnName.Camel)
			inLine = fmt.Sprintf("\t%s.%s = in.%s", ep.Abbr, column.ColumnName.Camel, column.ColumnName.Camel)
		case "int64", "int", "null.Int":
			typeValue = "int64"
			outLine = fmt.Sprintf("\tproto%s.%s = int64(%s.%s)", ep.Camel, column.ColumnName.Camel, ep.Abbr, column.ColumnName.Camel)
			inLine = fmt.Sprintf("\t%s.%s = int(in.%s)", ep.Abbr, column.ColumnName.Camel, column.ColumnName.Camel)
			if column.GoType == "null.Int" {
				outLine = fmt.Sprintf("\tproto%s.%s = %s.%s.Int64", ep.Camel, column.ColumnName.Camel, ep.Abbr, column.ColumnName.Camel)
				inLine = fmt.Sprintf("\t%s.%s.Scan(in.%s)", ep.Abbr, column.ColumnName.Camel, column.ColumnName.Camel)
			}
		case "uint32":
			typeValue = "uint32"
			outLine = fmt.Sprintf("\tproto%s.%s = uint32(%s.%s)", ep.Camel, column.ColumnName.Camel, ep.Abbr, column.ColumnName.Camel)
			inLine = fmt.Sprintf("\t%s.%s = in.%s", ep.Abbr, column.ColumnName.Camel, column.ColumnName.Camel)
		case "uint64":
			typeValue = "uint64"
			outLine = fmt.Sprintf("\tproto%s.%s = uint64(%s.%s)", ep.Camel, column.ColumnName.Camel, ep.Abbr, column.ColumnName.Camel)
			inLine = fmt.Sprintf("\t%s.%s = in.%s", ep.Abbr, column.ColumnName.Camel, column.ColumnName.Camel)
		case "bool", "null.Bool":
			typeValue = "bool"
			outLine = fmt.Sprintf("\tproto%s.%s = %s.%s", ep.Camel, column.ColumnName.Camel, ep.Abbr, column.ColumnName.Camel)
			inLine = fmt.Sprintf("\t%s.%s = in.%s", ep.Abbr, column.ColumnName.Camel, column.ColumnName.Camel)
			if column.GoType == "null.Bool" {
				outLine = fmt.Sprintf("\tproto%s.%s = %s.%s.Bool", ep.Camel, column.ColumnName.Camel, ep.Abbr, column.ColumnName.Camel)
				inLine = fmt.Sprintf("\t%s.%s.Scan(in.%s)", ep.Abbr, column.ColumnName.Camel, column.ColumnName.Camel)
			}
		case "[]byte":
			typeValue = "bytes"
			outLine = fmt.Sprintf("\tproto%s.%s = %s.%s", ep.Camel, column.ColumnName.Camel, ep.Abbr, column.ColumnName.Camel)
			inLine = fmt.Sprintf("\t%s.%s = in.%s", ep.Abbr, column.ColumnName.Camel, column.ColumnName.Camel)
		case "*json.RawMessage":
			typeValue = "bytes"
			outLine = fmt.Sprintf("\tproto%s.%s, _ = %s.%s.MarshalJSON()", ep.Camel, column.ColumnName.Camel, ep.Abbr, column.ColumnName.Camel)
			inLine = fmt.Sprintf("\t%s.%s.UnmarshalJSON(in.%s)", ep.Abbr, column.ColumnName.Camel, column.ColumnName.Camel)
		default:
			typeValue = "string"
			outLine = fmt.Sprintf("\tproto%s.%s = %s.%s", ep.Camel, column.ColumnName.Camel, ep.Abbr, column.ColumnName.Camel)
			inLine = fmt.Sprintf("\t%s.%s = in.%s", ep.Abbr, column.ColumnName.Camel, column.ColumnName.Camel)
			if column.GoType == "null.String" {
				outLine = fmt.Sprintf("\tproto%s.%s = %s.%s.String", ep.Camel, column.ColumnName.Camel, ep.Abbr, column.ColumnName.Camel)
				inLine = fmt.Sprintf("\t%s.%s.Scan(in.%s)", ep.Abbr, column.ColumnName.Camel, column.ColumnName.Camel)
			}
			if column.GoType == "null.Time" {
				outLine = fmt.Sprintf("\tproto%s.%s = %s.%s.Time.Format(time.RFC3339)", ep.Camel, column.ColumnName.Camel, ep.Abbr, column.ColumnName.Camel)
			}
			if column.GoType == "null.Time" {
				inLine = fmt.Sprintf("\t%s.%s.Scan(in.%s)", ep.Abbr, column.ColumnName.Camel, column.ColumnName.Camel)
			}
		}
		lines = append(lines, fmt.Sprintf("\t%s %s = %d;", typeValue, column.ColumnName.Camel, idx))
		translateInLines = append(translateInLines, inLine)
		translateOutLines = append(translateOutLines, outLine)
	}
	lines = append(lines, "}")
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("message %sResponse {", ep.Camel))
	lines = append(lines, fmt.Sprintf("\t%s %s = 1;", ep.Name.Camel, ep.Name.Camel))
	lines = append(lines, fmt.Sprintf("\tResult result = 2;"))
	lines = append(lines, "}")
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("message %sRepeatResponse {", ep.Camel))
	lines = append(lines, fmt.Sprintf("\trepeated %s %s = 1;", ep.Name.Camel, ep.Name.Camel))
	lines = append(lines, fmt.Sprintf("\tResult result = 2;"))
	lines = append(lines, "}")
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("service %sService {", ep.Name.Camel))
	lines = append(lines, fmt.Sprintf("\trpc Get%s(IDIn) returns (%sResponse);", ep.Name.Camel, ep.Name.Camel))
	lines = append(lines, fmt.Sprintf("\trpc Search%s(%s) returns (%sRepeatResponse);", ep.Name.Camel, ep.Name.Camel, ep.Name.Camel))
	lines = append(lines, fmt.Sprintf("\trpc Post%s(%s) returns (%sResponse);", ep.Name.Camel, ep.Name.Camel, ep.Name.Camel))
	lines = append(lines, fmt.Sprintf("\trpc Put%s(%s) returns (Result);", ep.Name.Camel, ep.Name.Camel))
	lines = append(lines, fmt.Sprintf("\trpc Patch%s(%s) returns (Result);", ep.Name.Camel, ep.Name.Camel))
	lines = append(lines, fmt.Sprintf("\trpc Delete%s(IDIn) returns (Result);", ep.Name.Camel))
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
	// save off translateIn/Out
	ep.GrpcTranslateIn = strings.Join(translateInLines, "\n\t")
	ep.GrpcTranslateOut = strings.Join(translateOutLines, "\n\t")
}
