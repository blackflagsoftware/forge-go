package menu

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	m "github.com/blackflagsoftware/forge-go/internal/model"
	"github.com/blackflagsoftware/forge-go/internal/util"
)

func NonSqlMenu(p *m.Project) {
	messages := []string{"Action:"}
	prompts := []string{"(1) New Field", "(2) Edit Field", "(3) Delete Field", "(4) Save & Exit", "(5) Cancel & Exit"}
	acceptablePrompts := []string{"1", "2", "3", "4", "5"}

	for {
		util.ClearScreen()
		fmt.Println("** File/MongoDB Storage Menu **")
		fmt.Println("")
		fmt.Print("Enter name of your entity (e) to exit: ")
		name := m.Name{}
		rawName := util.ParseInput()
		if strings.ToLower(rawName) == "e" {
			break
		}
		entityName := name.BuildName(rawName, p.ProjectFile.KnownAliases)
		p.ProjectFile.KnownAliases = append(p.ProjectFile.KnownAliases, entityName)
		p.SaveProjectFile()
		entity := m.Entity{Name: name}
		cancelExit := false
		for {
			saveExit := false
			util.ClearScreen()
			printSavedColumns(entity.Columns)
			fmt.Println("")
			selection := util.BasicPrompt(messages, prompts, acceptablePrompts, "", nil)
			switch selection {
			case "1":
				// add
				column := newField()
				entity.Columns = append(entity.Columns, column)
			case "2":
				// edit
				editField(&entity)
			case "3":
				// delete
				deleteField(&entity)
			case "4":
				saveExit = true
			case "5":
				cancelExit = true
			}
			if saveExit || cancelExit {
				break
			}
		}
		if !cancelExit {
			p.Entities = append(p.Entities, entity)
			anotherEndpoint := util.AskYesOrNo("Add another entity?")
			if !anotherEndpoint {
				break
			}
		}
	}
}

func newField() m.Column {
	messages := []string{"Field Type:"}
	prompts := []string{"(1) String", "(2) Integer", "(3) Decimal", "(4) Timestamp", "(5) Boolean", "(6) UUID"}
	acceptablePrompts := []string{"1", "2", "3", "4", "5", "6"}
	column := m.Column{}
	fmt.Print("Field Name (e) to exit: ")
	name := util.ParseInput()
	if strings.ToLower(name) == "e" {
		return column
	}
	column.ColumnName.BuildName(name, []string{})
	selection := util.BasicPrompt(messages, prompts, acceptablePrompts, "e", util.ClearScreen)

	switch selection {
	case "1":
		column.GoType = "null.String"
		column.DBType = "string"
		column.GoTypeNonSql = "string"
	case "2":
		column.GoType = "null.Int"
		column.DBType = "int"
		column.GoTypeNonSql = "int"
	case "3":
		column.GoType = "null.Float"
		column.DBType = "float"
		column.GoTypeNonSql = "float"
	case "4":
		column.GoType = "null.Time"
		column.DBType = "time"
		column.GoTypeNonSql = "string"
		// entity.GrpcImport = "\"time\""
	case "5":
		column.GoType = "null.Bool"
		column.DBType = "boolean"
		column.GoTypeNonSql = "boolean"
	case "6":
		column.GoType = "string"
		column.DBType = "string"
		column.GoTypeNonSql = "string"
		column.IsNoSQLUUID = true
	case "e", "E":
		return column
	}
	if util.AskYesOrNo("Is this a primary or part of composite/primary key") {
		column.PrimaryKey = true
	}
	return column
}

func editField(entity *m.Entity) {
	types := map[string]string{"string": "String", "int": "Integer", "float": "Decimal", "time": "Timestamp", "boolean": "Boolean", "uuid": "UUID"}
	typeIdx := map[string]int{"string": 1, "int": 2, "float": 3, "time": 4, "boolean": 5, "uuid": 6}
	for {
		fmt.Println("Enter # to enter 'edit' mode:")
		fmt.Println("")
		idx := printColumns(entity.Columns)
		if idx == -1 {
			break
		}
		// make new column and do a deep copy later
		newColumn := m.Column{}
		fmt.Println("")
		fmt.Println("Defaults in '[]', press 'enter' accept default")
		newName := util.ParseInputStringWithMessageCompare(fmt.Sprintf("Name [%s]: ", entity.Columns[idx].ColumnName.RawName), entity.Columns[idx].ColumnName.RawName)
		newColumn.ColumnName.BuildName(newName, []string{})
		dbType := entity.Columns[idx].DBType
		if dbType == "string" && entity.Columns[idx].IsNoSQLUUID {
			dbType = "uuid"
		}
		displayType := types[dbType]
		message := fmt.Sprintf("(1) String, (2) Integer, (3) Decimal, (4) Timestamp,  (5) Boolean, (6) UUID\nField Type [%s]: ", displayType)
		existingValue := typeIdx[dbType]
		newType := util.ParseInputIntWithMessageCompare(message, existingValue)

		switch newType {
		case 1:
			newColumn.GoType = "null.String"
			newColumn.DBType = "string"
			newColumn.GoTypeNonSql = "string"
		case 2:
			newColumn.GoType = "null.Int"
			newColumn.DBType = "int"
			newColumn.GoTypeNonSql = "int"
		case 3:
			newColumn.GoType = "null.Float"
			newColumn.DBType = "float"
			newColumn.GoTypeNonSql = "float"
		case 4:
			newColumn.GoType = "null.Time"
			newColumn.DBType = "time"
			newColumn.GoTypeNonSql = "string"
		case 5:
			newColumn.GoType = "null.Bool"
			newColumn.DBType = "boolean"
			newColumn.GoTypeNonSql = "boolean"
		case 6:
			newColumn.GoType = "string"
			newColumn.DBType = "string"
			newColumn.GoTypeNonSql = "string"
			newColumn.IsNoSQLUUID = true
		}
		newColumn.PrimaryKey = util.ParseInputBoolWithMessageCompare(fmt.Sprintf("Primary Key (y/n) [%t]: ", entity.Columns[idx].PrimaryKey), entity.Columns[idx].PrimaryKey)
		// replace the column
		entity.Columns[idx] = newColumn
	}
}

func deleteField(entity *m.Entity) {
	for {
		fmt.Println("Enter # to enter 'edit' mode:")
		fmt.Println("")
		idx := printColumns(entity.Columns)
		if idx == -1 {
			break
		}
		entity.Columns = append(entity.Columns[:idx], entity.Columns[idx+1:]...)
	}
}

func printColumns(columns []m.Column) int {
	for {
		for i, c := range columns {
			printIdx := i + 1
			fmt.Printf("%d - %s\n", printIdx, c.ColumnName.Camel)
		}
		fmt.Println("e - exit")
		selection := util.ParseInput()
		if strings.ToLower(selection) == "e" {
			return -1 // send back signal to exit
		}
		selectionIdx, err := strconv.Atoi(selection)
		if err != nil {
			fmt.Println("Not a valid number, press 'enter' to continue")
			util.ParseInput()
			util.ClearScreen()
			continue
		}
		if selectionIdx < 1 || selectionIdx > len(columns) {
			fmt.Println("Not in valid range, press 'enter' to continue")
			util.ParseInput()
			util.ClearScreen()
			continue
		}
		return selectionIdx - 1 // go back down to zero based array
	}
}

func printSavedColumns(cols []m.Column) {
	fmt.Println("--- Saved Columns ---")
	tab := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	fmt.Fprintln(tab, "|Name\t|Type\t|Primary Key")
	fmt.Fprintln(tab, "+----\t+----\t+-----------")
	for _, col := range cols {
		fmt.Fprintf(tab, " %s\t %s\t %t\n", col.ColumnName.Camel, col.DBType, col.PrimaryKey)
	}
	tab.Flush()
	fmt.Println("")
}
