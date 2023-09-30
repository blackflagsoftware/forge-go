package menu

import (
	"fmt"
	"os"
	"text/tabwriter"

	m "github.com/blackflagsoftware/forge-go/internal/model"
)

func PrintSqlColumns(cols []m.Column) {
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
