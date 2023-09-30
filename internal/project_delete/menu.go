package project

// import (
// 	"bytes"
// 	"fmt"
// 	"io/ioutil"
// 	"os"
// 	"path/filepath"
// 	"strconv"
// 	"strings"

// 	con "github.com/blackflagsoftware/forge-go/internal/constant"
// 	e "github.com/blackflagsoftware/forge-go/internal/entity"
// 	m "github.com/blackflagsoftware/forge-go/internal/model"
// 	n "github.com/blackflagsoftware/forge-go/internal/name"
// 	pf "github.com/blackflagsoftware/forge-go/internal/projectfile"
// 	s "github.com/blackflagsoftware/forge-go/internal/sql"
// 	"github.com/blackflagsoftware/forge-go/internal/util"
// )

// func (p *Project) ProjectMenu() {
// 	if p.ProjectFile.Storage == "s" {
// 		p.SqlMenu()
// 	} else {
// 		p.NonSqlMenu()
// 	}
// }

// func (p *Project) NonSqlMenu() {
// 	messages := []string{"Field Type:"}
// 	prompts := []string{"(1) String", "(2) Integer", "(3) Decimal", "(4) Timestamp", "(5) Boolean", "(6) UUID"}
// 	acceptablePrompts := []string{"1", "2", "3", "4", "5", "6"}

// OuterLoop:
// 	for {
// 		util.ClearScreen()
// 		fmt.Println("** File/MongoDB Storage Menu **")
// 		fmt.Println("")
// 		fmt.Print("Enter name of your entity (e) to exit: ")
// 		name := n.Name{}
// 		rawName := util.ParseInput()
// 		if strings.ToLower(rawName) == "e" {
// 			break
// 		}
// 		entityName := name.BuildName(rawName, p.ProjectFile.KnownAliases)
// 		p.ProjectFile.KnownAliases = append(p.ProjectFile.KnownAliases, entityName)
// 		entity := e.Entity{Name: name}
// 		for {
// 			s.PrintSqlColumns(entity.Columns)
// 			column := m.Column{}
// 			fmt.Print("Field Name (e) to exit: ")
// 			name := util.ParseInput()
// 			if strings.ToLower(name) == "e" {
// 				break OuterLoop
// 			}
// 			column.ColumnName.BuildName(name, []string{})
// 			selection := util.BasicPrompt(messages, prompts, acceptablePrompts, "e", util.ClearScreen)

// 			switch selection {
// 			case "1":
// 				column.GoType = "null.String"
// 			case "2":
// 				column.GoType = "null.Int"
// 				column.DBType = "int" // need this in grpc_template.go
// 			case "3":
// 				column.GoType = "null.Float"
// 			case "4":
// 				column.GoType = "null.Time"
// 			case "5":
// 				column.GoType = "null.Bool"
// 			case "6":
// 				column.GoType = "string"
// 			case "e", "E":
// 				break OuterLoop
// 			}
// 			if util.AskYesOrNo("Is this a primary or part of composite/primary key") {
// 				column.PrimaryKey = true
// 			}
// 			entity.Columns = append(entity.Columns, column)
// 			anotherColumn := util.AskYesOrNo("Add another field?")
// 			if !anotherColumn {
// 				break
// 			}
// 		}
// 		p.Entities = append(p.Entities, entity)
// 		anotherEndpoint := util.AskYesOrNo("Add another entity?")
// 		if !anotherEndpoint {
// 			break
// 		}
// 	}
// 	p.StartTemplating()
// 	p.Entities = []e.Entity{}
// 	fmt.Println("")
// 	fmt.Println("Entities have been processed, press 'enter' to continue")
// 	util.ParseInput()
// }

// func (p *Project) SqlMenu() {
// 	for {
// 		util.ClearScreen()
// 		mainMesssge := []string{"** SQL Storage Menu **", "How would you like create your entity?", "", "Project settings:"}
// 		mainMesssge = append(mainMesssge, fmt.Sprintf("Storage: %s", pf.StorageTypeToProper(p.ProjectFile.Storage)))
// 		if p.ProjectFile.Storage == "s" {
// 			mainMesssge = append(mainMesssge, fmt.Sprintf("  SQL Engine: %s", pf.SqlTypeToProper(p.ProjectFile.SqlStorage)))
// 			mainMesssge = append(mainMesssge, fmt.Sprintf("  Use ORM: %t", p.ProjectFile.UseORM))
// 		}
// 		mainMesssge = append(mainMesssge, fmt.Sprintf("TagForamt: %s", pf.TagFormatToProper(p.ProjectFile.TagFormat)))
// 		prompts := []string{"(1) File as input", "(2) Paste as input", "(3) Prompt as input", "(4) Blank Struct", "(5) Admin Screen"}
// 		acceptablePrompts := []string{"1", "2", "3", "4", "5"}
// 		selection := util.BasicPrompt(mainMesssge, prompts, acceptablePrompts, "e", util.ClearScreen)
// 		if selection == "e" {
// 			break
// 		}
// 		switch selection {
// 		case "1":
// 			p.FileMenu()
// 		case "2":
// 			p.PasteMenu()
// 		case "3":
// 			p.PromptMenu()
// 		case "4":
// 			p.BlankMenu()
// 		case "5":
// 			p.AdminMenu()
// 		}
// 		if len(p.Entities) > 0 {
// 			p.StartTemplating()
// 			// remove entities already processed
// 			p.Entities = []e.Entity{}
// 			fmt.Println("")
// 			fmt.Println("Entities have been processed, press 'enter' to continue")
// 			util.ParseInput()
// 		}
// 	}
// }

// func (p *Project) FileMenu() {
// 	for {
// 		util.ClearScreen()
// 		fmt.Println("** File **")
// 		fmt.Println("")
// 		fmt.Print("Enter full path to file or (e) to exit: ")
// 		selection := util.ParseInput()
// 		if strings.ToLower(selection) == "e" {
// 			return
// 		}
// 		filePath := filepath.Clean(selection)
// 		if _, err := os.Stat(filePath); !os.IsNotExist(err) {
// 			entities := processFile(p, filePath)
// 			if len(entities) > 0 {
// 				p.Entities = append(p.Entities, entities...)
// 			}
// 			return
// 		}
// 		fmt.Println("File does not exists, press 'enter' to continue")
// 		util.ParseInput()
// 	}
// }

// func (p *Project) PasteMenu() {
// PasteLoop:
// 	for {
// 		sql := []string{}
// 		util.ClearScreen()
// 		fmt.Println("** Paste **")
// 		fmt.Println("")
// 		fmt.Println("Enter table sql schema or (e) to exit:")
// 		fmt.Println("")
// 		for {
// 			line := util.ParseInput()
// 			if line == "" || strings.ToLower(line) == "e" {
// 				break PasteLoop
// 			} else if line[:] == ")" || line[:] == ");" || line[len(line)-1:] == ";" {
// 				sql = append(sql, line)
// 				break
// 			}
// 			sql = append(sql, line)
// 		}
// 		sqlEntity, err := s.ParseSqlLines(sql)
// 		if err != nil {
// 			return
// 		}
// 		entity := e.Entity{}
// 		name := entity.Name.BuildName(sqlEntity.Name, p.ProjectFile.KnownAliases)
// 		p.ProjectFile.KnownAliases = append(p.ProjectFile.KnownAliases, name)
// 		entity.Columns = sqlEntity.Columns
// 		entity.ColumnExistence = sqlEntity.ColExistence
// 		if sqlEntity.ColExistence.TimeColumn {
// 			entity.GrpcImport = "\"time\""
// 		}
// 		p.Entities = append(p.Entities, entity)
// 		cont := util.AskYesOrNo("Paste another table sql schema")
// 		if !cont {
// 			break
// 		}
// 	}
// }

// func (p *Project) PromptMenu() {
// 	for {
// 		sql := []string{}
// 		util.ClearScreen()
// 		fmt.Println("** Prompt **")
// 		fmt.Println("")
// 		fmt.Print("Enter entity name or (e) to exit: ")
// 		entityName := util.ParseInput()
// 		if strings.ToLower(entityName) == "e" {
// 			break
// 		}
// 		sql = append(sql, fmt.Sprintf("create table %s (", entityName))
// 		sql = append(sql, processColumns()...)
// 		sql = append(sql, ")")
// 		sqlEntity, err := s.ParseSqlLines(sql)
// 		if err != nil {
// 			return
// 		}
// 		entity := e.Entity{}
// 		entity.Name.BuildName(sqlEntity.Name, p.ProjectFile.KnownAliases)
// 		entity.Columns = sqlEntity.Columns
// 		entity.ColumnExistence = sqlEntity.ColExistence
// 		p.Entities = append(p.Entities, entity)
// 		cont := util.AskYesOrNo("Prompt for another entity")
// 		if !cont {
// 			break
// 		}
// 	}
// 	p.saveOutSql()
// }

// func (p *Project) BlankMenu() {
// 	for {
// 		util.ClearScreen()
// 		fmt.Println("** Blank **")
// 		fmt.Println("")
// 		fmt.Print("Enter entity name or (e) to exit: ")
// 		entityName := util.ParseInput()
// 		if entityName == "e" {
// 			break
// 		}
// 		entity := e.Entity{}
// 		name := entity.Name.BuildName(entityName, p.ProjectFile.KnownAliases)
// 		p.ProjectFile.KnownAliases = append(p.ProjectFile.KnownAliases, name)
// 		p.Entities = append(p.Entities, entity)
// 		p.UseBlank = true
// 		cont := util.AskYesOrNo("Another blank entity")
// 		if !cont {
// 			break
// 		}
// 	}
// }

// func (p *Project) AdminMenu() {
// 	mainMessage := []string{"** Admin Menu **", "", "Please make a selection:"}
// 	prompts := []string{"(1) Change Storage Type", "(2) Change TagFormat", "(3) Add a module", "(4) SQL only - use an ORM"}
// 	acceptablePrompts := []string{"1", "2", "3", "4"}
// 	for {
// 		util.ClearScreen()
// 		sel := util.BasicPrompt(mainMessage, prompts, acceptablePrompts, "e", util.ClearScreen)
// 		if sel == "e" {
// 			break
// 		}
// 		switch sel {
// 		case "1":
// 			p.StorageMenu()
// 		case "2":
// 			p.TagFormatMenu()
// 		case "3":
// 			p.ModuleMenu()
// 		case "4":
// 			p.OrmMenu()
// 		}
// 		// fmt.Println("temp 'enter' here")
// 		// util.ParseInput()
// 	}
// }

// func (p *Project) StorageMenu() {
// 	util.ClearScreen()
// 	mainMesssge := []string{"Storage Type", fmt.Sprintf("Current Value: %s", pf.StorageTypeToProper(p.ProjectFile.Storage)), "Do you wish to change the Storage Type?"}
// 	prompts := []string{"(s) SQL", "(f) File", "(m) MongoDB"}
// 	acceptablePrompts := []string{"s", "f", "m"}
// 	response := util.BasicPrompt(mainMesssge, prompts, acceptablePrompts, "e", util.ClearScreen)
// 	if response == "e" {
// 		return
// 	}
// 	p.ProjectFile.Storage = response
// 	if p.ProjectFile.Storage == "s" {
// 		mainMesssge = []string{"SQL Option", "Choice which SQL implementation"}
// 		prompts = []string{"(p) Postgres", "(m) Mysql", "(s) Sqlite"}
// 		acceptablePrompts = []string{"p", "m", "s"} // haha... get it?!??!!  pms... haha
// 		p.ProjectFile.SqlStorage = util.BasicPrompt(mainMesssge, prompts, acceptablePrompts, "", util.ClearScreen)
// 		// p.UseORM = util.AskYesOrNo("Would you like to use an ORM") only have this as an option in an "admin" screen
// 	}
// }

// func (p *Project) TagFormatMenu() {
// 	util.ClearScreen()
// 	mainMesssge := []string{"Tag Format", fmt.Sprintf("Current Value: %s", pf.TagFormatToProper(p.ProjectFile.TagFormat)), "Do you wish to change the Tag Format?"}
// 	prompts := []string{"(s) Snake Case (tag_format)", "(c) Camel Case (tagFormat)", "(p) Pascal Case (TagFormat)", "(k) Kebab Case (tag-format)", "(l) Lower Case (tag format)", "(u) Upper (TAG FORMAT)"}
// 	acceptablePrompts := []string{"s", "c", "p", "k", "l", "u"}
// 	tagFormat := util.BasicPrompt(mainMesssge, prompts, acceptablePrompts, "e", util.ClearScreen)
// 	if tagFormat == "e" {
// 		return
// 	}
// 	switch tagFormat {
// 	case "s":
// 		p.ProjectFile.TagFormat = "snakeCase"
// 	case "k":
// 		p.ProjectFile.TagFormat = "kebabCase"
// 	case "c":
// 		p.ProjectFile.TagFormat = "camelCase"
// 	case "p":
// 		p.ProjectFile.TagFormat = "pascalCase"
// 	case "u":
// 		p.ProjectFile.TagFormat = "upperCase"
// 	case "l":
// 		p.ProjectFile.TagFormat = "lowerCase"
// 	}
// }

// func (p *Project) OrmMenu() {
// 	util.ClearScreen()
// 	mainMesssge := []string{"Storage Type", fmt.Sprintf("Current Value: %t", p.ProjectFile.UseORM), "Do you wish to change this value?"}
// 	prompts := []string{"(t) True", "(f) False"}
// 	acceptablePrompts := []string{"t", "f"}
// 	response := util.BasicPrompt(mainMesssge, prompts, acceptablePrompts, "e", util.ClearScreen)
// 	if response == "e" {
// 		return
// 	}
// 	p.ProjectFile.UseORM = true
// 	if response == "f" {
// 		p.ProjectFile.UseORM = false
// 	}
// }

// func (p *Project) ModuleMenu() {
// 	mainMessage := []string{"** Add Module **", "", "Which module do you wish to add?"}
// 	prompts := []string{"(1) Login"}
// 	acceptablePrompt := []string{"1"}
// 	for {
// 		util.ClearScreen()
// 		sel := util.BasicPrompt(mainMessage, prompts, acceptablePrompt, "e", util.ClearScreen)
// 		if sel == "e" {
// 			break
// 		}
// 		switch sel {
// 		case "1":
// 			p.ModuleAddLogin()
// 		}
// 		fmt.Println("temp 'enter' here")
// 		util.ParseInput()
// 	}
// }

// func processColumns() []string {
// 	messages := []string{"Column DB Type:"}
// 	prompts := []string{"(1) Varchar", "(2) Decimal", "(3) Integer", "(4) Timestamp", "(5) Boolean", "(6) Json", "(7) UUID", "(8) Auto Increment", "(9) Text", "(10) Char", "(11) Date"}
// 	acceptablePrompt := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "12"}
// 	sql := []string{}
// 	columns := []m.Column{}
// 	primaryKeys := []string{}
// 	for {
// 		util.ClearScreen()
// 		s.PrintSqlColumns(columns)
// 		col := m.Column{}
// 		fmt.Println("")
// 		fmt.Print("Enter Column Name or (e) to exit: ")
// 		name := util.ParseInput()
// 		if strings.ToLower(name) == "e" {
// 			break
// 		}
// 		col.ColumnName.BuildName(name, []string{})
// 		sel := util.BasicPrompt(messages, prompts, acceptablePrompt, "", util.ClearScreen)
// 		if strings.ToLower(sel) == "" {
// 			fmt.Println("Empty, try again!  Press enter to continue")
// 			util.ParseInput()
// 			continue
// 		}
// 		if util.AskYesOrNo("Is this a primary or part of composite/primary key") {
// 			primaryKeys = append(primaryKeys, col.ColumnName.RawName)
// 		}
// 		switch sel {
// 		case "1":
// 			col.DBType = "varchar"
// 			col.Length = askLengthPrompt(fmt.Sprintf("What is the %s length? ", col.DBType))
// 			col.DBType = fmt.Sprintf("%s(%d)", col.DBType, col.Length)
// 			col.Null, col.DefaultValue = askNullDefaultPrompt(true, true)
// 		case "2":
// 			col.DBType = "numeric"
// 			col.Null, col.DefaultValue = askNullDefaultPrompt(true, true)
// 		case "3":
// 			col.DBType = "int"
// 			col.Null, col.DefaultValue = askNullDefaultPrompt(true, true)
// 		case "4":
// 			col.DBType = "timestamp"
// 			col.Null, col.DefaultValue = askNullDefaultPrompt(true, true)
// 		case "5":
// 			col.DBType = "bool"
// 			col.Null, col.DefaultValue = askNullDefaultPrompt(true, true)
// 		case "6":
// 			col.DBType = "json"
// 			col.Null, col.DefaultValue = askNullDefaultPrompt(false, true)
// 		case "7":
// 			col.DBType = "uuid"
// 		case "8":
// 			col.DBType = "autoincrement"
// 		case "9":
// 			col.DBType = "text"
// 			col.Null, col.DefaultValue = askNullDefaultPrompt(true, false)
// 		case "10":
// 			col.DBType = "char"
// 			col.Length = askLengthPrompt(fmt.Sprintf("What is the %s length? ", col.DBType))
// 			col.DBType = fmt.Sprintf("%s(%d)", col.DBType, col.Length)
// 			col.Null, col.DefaultValue = askNullDefaultPrompt(true, true)
// 		case "11":
// 			col.DBType = "date"
// 			col.Null, col.DefaultValue = askNullDefaultPrompt(true, false)
// 		default:
// 			fmt.Println("Not a valid selection, try again!")
// 			util.ParseInput()
// 			continue
// 		}
// 		columns = append(columns, col)
// 		sql = append(sql, buildSqlColumn(col))
// 		if !util.AskYesOrNo("Add another column") {
// 			break
// 		}
// 	}
// 	if len(primaryKeys) > 0 {
// 		sql = append(sql, fmt.Sprintf("primary key(%s)", strings.Join(primaryKeys, ", ")))
// 	}
// 	// add commas to end of all line except last one
// 	for i := 0; i < len(sql)-1; i++ {
// 		sql[i] = sql[i] + ","
// 	}
// 	return sql
// }

// func askLengthPrompt(msg string) int64 {
// 	for {
// 		fmt.Print(msg)
// 		length := util.ParseInput()
// 		lenInt, errParse := strconv.ParseInt(length, 10, 64)
// 		if errParse != nil {
// 			fmt.Println("Not an integer, you can can do better than that!")
// 			continue
// 		}
// 		return lenInt
// 	}
// }

// func askNullDefaultPrompt(askNull, askDefault bool) (nullAble bool, defaultValue string) {
// 	if askNull {
// 		nullAble = util.AskYesOrNo("Can this column be null")
// 	}
// 	if !nullAble && askDefault {
// 		if util.AskYesOrNo("Does this column have a default value") {
// 			fmt.Print("What is the default value? ")
// 			defaultValue = util.ParseInput()
// 		}
// 	}
// 	return
// }

// func buildSqlColumn(col m.Column) string {
// 	defaultStmt := " default null"
// 	if !col.Null {
// 		defaultStmt = " not null"
// 		if col.DefaultValue != "" {
// 			defaultStmt = fmt.Sprintf(" default '%s'", col.DefaultValue)
// 		}
// 	}
// 	return fmt.Sprintf("%s %s%s", col.ColumnName.RawName, col.DBType, defaultStmt)
// }

// func processFile(p *Project, filePath string) (entities []e.Entity) {
// 	bContent, errRead := ioutil.ReadFile(filePath)
// 	if errRead != nil {
// 		fmt.Println("processFile - error:", errRead)
// 		return
// 	}
// 	bArray := bytes.Split(bContent, []byte("\n"))
// 	arraySqlStmt := [][]string{}
// 	sqlStmt := []string{}
// 	tableCount := 0
// 	for _, bLine := range bArray {
// 		if bytes.Contains(bytes.ToLower(bLine), []byte("create table")) {
// 			// found an create stmt, add all lines already acquired and add to arraySql if not empty
// 			if len(sqlStmt) > 0 {
// 				tableCount++
// 				arraySqlStmt = append(arraySqlStmt, sqlStmt)
// 				sqlStmt = []string{string(bLine)}
// 				continue
// 			}
// 		}
// 		// else just add if not empty
// 		if len(bytes.TrimSpace(bLine)) != 0 {
// 			sqlStmt = append(sqlStmt, string(bLine))
// 		}
// 	}
// 	// add the last set of lines if not empty
// 	if len(sqlStmt) > 0 {
// 		tableCount++
// 		arraySqlStmt = append(arraySqlStmt, sqlStmt)
// 	}
// 	fmt.Println("")
// 	fmt.Printf("Processed %d tables, press any key to continue or (e) to exit", tableCount)
// 	cont := util.ParseInput()
// 	if strings.ToLower(cont) == "e" {
// 		return
// 	}
// 	for i := range arraySqlStmt {
// 		sqlEntity, err := s.ParseSqlLines(arraySqlStmt[i])
// 		if err != nil {
// 			return
// 		}
// 		entity := e.Entity{}
// 		name := entity.Name.BuildName(sqlEntity.Name, p.ProjectFile.KnownAliases)
// 		p.ProjectFile.KnownAliases = append(p.ProjectFile.KnownAliases, name)
// 		if sqlEntity.ColExistence.TimeColumn {
// 			entity.GrpcImport = "\"time\""
// 		}
// 		entity.Columns = sqlEntity.Columns
// 		entity.ColumnExistence = sqlEntity.ColExistence
// 		entities = append(entities, entity)
// 	}
// 	return
// }

// func (p *Project) saveOutSql() {
// 	sqlProvider := ""
// 	switch p.ProjectFile.Storage {
// 	case "m":
// 		sqlProvider = con.MYSQL
// 	case "p":
// 		sqlProvider = con.POSTGRESQL
// 	case "s":
// 		sqlProvider = con.SQLITE3
// 	}
// 	fileName := "./prompt_schema"
// 	lines := []string{}
// 	for e, ep := range p.Entities {
// 		primaryKeys := []string{}
// 		lines = append(lines, fmt.Sprintf("create table if not exists %s (", ep.Name.Lower))
// 		for i, c := range ep.Columns {
// 			null := " null"
// 			defaultValue := ""
// 			length := ""
// 			if !c.Null {
// 				null = " not null"
// 			}
// 			dbType := c.DBType
// 			if dbType == "autoincrement" || dbType == "serial" {
// 				null = ""
// 				if sqlProvider == con.SQLITE3 {
// 					dbType = "integer primary key autoincrement"
// 				}
// 				if sqlProvider == con.MYSQL {
// 					dbType = "integer auto_increment"
// 				}
// 				if sqlProvider == con.POSTGRESQL {
// 					dbType = "serial"
// 				}
// 			}
// 			if c.PrimaryKey && !(dbType == "autoincrement" && ep.SQLProvider == con.SQLITE3) {
// 				primaryKeys = append(primaryKeys, c.ColumnName.Lower)
// 			}
// 			if c.DefaultValue != "" {
// 				if c.DBType == "varchar" || c.DBType == "char" || c.DBType == "text" {
// 					defaultValue = fmt.Sprintf(" default '%s'", c.DefaultValue)
// 				} else {
// 					defaultValue = fmt.Sprintf(" default %s", c.DefaultValue)
// 				}
// 			}
// 			if c.Length > 0 {
// 				length = fmt.Sprintf("(%d)", c.Length)
// 			}
// 			if i < len(ep.Columns)-1 || (i == len(ep.Columns)-1 && len(primaryKeys) > 0) {
// 				lines = append(lines, fmt.Sprintf("\t%s %s%s%s%s,", c.ColumnName.Lower, dbType, length, null, defaultValue))
// 			} else {
// 				lines = append(lines, fmt.Sprintf("\t%s %s%s%s%s", c.ColumnName.Lower, dbType, length, null, defaultValue))
// 			}
// 		}
// 		if sqlProvider != con.SQLITE3 && len(primaryKeys) > 0 {
// 			lines = append(lines, fmt.Sprintf("\tprimary key(%s)", strings.Join(primaryKeys, ", ")))
// 		}
// 		lines = append(lines, ");")
// 		if e > 0 {
// 			lines = append(lines, "")
// 		}
// 	}
// 	lines = append(lines, "\n")
// 	if len(lines) > 0 {
// 		file, errOpen := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
// 		if errOpen != nil {
// 			fmt.Println("Unable to save schema to:", fileName)
// 			return
// 		}
// 		defer file.Close()
// 		if _, errWrite := file.WriteString(strings.Join(lines, "\n")); errWrite != nil {
// 			fmt.Println("Unable to write lines to:", fileName)
// 		}
// 	}
// }

// func (p *Project) ModuleAddLogin() {
// 	for _, m := range p.ProjectFile.Modules {
// 		if m == "login" {
// 			fmt.Println("The 'Login' module has already been added to this project, press 'enter' to continue")
// 			util.ParseInput()
// 			return
// 		}
// 	}
// 	// add all the modules
// 	p.AddLogin()
// 	// mark it as complete
// 	p.ProjectFile.Modules = append(p.ProjectFile.Modules, "login")
// }
