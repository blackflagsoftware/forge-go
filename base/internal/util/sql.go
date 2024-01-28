package util

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
)

type (
	SearchBuilder struct {
		Params []string
		Values []interface{}
	}
)

func BuildSearchString(search []ParamSearch) (string, []interface{}) {
	sb := SearchBuilder{}
	for _, s := range search {
		switch s.Compare {
		case "LIKE":
			sb.AppendLike(s.Column, s.Value.(string))
		case "NULL":
			sb.AppendNull(s.Column, true)
		case "NOT NULL":
			sb.AppendNull(s.Column, false)
		case "IN":
			sb.AppendIn(s.Column, s.Value)
		default:
			if s.Compare != "" {
				sb.AppendCompare(s.Column, s.Compare, s.Value)
			}
		}
	}
	return sb.String(), sb.Values
}

func (s *SearchBuilder) AppendCompare(param, compare string, value interface{}) {
	s.Params = append(s.Params, fmt.Sprintf("%s %s ?", param, compare))
	s.Values = append(s.Values, value)
}

func (s *SearchBuilder) AppendLike(param, value string) {
	s.Params = append(s.Params, fmt.Sprintf("%s LIKE '%%%s%%'", param, value))
}

func (s *SearchBuilder) AppendNull(param string, wantNull bool) {
	nullStmt := "IS NOT NULL"
	if wantNull {
		nullStmt = "IS NULL"
	}
	s.Params = append(s.Params, fmt.Sprintf("%s %s", param, nullStmt))
}

// this will produce a string that represents an "IN" clause as long as the incoming arg 'value' is of type slice or array
// the output will always be sql string array, most DB engines will deal with the single quotes even if the underlying columns is not a text type
// []string{1, 2, 3} => IN ('1', '2', '3'), this is what is expected
func (s *SearchBuilder) AppendIn(param string, value interface{}) {
	slice := reflect.ValueOf(value)
	if slice.Kind() == reflect.Slice || slice.Kind() == reflect.Array {
		b := make([]interface{}, slice.Len())
		for i := 0; i < slice.Len(); i++ {
			b[i] = slice.Index(i).Interface()
		}
		inListArray := []string{}
		for _, i := range b {
			inListArray = append(inListArray, fmt.Sprintf("'%v'", i))
		}
		inList := strings.Join(inListArray, ", ")
		s.Params = append(s.Params, fmt.Sprintf("%s IN (%s)", param, inList))
		return
	}
	fmt.Println("Not a slice")
}

func (s *SearchBuilder) String() string {
	if len(s.Params) > 0 {
		return fmt.Sprintf("WHERE %s", strings.Join(s.Params, "\n\t\tAND "))
	}
	return ""
}

func TxnFinish(tx *sqlx.Tx, err *error) {
	if p := recover(); p != nil {
		tx.Rollback()
		panic(p)
	} else if *err != nil {
		tx.Rollback()
	} else {
		if errCommit := tx.Commit(); errCommit != nil {
			err = &errCommit
		}
	}
}
