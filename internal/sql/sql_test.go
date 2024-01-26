package sql

import (
	"reflect"
	"testing"

	m "github.com/blackflagsoftware/forge-go/internal/model"
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
		wantColumns []m.Column
		wantErr     bool
	}{
		{
			"columnParse - #1",
			args{
				columnPart: "id int auto_increment, user_id int not null, first_name varchar(50) not null, last_name varchar(50) not null, phone varchar(20), age int, active boolean not null default true, primary key(id)",
			},
			[]m.Column{
				{ColumnName: m.Name{RawName: "id", Lower: "id", Camel: "Id", LowerCamel: "id", Upper: "ID", Abbr: "id", EnvVar: "ID", AllLower: "id"}, DBType: "autoincrement", GoType: "int", GoTypeNonSql: "int", Null: false, DefaultValue: "", Length: 0, PrimaryKey: true},
				{ColumnName: m.Name{RawName: "user_id", Lower: "user_id", Camel: "UserId", LowerCamel: "userId", Upper: "USERID", Abbr: "use", EnvVar: "USER_ID", AllLower: "userid"}, DBType: "int", GoType: "null.Int", GoTypeNonSql: "int", Null: false, DefaultValue: "", Length: 0, PrimaryKey: false},
				{ColumnName: m.Name{RawName: "first_name", Lower: "first_name", Camel: "FirstName", LowerCamel: "firstName", Upper: "FIRSTNAME", Abbr: "fir", EnvVar: "FIRST_NAME", AllLower: "firstname"}, DBType: "varchar", GoType: "null.String", GoTypeNonSql: "string", Null: false, DefaultValue: "", Length: 50, PrimaryKey: false},
				{ColumnName: m.Name{RawName: "last_name", Lower: "last_name", Camel: "LastName", LowerCamel: "lastName", Upper: "LASTNAME", Abbr: "las", EnvVar: "LAST_NAME", AllLower: "lastname"}, DBType: "varchar", GoType: "null.String", GoTypeNonSql: "string", Null: false, DefaultValue: "", Length: 50, PrimaryKey: false},
				{ColumnName: m.Name{RawName: "phone", Lower: "phone", Camel: "Phone", LowerCamel: "phone", Upper: "PHONE", Abbr: "pho", EnvVar: "PHONE", AllLower: "phone"}, DBType: "varchar", GoType: "null.String", GoTypeNonSql: "string", Null: true, DefaultValue: "", Length: 20, PrimaryKey: false},
				{ColumnName: m.Name{RawName: "age", Lower: "age", Camel: "Age", LowerCamel: "age", Upper: "AGE", Abbr: "age", EnvVar: "AGE", AllLower: "age"}, DBType: "int", GoType: "null.Int", GoTypeNonSql: "int", Null: true, DefaultValue: "", Length: 0, PrimaryKey: false},
				{ColumnName: m.Name{RawName: "active", Lower: "active", Camel: "Active", LowerCamel: "active", Upper: "ACTIVE", Abbr: "act", EnvVar: "ACTIVE", AllLower: "active"}, DBType: "boolean", GoType: "null.Bool", GoTypeNonSql: "bool", Null: false, DefaultValue: "true", Length: 0, PrimaryKey: false},
			},
			false,
		},
		{
			"columnParse - #2",
			args{
				columnPart: "id int NOT NULL AUTO_INCREMENT, uid varchar(255) NOT NULL, tenant_uid varchar(255) NOT NULL, name varchar(255) NOT NULL, plan_currency varchar(20) DEFAULT '$', utility_name varchar(255) DEFAULT NULL, utility_code varchar(255) DEFAULT NULL, PRIMARY KEY (id)",
			},
			[]m.Column{
				{ColumnName: m.Name{RawName: "id", Lower: "id", Camel: "Id", LowerCamel: "id", Upper: "ID", Abbr: "id", EnvVar: "ID", AllLower: "id"}, DBType: "autoincrement", GoType: "int", GoTypeNonSql: "int", Null: false, DefaultValue: "", Length: 0, PrimaryKey: true},
				{ColumnName: m.Name{RawName: "uid", Lower: "uid", Camel: "Uid", LowerCamel: "uid", Upper: "UID", Abbr: "uid", EnvVar: "UID", AllLower: "uid"}, DBType: "varchar", GoType: "null.String", GoTypeNonSql: "string", Null: false, DefaultValue: "", Length: 255, PrimaryKey: false},
				{ColumnName: m.Name{RawName: "tenant_uid", Lower: "tenant_uid", Camel: "TenantUid", LowerCamel: "tenantUid", Upper: "TENANTUID", Abbr: "ten", EnvVar: "TENANT_UID", AllLower: "tenantuid"}, DBType: "varchar", GoType: "null.String", GoTypeNonSql: "string", Null: false, DefaultValue: "", Length: 255, PrimaryKey: false},
				{ColumnName: m.Name{RawName: "name", Lower: "name", Camel: "Name", LowerCamel: "name", Upper: "NAME", Abbr: "nam", EnvVar: "NAME", AllLower: "name"}, DBType: "varchar", GoType: "null.String", GoTypeNonSql: "string", Null: false, DefaultValue: "", Length: 255, PrimaryKey: false},
				{ColumnName: m.Name{RawName: "plan_currency", Lower: "plan_currency", Camel: "PlanCurrency", LowerCamel: "planCurrency", Upper: "PLANCURRENCY", Abbr: "pla", EnvVar: "PLAN_CURRENCY", AllLower: "plancurrency"}, DBType: "varchar", GoType: "null.String", GoTypeNonSql: "string", Null: true, DefaultValue: "'$'", Length: 20, PrimaryKey: false},
				{ColumnName: m.Name{RawName: "utility_name", Lower: "utility_name", Camel: "UtilityName", LowerCamel: "utilityName", Upper: "UTILITYNAME", Abbr: "uti", EnvVar: "UTILITY_NAME", AllLower: "utilityname"}, DBType: "varchar", GoType: "null.String", GoTypeNonSql: "string", Null: true, DefaultValue: "", Length: 255, PrimaryKey: false},
				{ColumnName: m.Name{RawName: "utility_code", Lower: "utility_code", Camel: "UtilityCode", LowerCamel: "utilityCode", Upper: "UTILITYCODE", Abbr: "uti", EnvVar: "UTILITY_CODE", AllLower: "utilitycode"}, DBType: "varchar", GoType: "null.String", GoTypeNonSql: "string", Null: true, DefaultValue: "", Length: 255, PrimaryKey: false},
			},
			false,
		},
		{
			"columnParse - #3",
			args{
				columnPart: "id int NOT NULL AUTO_INCREMENT, tenant_id int DEFAULT NULL, name varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_unicode_ci DEFAULT NULL, created_at datetime DEFAULT NULL",
			},
			[]m.Column{
				{ColumnName: m.Name{RawName: "id", Lower: "id", Camel: "Id", LowerCamel: "id", Upper: "ID", Abbr: "id", EnvVar: "ID", AllLower: "id"}, DBType: "autoincrement", GoType: "int", GoTypeNonSql: "int", Null: false, DefaultValue: "", Length: 0, PrimaryKey: false},
				{ColumnName: m.Name{RawName: "tenant_id", Lower: "tenant_id", Camel: "TenantId", LowerCamel: "tenantId", Upper: "TENANTID", Abbr: "ten", EnvVar: "TENANT_ID", AllLower: "tenantid"}, DBType: "int", GoType: "null.Int", GoTypeNonSql: "int", Null: true, DefaultValue: "", Length: 0, PrimaryKey: false},
				{ColumnName: m.Name{RawName: "name", Lower: "name", Camel: "Name", LowerCamel: "name", Upper: "NAME", Abbr: "nam", EnvVar: "NAME", AllLower: "name"}, DBType: "varchar", GoType: "null.String", GoTypeNonSql: "string", Null: true, DefaultValue: "", Length: 255, PrimaryKey: false},
				{ColumnName: m.Name{RawName: "created_at", Lower: "created_at", Camel: "CreatedAt", LowerCamel: "createdAt", Upper: "CREATEDAT", Abbr: "cre", EnvVar: "CREATED_AT", AllLower: "createdat"}, DBType: "datetime", GoType: "null.Time", GoTypeNonSql: "time.Time", Null: true, DefaultValue: "", Length: 0, PrimaryKey: false},
			},
			false,
		},
		{
			"columnParse - #4",
			args{
				columnPart: "id int NOT NULL AUTO_INCREMENT, tenant_id int DEFAULT NULL, name varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_unicode_ci DEFAULT NULL, score numeric(12,2) default 10.2, created_at datetime DEFAULT NULL",
			},
			[]m.Column{
				{ColumnName: m.Name{RawName: "id", Lower: "id", Camel: "Id", LowerCamel: "id", Upper: "ID", Abbr: "id", EnvVar: "ID", AllLower: "id"}, DBType: "autoincrement", GoType: "int", GoTypeNonSql: "int", Null: false, DefaultValue: "", Length: 0, PrimaryKey: false},
				{ColumnName: m.Name{RawName: "tenant_id", Lower: "tenant_id", Camel: "TenantId", LowerCamel: "tenantId", Upper: "TENANTID", Abbr: "ten", EnvVar: "TENANT_ID", AllLower: "tenantid"}, DBType: "int", GoType: "null.Int", GoTypeNonSql: "int", Null: true, DefaultValue: "", Length: 0, PrimaryKey: false},
				{ColumnName: m.Name{RawName: "name", Lower: "name", Camel: "Name", LowerCamel: "name", Upper: "NAME", Abbr: "nam", EnvVar: "NAME", AllLower: "name"}, DBType: "varchar", GoType: "null.String", GoTypeNonSql: "string", Null: true, DefaultValue: "", Length: 255, PrimaryKey: false},
				{ColumnName: m.Name{RawName: "score", Lower: "score", Camel: "Score", LowerCamel: "score", Upper: "SCORE", Abbr: "sco", EnvVar: "SCORE", AllLower: "score"}, DBType: "numeric(12,2)", GoType: "null.Float", GoTypeNonSql: "float64", Null: true, DefaultValue: "10.2", Length: 0, PrimaryKey: false},
				{ColumnName: m.Name{RawName: "created_at", Lower: "created_at", Camel: "CreatedAt", LowerCamel: "createdAt", Upper: "CREATEDAT", Abbr: "cre", EnvVar: "CREATED_AT", AllLower: "createdat"}, DBType: "datetime", GoType: "null.Time", GoTypeNonSql: "time.Time", Null: true, DefaultValue: "", Length: 0, PrimaryKey: false},
			},
			false,
		},
		{
			"columnParse - #5",
			args{
				columnPart: "client_id int not null, user_id int not null, primary key(client_id, user_id)",
			},
			[]m.Column{
				{ColumnName: m.Name{RawName: "client_id", Lower: "client_id", Camel: "ClientId", LowerCamel: "clientId", Upper: "CLIENTID", Abbr: "cli", EnvVar: "CLIENT_ID", AllLower: "clientid"}, DBType: "int", GoType: "null.Int", GoTypeNonSql: "int", Null: false, DefaultValue: "", Length: 0, PrimaryKey: true},
				{ColumnName: m.Name{RawName: "user_id", Lower: "user_id", Camel: "UserId", LowerCamel: "userId", Upper: "USERID", Abbr: "use", EnvVar: "USER_ID", AllLower: "userid"}, DBType: "int", GoType: "null.Int", GoTypeNonSql: "int", Null: false, DefaultValue: "", Length: 0, PrimaryKey: true},
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
		col              *m.Column
		wantNull         bool
		wantDefaultValue string
		wantDBType       string
		wantPrimaryKey   bool
	}{
		{
			"theRestParse - empty",
			args{
				"",
			},
			&m.Column{},
			true,
			"",
			"",
			false,
		},
		{
			"theRestParse - null",
			args{
				"null",
			},
			&m.Column{},
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
			&m.Column{},
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
			&m.Column{},
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
			&m.Column{},
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
			&m.Column{},
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
			&m.Column{},
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
			&m.Column{},
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

func Test_replaceNumericPart(t *testing.T) {
	type args struct {
		columnPart string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"nothing - success",
			args{"create table user (id serial, last_name varchar(100), age int, active bool, primary key(id))"},
			"create table user (id serial, last_name varchar(100), age int, active bool, primary key(id))",
		},
		{
			"replace one - success",
			args{"create table user (id serial, last_name varchar(100), age int, active bool, amount numeric(16,2), primary key(id))"},
			"create table user (id serial, last_name varchar(100), age int, active bool, amount numeric(16_2), primary key(id))",
		},
		{
			"replace all - success",
			args{"create table user (id serial, last_name varchar(100), age int, active bool, amount numeric(16,2), debit decimal(14,4), primary key(id))"},
			"create table user (id serial, last_name varchar(100), age int, active bool, amount numeric(16_2), debit decimal(14_4), primary key(id))",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := replaceNumericPart(tt.args.columnPart); got != tt.want {
				t.Errorf("replaceNumericPart() = %v, want %v", got, tt.want)
			}
		})
	}
}
