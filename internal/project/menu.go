package project

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	c "github.com/blackflagsoftware/forge-go/internal/column"
	e "github.com/blackflagsoftware/forge-go/internal/entity"
	n "github.com/blackflagsoftware/forge-go/internal/name"
	pf "github.com/blackflagsoftware/forge-go/internal/projectfile"
	s "github.com/blackflagsoftware/forge-go/internal/sql"
	"github.com/blackflagsoftware/forge-go/internal/util"
)

func (p *Project) ProjectMenu() {
	if p.ProjectFile.Storage == "s" {
		p.SqlMenu()
	} else {
		p.NonSqlMenu()
	}
	// run templating process on endpoints
	p.StartTemplating()
}

func (p *Project) NonSqlMenu() {
	messages := []string{"Field Type:"}
	prompts := []string{"(1) String", "(2) Integer", "(3) Decimal", "(4) Timestamp", "(5) Boolean", "(6) UUID"}
	acceptablePrompts := []string{"1", "2", "3", "4", "5", "6"}

	for {
		util.ClearScreen()
		fmt.Println("** File/MongoDB Storage Menu **")
		fmt.Println("")
		fmt.Print("Enter name of your object (e) to exit: ")
		name := n.Name{}
		rawName := util.ParseInput()
		if strings.ToLower(rawName) == "e" {
			break
		}
		name.BuildName(rawName, p.ProjectFile.KnownAliases)
		entity := e.Entity{Name: name}
		for {
			s.PrintSqlColumns(entity.Columns)
			column := c.Column{}
			fmt.Print("Field Name: (e) to exit")
			name := util.ParseInput()
			if strings.ToLower(name) == "e" {
				break
			}
			column.ColumnName.BuildName(name, p.ProjectFile.KnownAliases)
			selection := util.BasicPrompt(messages, prompts, acceptablePrompts, "e", util.ClearScreen)

			switch selection {
			case "1":
				column.GoType = "string"
			case "2":
				column.GoType = "int"
			case "3":
				column.GoType = "float64"
			case "4":
				column.GoType = "time.Time"
			case "5":
				column.GoType = "bool"
			case "6":
				column.GoType = "string"
			case "e", "E":
				break
			}
			entity.Columns = append(entity.Columns, column)
			anotherColumn := util.AskYesOrNo("Add another field?")
			if !anotherColumn {
				break
			}
		}
		p.Entities = append(p.Entities, entity)
		anotherEndpoint := util.AskYesOrNo("Add another object?")
		if !anotherEndpoint {
			break
		}
	}
}

func (p *Project) SqlMenu() {
	for {
		util.ClearScreen()
		mainMesssge := []string{"** SQL Storage Menu **", "How would you like create your entity?", "", "Project settings:"}
		mainMesssge = append(mainMesssge, fmt.Sprintf("Storage: %s", pf.StorageTypeToProper(p.ProjectFile.Storage)))
		if p.ProjectFile.Storage == "s" {
			mainMesssge = append(mainMesssge, fmt.Sprintf("  SQL Engine: %s", pf.SqlTypeToProper(p.ProjectFile.SqlStorage)))
			mainMesssge = append(mainMesssge, fmt.Sprintf("  Use ORM: %t", p.ProjectFile.UseORM))
		}
		mainMesssge = append(mainMesssge, fmt.Sprintf("TagForamt: %s", pf.TagFormatToProper(p.ProjectFile.TagFormat)))
		prompts := []string{"(1) File as input", "(2) Paste as input", "(3) Prompt as input", "(4) Blank Struct", "(5) Admin Screen"}
		acceptablePrompts := []string{"1", "2", "3", "4", "5"}
		selection := util.BasicPrompt(mainMesssge, prompts, acceptablePrompts, "e", util.ClearScreen)
		if selection == "e" {
			break
		}
		switch selection {
		case "1":
			p.FileMenu()
		case "2":
			p.PasteMenu()
		case "3":
			p.PromptMenu()
		case "4":
			p.BlankMenu()
		case "5":
			p.AdminMenu()
		}
	}
}

func (p *Project) FileMenu() {
	for {
		util.ClearScreen()
		fmt.Println("** File **")
		fmt.Println("")
		fmt.Print("Enter path to file or (e) to exit: ")
		selection := util.ParseInput()
		if strings.ToLower(selection) == "e" {
			return
		}
		filePath := filepath.Clean(selection)
		if _, err := os.Stat(filePath); !os.IsNotExist(err) {
			entities := processFile(p, filePath)
			if len(entities) > 0 {
				p.Entities = append(p.Entities, entities...)
			}
			return
		}
		fmt.Println("File does not exists, press 'enter' to continue")
		util.ParseInput()
	}
}

func (p *Project) PasteMenu() {
	for {
		sql := []string{}
		util.ClearScreen()
		fmt.Println("** Paste **")
		fmt.Println("")
		fmt.Print("Enter table sql schema or (e) to exit: ")
		for {
			line := util.ParseInput()
			if line == "" || strings.ToLower(line) == "e" {
				break
			} else if line[:] == ")" || line[:] == ");" || line[len(line)-1:] == ";" {
				sql = append(sql, line)
				break
			}
			sql = append(sql, line)
		}
		sqlEntity := s.ParseSqlLines(sql)
		entity := e.Entity{}
		entity.Name.BuildName(sqlEntity.Name, p.ProjectFile.KnownAliases)
		entity.Columns = sqlEntity.Columns
		entity.ColumnExistence = sqlEntity.ColExistence
		if sqlEntity.ColExistence.TimeColumn {
			entity.GrpcImport = "\"time\""
		}
		p.Entities = append(p.Entities, entity)
		cont := util.AskYesOrNo("Paste another table sql schema")
		if !cont {
			break
		}
	}
}

func (p *Project) PromptMenu() {
	for {
		sql := []string{}
		util.ClearScreen()
		fmt.Println("** Prompt **")
		fmt.Println("")
		fmt.Print("Enter entity name or (e) to exit: ")
		objectName := util.ParseInput()
		sql = append(sql, fmt.Sprintf("create table %s (", objectName))
		sql = append(sql, processColumns()...)
		sql = append(sql, ")")
		sqlEntity := s.ParseSqlLines(sql)
		entity := e.Entity{}
		entity.Name.BuildName(sqlEntity.Name, p.ProjectFile.KnownAliases)
		entity.Columns = sqlEntity.Columns
		entity.ColumnExistence = sqlEntity.ColExistence
		p.Entities = append(p.Entities, entity)
		cont := util.AskYesOrNo("Prompt for another entity")
		if !cont {
			break
		}
	}
}

func (p *Project) BlankMenu() {
	for {
		util.ClearScreen()
		fmt.Println("** Blank **")
		fmt.Println("")
		fmt.Print("Enter entity name or (e) to exit: ")
		objectName := util.ParseInput()
		if objectName == "e" {
			break
		}
		entity := e.Entity{}
		entity.Name.BuildName(objectName, p.ProjectFile.KnownAliases)
		p.Entities = append(p.Entities, entity)
		p.UseBlank = true
		cont := util.AskYesOrNo("Another blank entity")
		if !cont {
			break
		}
	}
}

func (p *Project) AdminMenu() {
	fmt.Println("Admin menu")
	util.ParseInput()
}

func processColumns() []string {
	messages := []string{"Column DB Type:"}
	prompts := []string{"(1) Varchar", "(2) Decimal", "(3) Integer", "(4) Timestamp", "(5) Boolean", "(6) Json", "(7) UUID", "(8) Auto Increment", "(9) Text", "(10) Char", "(11) Date"}
	acceptablePrompt := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "12"}
	sql := []string{}
	columns := []c.Column{}
	primaryKeys := []string{}
	for {
		util.ClearScreen()
		s.PrintSqlColumns(columns)
		col := c.Column{}
		fmt.Println("")
		fmt.Print("Enter Column Name or (e) to exit: ")
		name := util.ParseInput()
		if strings.ToLower(name) == "e" {
			break
		}
		col.ColumnName.BuildName(name, []string{})
		sel := util.BasicPrompt(messages, prompts, acceptablePrompt, "", util.ClearScreen)
		if strings.ToLower(sel) == "" {
			fmt.Println("Empty, try again!  Press enter to continue")
			util.ParseInput()
			continue
		}
		if util.AskYesOrNo("Is this a primary or part of composite/primary key") {
			primaryKeys = append(primaryKeys, col.ColumnName.RawName)
		}
		switch sel {
		case "1":
			col.DBType = "varchar"
			col.Length = askLengthPrompt(fmt.Sprintf("What is the %s length?", col.DBType))
			col.DBType = fmt.Sprintf("%s(%d)", col.DBType, col.Length)
			col.Null, col.DefaultValue = askNullDefaultPrompt(true, true)
		case "2":
			col.DBType = "numeric"
			col.Null, col.DefaultValue = askNullDefaultPrompt(true, true)
		case "3":
			col.DBType = "int"
			col.Null, col.DefaultValue = askNullDefaultPrompt(true, true)
		case "4":
			col.DBType = "timestamp"
			col.Null, col.DefaultValue = askNullDefaultPrompt(true, true)
		case "5":
			col.DBType = "bool"
			col.Null, col.DefaultValue = askNullDefaultPrompt(true, true)
		case "6":
			col.DBType = "json"
			col.Null, col.DefaultValue = askNullDefaultPrompt(false, true)
		case "7":
			col.DBType = "uuid"
		case "8":
			col.DBType = "autoincrement"
		case "9":
			col.DBType = "text"
			col.Null, col.DefaultValue = askNullDefaultPrompt(true, false)
		case "10":
			col.DBType = "char"
			col.Length = askLengthPrompt(fmt.Sprintf("What is the %s length?", col.DBType))
			col.DBType = fmt.Sprintf("%s(%d)", col.DBType, col.Length)
			col.Null, col.DefaultValue = askNullDefaultPrompt(true, true)
		case "11":
			col.DBType = "date"
			col.Null, col.DefaultValue = askNullDefaultPrompt(true, false)
		default:
			fmt.Println("Not a valid selection, try again!")
			util.ParseInput()
			continue
		}
		sql = append(sql, buildSqlColumn(col))
		if !util.AskYesOrNo("Add another column") {
			break
		}
	}
	if len(primaryKeys) > 0 {
		sql = append(sql, fmt.Sprintf("primary key(%s)", strings.Join(primaryKeys, ", ")))
	}
	// add commas to end of all line except last one
	for i := 0; i < len(sql)-1; i++ {
		sql[i] = sql[i] + ","
	}
	return sql
}

func askLengthPrompt(msg string) int64 {
	for {
		fmt.Print(msg)
		length := util.ParseInput()
		lenInt, errParse := strconv.ParseInt(length, 10, 64)
		if errParse != nil {
			fmt.Println("Not an integer, you can can do better than that!")
			continue
		}
		return lenInt
	}
}

func askNullDefaultPrompt(askNull, askDefault bool) (nullAble bool, defaultValue string) {
	if askNull {
		nullAble = util.AskYesOrNo("Can this column be null")
	}
	if !nullAble && askDefault {
		if util.AskYesOrNo("Does this column have a default value") {
			fmt.Print("What is the default value? ")
			defaultValue = util.ParseInput()
		}
	}
	return
}

func buildSqlColumn(col c.Column) string {
	defaultStmt := " default null"
	if !col.Null {
		defaultStmt = " not null"
		if col.DefaultValue != "" {
			defaultStmt = fmt.Sprintf(" default '%s'", col.DefaultValue)
		}
	}
	return fmt.Sprintf("%s %s%s", col.ColumnName.RawName, col.DBType, defaultStmt)
}

func processFile(p *Project, filePath string) (entities []e.Entity) {
	bContent, errRead := ioutil.ReadFile(filePath)
	if errRead != nil {
		fmt.Println("processFile - error:", errRead)
		return
	}
	bArray := bytes.Split(bContent, []byte("\n"))
	arraySqlStmt := [][]string{}
	sqlStmt := []string{}
	tableCount := 0
	for _, bLine := range bArray {
		if bytes.Contains(bytes.ToLower(bLine), []byte("create")) {
			// found an create stmt, add all lines already acquired and add to arraySql if not empty
			if len(sqlStmt) > 0 {
				tableCount++
				arraySqlStmt = append(arraySqlStmt, sqlStmt)
				sqlStmt = []string{string(bLine)}
				continue
			}
			// else just add if not empty
			if len(bytes.TrimSpace(bLine)) != 0 {
				sqlStmt = append(sqlStmt, string(bLine))
			}
		}
	}
	// add the last set of lines if not empty
	if len(sqlStmt) > 0 {
		tableCount++
		arraySqlStmt = append(arraySqlStmt, sqlStmt)
	}
	fmt.Println("")
	fmt.Printf("Processed %d tables, press any key to continue or (e) to exit", tableCount)
	cont := util.ParseInput()
	if strings.ToLower(cont) == "e" {
		return
	}
	for i := range arraySqlStmt {
		sqlEntity := s.ParseSqlLines(arraySqlStmt[i])
		entity := e.Entity{}
		entity.Name.BuildName(sqlEntity.Name, p.ProjectFile.KnownAliases)
		entity.Columns = sqlEntity.Columns
		entity.ColumnExistence = sqlEntity.ColExistence
		entities = append(entities, entity)
	}
	return
}
