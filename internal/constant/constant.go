package constant

const (
	POSTGRESQL      = "Postgres"
	MYSQL           = "MySql"
	SQLITE3         = "Sqlite"
	POSTGRESQLLOWER = "postgres"
	MYSQLLOWER      = "mysql"
	SQLITE3LOWER    = "sqlite"

	MODEL_INCLUDE_NULL   = "\n\t\"gopkg.in/guregu/null.v3\""
	MODEL_COLUMN_W_GORM  = "\t\t%s\t%s\t`db:\"%s\" json:\"%s\" gorm:\"column:%s\"`"
	MODEL_COLUMN_WO_GORM = "\t\t%s\t%s\t`db:\"%s\" json:\"%s\"`"

	REST_PRIMARY_INT = `	%sStr := c.Param("%s")
	%s, err := strconv.ParseInt(%sStr, 10, 64)
	if err != nil {
		bindErr := ae.ParseError("Invalid param value, not a number")
		return c.JSON(bindErr.StatusCode, util.NewOutput(bindErr.BodyError(), &bindErr, nil))
	}`  // Lower, Lower, Lower, Lower
	REST_PRIMARY_STR = `	%s := c.Param("%s")` // Lower, Lower
	REST_GET_DELETE  = `	%s := &%s{%s}`       // CamelLower, Camel, RestArgSet
	MANAGER_GET_INT  = `	if %s.%s < 1 {
		return ae.MissingParamError("%s")
	}
	`  // Abbr, Camel, Camel

	MANAGER_GET_STRING = `	if %s.%s == "" {
		return ae.MissingParamError("%s")
	}
	`  // Abbr, Camel, Camel
	MANAGER_POST_AUTOINCREMENT = `if %s.%s < 1 {
		return ae.MissingParamError("%s")
	}
	`  // Abbr, Lower, Camel
	MANAGER_POST_UUID = `if %s.%s == "" {
		return ae.MissingParamError("%s")
	}
	`  // Abbr, Lower, Camel
	MANAGER_POST_NULL = `if !%s.%s.Valid {
		return ae.MissingParamError("%s")
	}
	`  // Abbr, Lower, Camel
	MANAGER_POST_VARCHAR_LEN = `if %s.%s.Valid && len(%s.%s.ValueOrZero()) > %d {
		return ae.StringLengthError("%s", %d)
	}
	`  // Abbr, ColumnCamel, Abbr, ColumnCamel, ColumnLength, ColumnCamel, ColumnLength
	MANAGER_GET_TIME = `	if %s.%s.IsZero() {
		return ae.MissingParamError("%s")
	}
	`  // Abbr, Lower, Camel
	MANAGER_GET_JSON = `	if !%s.%s.ValidJson() {
		return ae.ParseError("%s is invalid JSON")
	}
	`  // Abbr, Lower, Camel
	MANAGER_PATCH_SEARCH_STRING = `	%s, ok%s := jParsed.Search("%s").Data().(string)
	if !ok%s {
		return ae.MissingParamError("%s")
	}
	`  // Lower, Camel, Camel, Camel, Camel
	MANAGER_PATCH_SEARCH_INT = `	%sFlt, ok%sFlt := jParsed.Search("%s").Data().(float64)
	if !ok%sFlt {
		return ae.MissingParamError("%s")
	}
	%s := int(%sFlt)
	`  // Lower, Camel, Camel, Camel, Camel, Lower, Lower
	MANAGER_PATCH_STRUCT_STMT = `	%s := &%s{%s}
	`  // Abbr, Camel, KeySearchList
	MANAGER_PATCH_GET_STMT = `	errGet := m.Get(%s)
	if errGet != nil {
		return errGet
	}
	`  // Abbr
	MANAGER_PATCH_DEFAULT_ASSIGN = `// %s
	if %sIn.%s.Valid {%s
		%s.%s = %sIn.%s
	}
	`  // ColCamel, Abbr, ColCamel, StringLenCheck, Abbr, ColCamel, Abbr. ColCamel
	MANAGER_PATCH_JSON_NULL_ASSIGN = `// %s
	if %sIn.%s != nil {
		if !util.ValidJson(*%sIn.%s) {
			return ae.ParseError("Invalid JSON syntax for %s")
		}
		%s.%s = %sIn.%s
	}
	`  // ColCamel, Abbr, ColCamel, Abbr, ColCamel, ColLowerCamel, Abbr, ColCamel, Abbr, ColCamel
	MANAGER_PATCH_TIME_NULL_ASSIGN = `// %s
	if %sIn.%s.Valid {
		_, errParse := time.Parse(time.RFC3339, %s.%s.Time.String())
		if errParse != nil {
			return ae.ParseError("%s: unable to parse time")
		}
		%s.%s = %sIn.%s
	}
	`  // ColCamel, Abbr, ColCamel, Abbr, ColCamel, ColCamel, Abbr, ColCamel, Abbr, ColCamel
	MANAGER_PATCH_VARCHAR_LEN = `
			if %sIn.%s.Valid && len(%sIn.%s.ValueOrZero()) > %d {
				return ae.StringLengthError("%s", %d)
		}`  // Abbr, ColumnCamel, Abbr, ColumnCamel, ColumnLength, ColumnCamel, ColumnLength
	SQL_POST_QUERY       = `_, errDB := d.DB.NamedExec(sqlPost, %s)`      // Abbr
	SQL_POST_QUERY_MYSQL = `result, errDB := d.DB.NamedExec(sqlPost, %s)` // Abbr
	SQL_LAST_ID_MYSQL    = `lastId, err := result.LastInsertId()
	if err != nil {
		return ae.DBError("%s Create: unable to get lastid.", err)
	}
	%s.%s = int(lastId)
	`  // Camel, Abbr, ColCamel
	SQL_POST_QUERY_POSTGRES = `rows, errDB := d.DB.NamedQuery(sqlPost, %s)` // Abbr
	SQL_LAST_ID_POSTGRES    = `var lastId int64
	if rows.Next() {
		rows.Scan(&lastId)
	}
	%s.%s = int(lastId)
	`  // Abbr, ColCamel

	// TESTS
	REST_TEST_INT_FAILURE = `
	func Test%sRestGetFailureInvalidInt(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/{{.GetDeleteUrl}}")
		c.SetParamNames("%s")
		c.SetParamValues("a")
	
		man := &MockManager%s{}
		h := NewRest%s(man)
	
		h.Get(c)
	
		be := ae.BodyError{}
		json.Unmarshal(rec.Body.Bytes(), &be)
	
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "Invalid param value, not a number", be.Detail)
	}`  // camel, get_delete_url, col_lower, camel, camel
	REST_TEST_INT_ZERO = `
	func Test%sRestGetFailureZeroInt(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/{{.GetDeleteUrl}}")
		c.SetParamNames("%s")
		c.SetParamValues("0")
	
		man := &MockManager%s{}
		h := NewRest%s(man)
	
		h.Get(c)
	
		be := ae.BodyError{}
		json.Unmarshal(rec.Body.Bytes(), &be)
	
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "zero value", be.Detail)
	}`  // camel, get_delete_url, col_lower, camel, camel

	MAIN_COMMON_PATH = `"{{.ProjectFile.ProjectPathEncoded}}\/internal\/{{.ProjectFile.SubPackage}}\/{{.Name.AllLower}}"
	\/\/ --- replace main header text - do not remove ---
`

	SERVER_ROUTE = `{{.Name.AllLower}}.InitializeRest(routeGroup)
	\/\/ --- replace server text - do not remove ---
`
	COMMON_IMPORT = `
import (
	"github.com\/labstack\/echo\/v4"
	
	\/\/ --- replace header text - do not remove ---
)`
	COMMON_HEADER = `{{.Name.Abbr}} "{{.ProjectFile.ProjectPathEncoded}}\/{{.ProjectFile.SubDirEncoded}}\/{{.Name.AllLower}}"
	\/\/ --- replace header text - do not remove ---`

	COMMON_SECTION = `\/\/ {{.Camel}}
func Setup{{.Camel}}(eg *echo.Group) {
	sl := {{.Abbr}}.InitStorage()
	ml := {{.Abbr}}.NewManager{{.Camel}}(sl)
	hl := {{.Abbr}}.NewRest{{.Camel}}(ml)
	hl.Load{{.Camel}}Routes(eg)
}
	
\/\/ --- replace section text - do not remove ---`

	GRPC_IMPORT_ONCE = `pb "{{.ProjectFile.ProjectPathEncoded}}\/pkg\/proto"`

	GRPC_IMPORT = `{{.Name.Abbr}} "{{.ProjectFile.ProjectPathEncoded}}\/{{.ProjectFile.SubDirEncoded}}\/{{.Name.AllLower}}"
	\/\/ --- replace grpc import - do not remove ---`

	GRPC_TEXT = `\/\/ {{.Name.Camel}}
	s{{.Name.Abbr}} := {{.Name.Abbr}}.InitStorage()
	m{{.Name.Abbr}} := {{.Name.Abbr}}.NewManager{{.Name.Camel}}(s{{.Name.Abbr}})
	h{{.Name.Abbr}} := {{.Name.Abbr}}.New{{.Name.Camel}}Grpc(m{{.Name.Abbr}})
	pb.Register{{.Name.Camel}}ServiceServer(s, h{{.Name.Abbr}})
	\/\/ --- replace grpc text - do not remove ---`

	MIGRATION_VERIFY_HEADER_MYSQL = `_ "github.com/go-sql-driver/mysql"`

	MIGRATION_VERIFY_MYSQL = `connectionStr := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, "mysql")
	if dbPass == "" {
		connectionStr = fmt.Sprintf("%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbHost, "mysql")
	}
	db, errOpen := sqlx.Open("mysql", connectionStr)
	if errOpen != nil {
		return fmt.Errorf("Unable to open DB for init: %s", errOpen)
	}
	defer db.Close()

	sqlDatabase := fmt.Sprintf("SELECT EXISTS (SELECT schema_name FROM information_schema.schemata WHERE schema_name = '%s')", expectedDB)
	var exists bool
	errGet := db.Get(&exists, sqlDatabase)
	if errGet != nil {
		return fmt.Errorf("Error Get schema: %s", errGet)
	}
	if !exists {
		sqlCreateDB := fmt.Sprintf("CREATE DATABASE %s", expectedDB)
		_, err := db.Exec(sqlCreateDB)
		if err != nil {
			return fmt.Errorf("Error in creating DB: %s with error: %s", expectedDB, err)
		}
	}
	sqlCreateUser := fmt.Sprintf("CREATE USER IF NOT EXISTS '%s' IDENTIFIED BY '%s'", dbUser, dbPass)
	_, errCreateUser := db.Exec(sqlCreateUser)
	if errCreateUser != nil {
		return fmt.Errorf("Error in creating user: %s", errCreateUser)
	}
	sqlGrantUser := fmt.Sprintf("GRANT ALL ON %s.* TO '%s'@'%%'", expectedDB, dbUser)
	_, errGrantUser := db.Exec(sqlGrantUser)
	if errGrantUser != nil {
		return fmt.Errorf("Error in grant user: %s", errGrantUser)
	}
	return nil`

	MIGRATION_CONNECTION_MYSQL = `connectionStr := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, "mysql")
	if dbPass == "" {
		connectionStr = fmt.Sprintf("%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbHost, "mysql")
	}
	db, errOpen := sqlx.Open("mysql", connectionStr)
	if errOpen != nil {
		fmt.Printf("Unable to open DB for migrations: %s\n", errOpen)
		os.Exit(1)
	}`

	MIGRATION_VERIFY_HEADER_POSTGRES = `_ "github.com/lib/pq"`

	MIGRATION_VERIFY_POSTGRES = `connectionStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable", dbUser, dbPass, "postgres", dbHost)
	if dbPass == "" {
		connectionStr = fmt.Sprintf("user=%s dbname=%s host=%s sslmode=disable", dbUser, "postgres", dbHost)
	}
	db, errOpen := sqlx.Open("postgres", connectionStr)
	if errOpen != nil {
		return fmt.Errorf("Unable to open DB for init: %s", errOpen)
	}
	defer db.Close()

	sqlDatabase := fmt.Sprintf("SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = '%s')", expectedDB)
	var exists bool
	errGet := db.Get(&exists, sqlDatabase)
	if errGet != nil {
		return fmt.Errorf("Error Get schema: %s", errGet)
	}
	if !exists {
		sqlCreateDB := fmt.Sprintf("CREATE DATABASE %s", expectedDB)
		_, err := db.Exec(sqlCreateDB)
		if err != nil {
			return fmt.Errorf("Error in creating DB: %s with error: %s", expectedDB, err)
		}
	}
	sqlUserExists := fmt.Sprintf("SELECT EXISTS(SELECT rolname FROM pg_roles WHERE rolname = '%s')", dbUser)
	errUser := db.Get(&exists, sqlUserExists)
	if errUser != nil {
		return fmt.Errorf("Error get user: %s", errUser)
	}
	if !exists {
		sqlCreateUser := fmt.Sprintf("CREATE USER IF NOT EXISTS %s WITH ENCRYPTED PASSWORD '%s", dbUser, dbPass)
		_, errCreateUser := db.Exec(sqlCreateUser)
		if errCreateUser != nil {
			return fmt.Errorf("Error in creating user: %s", errCreateUser)
		}
		sqlGrantUser := fmt.Sprintf("GRANT ALL PRIVILEGES ON DATABASE %s TO %s", expectedDB, dbUser)
		_, errGrantUser := db.Exec(sqlGrantUser)
		if errGrantUser != nil {
			return fmt.Errorf("Error in grant user: %s", errGrantUser)
		}
	}
	return nil`

	MIGRATION_CONNECTION_POSTGRES = `connectionStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable", dbUser, dbPass, dbDB, dbHost)
	if dbPass == "" {
		connectionStr = fmt.Sprintf("user=%s dbname=%s host=%s sslmode=disable", dbUser, dbDB, dbHost)
	}
	db, errOpen := sqlx.Open("postgres", connectionStr)
	if errOpen != nil {
		fmt.Printf("Unable to open DB for migrations: %s\n", errOpen)
		os.Exit(1)
	}`

	MIGRATION_VERIFY_HEADER_SQLITE = `_ "github.com/mattn/go-sqlite3"`

	MIGRATION_VERIFY_SQLITE = `return nil`

	MIGRATION_CONNECTION_SQLITE = `connectionStr := fmt.Sprintf("%s?cache=shared&mode=wrc", dbHost)
	db, err := sqlx.Open("sqlite3", connectionStr)
	if err != nil {
		fmt.Println("Could not connect with connection string:", connectionStr)
		os.Exit(1)
	}
	db.SetMaxOpenConns(1)
	`

	MIGRATION_CALL = `if config.UseMigration {
		err := os.MkdirAll(config.MigrationPath, 0744)
		if err != nil {
			fmt.Printf("Unable to make scripts\/migrations directory structure: %s\\n", err)
		}

		errVerify := mig.VerifyDBInit(config.DBDB, config.DBHost, config.DBUser, config.DBPass)
		if errVerify != nil {
			panic(errVerify)
		}
		mig.RunMigration(config.MigrationPath, config.DBHost, config.DBUser, config.DBPass, config.DBDB)
	}
`
	MIGRATION_GRPC_HEADER_ONCE = `
	"os"`
)
