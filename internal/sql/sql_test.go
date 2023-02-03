package sql

import (
	"reflect"
	"testing"

	c "github.com/blackflagsoftware/forge-go/internal/column"
	n "github.com/blackflagsoftware/forge-go/internal/name"
)

func Test_tableNameParse(t *testing.T) {
	type args struct {
		lines string
	}
	tests := []struct {
		name           string
		args           args
		wantTableName  string
		wantColumnPart string
		wantErr        bool
	}{
		{
			"tableNameParse: success",
			args{lines: "create table my_table_name (test columns);"},
			"my_table_name",
			"test columns",
			false,
		},
		{
			"tableNameParse: failed",
			args{lines: "update table my_table_name (test columns);"},
			"",
			"",
			true,
		},
		{
			"tableNameParse extra: success",
			args{lines: "create table if not exists my_table_name (test columns);"},
			"my_table_name",
			"test columns",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTableName, gotColumns, err := tableNameParse(tt.args.lines)
			if (err != nil) != tt.wantErr {
				t.Errorf("tableParse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotTableName != tt.wantTableName {
				t.Errorf("tableParse() gotTableName = %v, want %v", gotTableName, tt.wantTableName)
			}
			if gotColumns != tt.wantColumnPart {
				t.Errorf("tableParse() gotColumns = %v, want %v", gotColumns, tt.wantColumnPart)
			}
		})
	}
}

func Test_columnsParse(t *testing.T) {
	type args struct {
		columnPart string
	}
	tests := []struct {
		name        string
		args        args
		wantColumns []c.Column
		wantErr     bool
	}{
		{
			"columnParse - #1",
			args{
				columnPart: "id int auto_increment, user_id int not null, first_name varchar(50) not null, last_name varchar(50) not null, phone varchar(20), age int, active boolean not null default true, primary key(id)",
			},
			[]c.Column{
				{ColumnName: n.Name{RawName: "id"}, DBType: "autoincrement", GoType: "int", GoTypeNonSql: "int", Null: false, DefaultValue: "", Length: 0, PrimaryKey: true},
				{ColumnName: n.Name{RawName: "user_id"}, DBType: "int", GoType: "null.Int", GoTypeNonSql: "int", Null: false, DefaultValue: "", Length: 0, PrimaryKey: false},
				{ColumnName: n.Name{RawName: "first_name"}, DBType: "varchar", GoType: "null.String", GoTypeNonSql: "string", Null: false, DefaultValue: "", Length: 50, PrimaryKey: false},
				{ColumnName: n.Name{RawName: "last_name"}, DBType: "varchar", GoType: "null.String", GoTypeNonSql: "string", Null: false, DefaultValue: "", Length: 50, PrimaryKey: false},
				{ColumnName: n.Name{RawName: "phone"}, DBType: "varchar", GoType: "null.String", GoTypeNonSql: "string", Null: true, DefaultValue: "", Length: 20, PrimaryKey: false},
				{ColumnName: n.Name{RawName: "age"}, DBType: "int", GoType: "null.Int", GoTypeNonSql: "int", Null: true, DefaultValue: "", Length: 0, PrimaryKey: false},
				{ColumnName: n.Name{RawName: "active"}, DBType: "boolean", GoType: "null.Bool", GoTypeNonSql: "bool", Null: false, DefaultValue: "true", Length: 0, PrimaryKey: false},
			},
			false,
		},
		{
			"columnParse - #2",
			args{
				columnPart: "id int NOT NULL AUTO_INCREMENT, uid varchar(255) NOT NULL, tenant_uid varchar(255) NOT NULL, name varchar(255) NOT NULL, plan_currency varchar(20) DEFAULT '$', utility_name varchar(255) DEFAULT NULL, utility_code varchar(255) DEFAULT NULL, PRIMARY KEY (id)",
			},
			[]c.Column{
				{ColumnName: n.Name{RawName: "id"}, DBType: "autoincrement", GoType: "int", GoTypeNonSql: "int", Null: false, DefaultValue: "", Length: 0, PrimaryKey: true},
				{ColumnName: n.Name{RawName: "uid"}, DBType: "varchar", GoType: "null.String", GoTypeNonSql: "string", Null: false, DefaultValue: "", Length: 255, PrimaryKey: false},
				{ColumnName: n.Name{RawName: "tenant_uid"}, DBType: "varchar", GoType: "null.String", GoTypeNonSql: "string", Null: false, DefaultValue: "", Length: 255, PrimaryKey: false},
				{ColumnName: n.Name{RawName: "name"}, DBType: "varchar", GoType: "null.String", GoTypeNonSql: "string", Null: false, DefaultValue: "", Length: 255, PrimaryKey: false},
				{ColumnName: n.Name{RawName: "plan_currency"}, DBType: "varchar", GoType: "null.String", GoTypeNonSql: "string", Null: true, DefaultValue: "'$'", Length: 20, PrimaryKey: false},
				{ColumnName: n.Name{RawName: "utility_name"}, DBType: "varchar", GoType: "null.String", GoTypeNonSql: "string", Null: true, DefaultValue: "", Length: 255, PrimaryKey: false},
				{ColumnName: n.Name{RawName: "utility_code"}, DBType: "varchar", GoType: "null.String", GoTypeNonSql: "string", Null: true, DefaultValue: "", Length: 255, PrimaryKey: false},
			},
			false,
		},
		{
			"columnParse - #3",
			args{
				columnPart: "id int NOT NULL AUTO_INCREMENT, tenant_id int DEFAULT NULL, name varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_unicode_ci DEFAULT NULL, created_at datetime DEFAULT NULL",
			},
			[]c.Column{
				{ColumnName: n.Name{RawName: "id"}, DBType: "autoincrement", GoType: "int", GoTypeNonSql: "int", Null: false, DefaultValue: "", Length: 0, PrimaryKey: false},
				{ColumnName: n.Name{RawName: "tenant_id"}, DBType: "int", GoType: "null.Int", GoTypeNonSql: "int", Null: true, DefaultValue: "", Length: 0, PrimaryKey: false},
				{ColumnName: n.Name{RawName: "name"}, DBType: "varchar", GoType: "null.String", GoTypeNonSql: "string", Null: true, DefaultValue: "", Length: 255, PrimaryKey: false},
				{ColumnName: n.Name{RawName: "created_at"}, DBType: "datetime", GoType: "null.Time", GoTypeNonSql: "time.Time", Null: true, DefaultValue: "", Length: 0, PrimaryKey: false},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotColumns, err := columnsParse(tt.args.columnPart)
			if (err != nil) != tt.wantErr {
				t.Errorf("columnsParse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotColumns, tt.wantColumns) {
				t.Errorf("columnsParse() = %v, want %v", gotColumns, tt.wantColumns)
			}
		})
	}
}

func Test_theRestParse(t *testing.T) {
	type args struct {
		theRest string
	}
	tests := []struct {
		name             string
		args             args
		col              *c.Column
		wantNull         bool
		wantDefaultValue string
		wantDBType       string
		wantPrimaryKey   bool
	}{
		{
			"theRestParse - null",
			args{
				"null",
			},
			&c.Column{},
			true,
			"",
			"",
			false,
		},
		{
			"theRestParse - not null",
			args{
				"not null",
			},
			&c.Column{},
			false,
			"",
			"",
			false,
		},
		{
			"theRestParse - autoincrement",
			args{
				"not null auto_increment",
			},
			&c.Column{},
			false,
			"",
			"autoincrement",
			false,
		},
		{
			"theRestParse - primary",
			args{
				"primary key not null",
			},
			&c.Column{},
			false,
			"",
			"",
			true,
		},
		{
			"theRestParse - default int",
			args{
				"default '0'",
			},
			&c.Column{},
			true,
			"'0'",
			"",
			false,
		},
		{
			"theRestParse - default false",
			args{
				"default false",
			},
			&c.Column{},
			true,
			"false",
			"",
			false,
		},
		{
			"theRestParse - addl stuff",
			args{
				"character set utf8mb3 collate utf8mb3_unicode_ci default null",
			},
			&c.Column{},
			true,
			"",
			"",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			theRestParse(tt.col, tt.args.theRest)
		})
		if tt.col.Null != tt.wantNull {
			t.Errorf("theRestParse() = column.Null: %v, want %v", tt.col.Null, tt.wantNull)
		}
		if tt.col.DefaultValue != tt.wantDefaultValue {
			t.Errorf("theRestParse() = column.DefaultValue: %v, want %v", tt.col.DefaultValue, tt.wantDefaultValue)
		}
		if tt.col.DBType != tt.wantDBType {
			t.Errorf("theRestParse() = column.DBType: %v, want %v", tt.col.DBType, tt.wantDBType)
		}
		if tt.col.PrimaryKey != tt.wantPrimaryKey {
			t.Errorf("theRestParse() = column.PrimaryKey: %v, want %v", tt.col.PrimaryKey, tt.wantPrimaryKey)
		}
	}
}
