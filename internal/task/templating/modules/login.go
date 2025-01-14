package modules

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"text/template"
	"time"

	m "github.com/blackflagsoftware/forge-go/internal/model"
	temp "github.com/blackflagsoftware/forge-go/internal/task/templating"
	"github.com/blackflagsoftware/forge-go/internal/util"
)

type (
	ScriptFile struct {
		FileName string
		Script   string
	}
)

func AddAuth(p *m.Project) {
	Config(*p)
	Errors(*p)
	ProtoFile(*p)
	TemplateFiles(*p)
	Readme(*p)
	// auth entry
	authEntity := m.Entity{ModuleName: "auth", SkipGrpc: true}
	authName := authEntity.Name.BuildName("auth", p.ProjectFile.KnownAliases)
	p.KnownAliases = append(p.KnownAliases, authName)
	p.Entities = append(p.Entities, authEntity)
	// authauthorize
	authAuthEntity := m.Entity{ModuleName: "authauthorize", SkipGrpc: true}
	authAuthName := authAuthEntity.Name.BuildName("auth_authorize", p.ProjectFile.KnownAliases)
	p.KnownAliases = append(p.KnownAliases, authAuthName)
	p.Entities = append(p.Entities, authAuthEntity)
	// authclient
	authClientEntity := m.Entity{ModuleName: "authclient", SkipGrpc: true}
	authClientName := authClientEntity.Name.BuildName("auth_client", p.ProjectFile.KnownAliases)
	p.KnownAliases = append(p.KnownAliases, authClientName)
	p.Entities = append(p.Entities, authClientEntity)
	// authclientcallback
	authClientCallbackEntity := m.Entity{ModuleName: "authclientcallback", SkipGrpc: true}
	authClientCallbackName := authClientCallbackEntity.Name.BuildName("auth_client_callback", p.ProjectFile.KnownAliases)
	p.KnownAliases = append(p.KnownAliases, authClientCallbackName)
	p.Entities = append(p.Entities, authClientCallbackEntity)
	// authclientsecret
	authClientSecretEntity := m.Entity{ModuleName: "authclientsecret", SkipGrpc: true}
	authClientSecretName := authClientSecretEntity.Name.BuildName("auth_client_secret", p.ProjectFile.KnownAliases)
	p.KnownAliases = append(p.KnownAliases, authClientSecretName)
	p.Entities = append(p.Entities, authClientSecretEntity)
	// authrefresh
	authRefreshEntity := m.Entity{ModuleName: "authrefresh", SkipGrpc: true}
	authRefreshName := authRefreshEntity.Name.BuildName("auth_refresh", p.ProjectFile.KnownAliases)
	p.KnownAliases = append(p.KnownAliases, authRefreshName)
	p.Entities = append(p.Entities, authRefreshEntity)
	// login entry
	loginEntity := m.Entity{ModuleName: "login"}
	loginName := loginEntity.Name.BuildName("login", p.ProjectFile.KnownAliases)
	p.KnownAliases = append(p.KnownAliases, loginName)
	p.Entities = append(p.Entities, loginEntity)
	// add role and login_role, this will add the entities
	buildRoleEntities(p)
	temp.StartTemplating(p)
	MigrationScripts(*p)
	p.Entities = []m.Entity{}
}

func Config(p m.Project) {
	configVarLine := `LoginPwdCost   string
	LoginResetDuration     	 string
	LoginExpiresAtDuration 	 string
	LoginAuthAlg           	 string
	LoginAuthSecret        	 string
	LoginAuthPublic        	 string
	LoginEmailHost         	 string
	LoginEmailPort         	 string
	LoginEmailPwd          	 string
	LoginEmailFrom         	 string
	LoginEmailResetUrl     	 string
	LoginAdminEmail        	 string
	LoginBasicAuthUser     	 string
	LoginBasicAuthPwd      	 string
	AuthAuthorizationExpires string
	AuthRefreshTokenExpires  string
	\/\/ --- replace config var text - do not remove ---
	`
	configInitLine := `LoginPwdCost           = GetEnvOrDefault("{{.ProjectNameEnv}}_LOGIN_PWD_COST", "10")             \/\/ algorithm cost
	LoginResetDuration     = GetEnvOrDefault("{{.ProjectNameEnv}}_LOGIN_RESET_DURATION", "7")        \/\/ in days
	LoginExpiresAtDuration = GetEnvOrDefault("{{.ProjectNameEnv}}_LOGIN_EXPIRES_AT_DURATION", "168") \/\/ in hours (7 days)
	LoginAuthAlg           = GetEnvOrDefault("{{.ProjectNameEnv}}_LOGIN_AUTH_ALG", "HMAC")           \/\/ HMAC, RSA, ECDSA or EdDSA (only use the 512 size)
	LoginAuthSecret        = GetEnvOrDefault("{{.ProjectNameEnv}}_LOGIN_AUTH_SECRET", "")            \/\/ base64 format: used by all 3, HMAC is insecure use only for dev
	LoginAuthPublic        = GetEnvOrDefault("{{.ProjectNameEnv}}_LOGIN_AUTH_PUBLIC", "")            \/\/ base64 format: only used by RSA or ECDSA
	LoginEmailHost         = GetEnvOrDefault("{{.ProjectNameEnv}}_EMAIL_HOST", "")
	LoginEmailPort         = GetEnvOrDefault("{{.ProjectNameEnv}}_EMAIL_PORT", "587")
	LoginEmailPwd          = GetEnvOrDefault("{{.ProjectNameEnv}}_EMAIL_PWD", "")
	LoginEmailFrom         = GetEnvOrDefault("{{.ProjectNameEnv}}_EMAIL_FROM", "")
	LoginEmailResetUrl     = GetEnvOrDefault("{{.ProjectNameEnv}}_EMAIL_RESET_URL", "")
	LoginAdminEmail        = GetEnvOrDefault("{{.ProjectNameEnv}}_ADMIN_EMAIL", "")
	LoginBasicAuthUser     = GetEnvOrDefault("{{.ProjectNameEnv}}_BASIC_AUTH_USER", "")
	LoginBasicAuthPwd      = GetEnvOrDefault("{{.ProjectNameEnv}}_BASIC_AUTH_PASS", "")
	AuthAuthorizationExpires = GetEnvOrDefault("AUTH_SERVER_AUTH_AUTHORIZATION_EXPIRES", "60") \/\/ in seconds
	AuthRefreshTokenExpires = GetEnvOrDefault("AUTH_REFRESH_TOKEN_EXPIRES", "86400")           \/\/ in seconds, set -1 to never expire; 0 - to always refresh; >0 in seconds to expire at
	\/\/ --- replace config init text - do not remove ---
	`
	configFile := fmt.Sprintf("%s/config/config.go", p.ProjectFile.FullPath)
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		fmt.Printf("%s is missing unable to write in hooks\n", configFile)
	} else {
		var configReplace bytes.Buffer
		// config var lines
		tConfig := template.Must(template.New("config").Parse(configVarLine))
		errConfig := tConfig.Execute(&configReplace, p)
		if errConfig != nil {
			fmt.Printf("%s: template error [%s]\n", configFile, errConfig)
		} else {
			cmdConfig := fmt.Sprintf(`perl -pi -e 's/\/\/ --- replace config var text - do not remove ---/%s/g' %s`, configReplace.String(), configFile)
			execConfig := exec.Command("bash", "-c", cmdConfig)
			errConfigCmd := execConfig.Run()
			if errConfigCmd != nil {
				fmt.Printf("%s: error in replace for config text [%s]\n", configFile, errConfigCmd)
			}
		}
		// config init lines
		tConfig = template.Must(template.New("config").Parse(configInitLine))
		errConfig = tConfig.Execute(&configReplace, p)
		if errConfig != nil {
			fmt.Printf("%s: template error [%s]\n", configFile, errConfig)
		} else {
			cmdConfig := fmt.Sprintf(`perl -pi -e 's/\/\/ --- replace config inti text - do not remove ---/%s/g' %s`, configReplace.String(), configFile)
			execConfig := exec.Command("bash", "-c", cmdConfig)
			errConfigCmd := execConfig.Run()
			if errConfigCmd != nil {
				fmt.Printf("%s: error in replace for config text [%s]\n", configFile, errConfigCmd)
			}
		}
		cmdConfigImport := fmt.Sprintf(`perl -pi -e 's/\/\/ --- replace config import text - do not remove ---/%s/g' %s`, "\"strconv\"", configFile)
		execConfig := exec.Command("bash", "-c", cmdConfigImport)
		errConfigImportCmd := execConfig.Run()
		if errConfigImportCmd != nil {
			fmt.Printf("%s: error in replace for config import text [%s]\n", configFile, errConfigImportCmd)
		}

	}

	configFunc := `
func GetPwdCost() int {
	cost, err := strconv.Atoi(LoginPwdCost)
	if err != nil {
		// TODO: unable to print to default log, might want to send error to another feedback loop
		fmt.Printf("GetPwdCost: unable to parse env var: %s", err)
		return 10
	}
	return cost
}

func GetResetDuration() int {
	durationInDays, err := strconv.Atoi(LoginResetDuration)
	if err != nil {
		// TODO: unable to print to default log, might want to send error to another feedback loop
		fmt.Printf("GetResetDuration: unable to parse env var: %s", err)
		return 7
	}
	return durationInDays
}

func GetExpiresAtDuration() int {
	durationInHours, err := strconv.Atoi(LoginExpiresAtDuration)
	if err != nil {
		// TODO: unable to print to default log, might want to send error to another feedback loop
		fmt.Printf("GetExpiredAtDuration: unable to parse env var: %s", err)
		return 7
	}
	return durationInHours
}

func GetEmailPort() int {
	loginEmailPort, err := strconv.Atoi(LoginEmailPort)
	if err != nil {
		// TODO: unable to print to default log, might want to send error to another feedback loop
		fmt.Printf("GetEmailPort: unable to parse env var: %s", err)
		return 7
	}
	return loginEmailPort
}
	
func GetAuthorizationExpires() int {
	authExpires, err := strconv.Atoi(AuthAuthorizationExpires)
	if err != nil {
		// TODO: unable to print to default log, might want to send error to another feedback loop
		fmt.Printf("GetAuthorizationExpires: unable to parse env var: %s", err)
		// TODO: change in forge too
		return 3600 // 1 hour
	}
	// TODO: change in forge too
	return authExpires * 60 * 60
}

func GetRefreshTokenExpires() int {
	authExpires, err := strconv.Atoi(AuthRefreshTokenExpires)
	if err != nil {
		// TODO: unable to print to default log, might want to send error to another feedback loop
		fmt.Printf("GetRefreshTokenExpires: unable to parse env var: %s", err)
		// TODO: change in forge too
		return 86400 // 1 day
	}
	// TODO: change in forge too
	return authExpires
}
`

	file, err := os.OpenFile(configFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening the config file")
		return
	}
	defer file.Close()
	if _, err := file.WriteString(configFunc); err != nil {
		fmt.Println("Error writing to the config file")
		return
	}
}

func Errors(p m.Project) {
	errorsFile := fmt.Sprintf("%s/internal/api_error/errors.go", p.ProjectFile.FullPath)
	errors := `	
func PasswordValiationError(msg string) ApiError {
	return NewApiError(
		http.StatusBadRequest,
		"Password Validation Error",
		fmt.Sprintf("Invalid password, reason: %s", msg),
		false,
		nil,
	)
}

func EmailValidError(msg string) ApiError {
	return NewApiError(
		http.StatusBadRequest,
		"Email Validation Error",
		fmt.Sprintf("Invalid email, reason: %s", msg),
		false,
		nil,
	)
}

func ResetTokenInvalidError() ApiError {
	return NewApiError(
		http.StatusBadRequest,
		"Invalid Reset Token",
		"Token missing/expired, please repeat the reset password process",
		false,
		nil,
	)
}

func LoginActiveError() ApiError {
	return NewApiError(
		http.StatusBadRequest,
		"Login Inactive",
		"The login user inactive, please contact the site administrator",
		false,
		nil,
	)
}

func EmailPasswordComboError() ApiError {
	return NewApiError(
		http.StatusBadRequest,
		"Invalid Email/Password Combination",
		fmt.Sprintf("The email/password entered was not valid, please try again"),
		false,
		nil,
	)
}

func DuplicateEmailError(email string) ApiError {
	return NewApiError(
		http.StatusBadRequest,
		"Duplicate Email",
		fmt.Sprintf("The email: %s already exists in this system", email),
		false,
		nil,
	)
}`
	file, err := os.OpenFile(errorsFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening the %s file: %v\n", errorsFile, err)
		return
	}
	defer file.Close()
	if _, err := file.WriteString(errors); err != nil {
		fmt.Println("Error writing to the api_error file")
		return
	}
}

func ProtoFile(p m.Project) {
	protoFile := fmt.Sprintf("%s/pkg/proto/%s.proto", p.ProjectFile.FullPath, p.ProjectFile.AppName)
	proto := `
message Login {
	string Id = 1;
	string EmailAddr = 2;
	string Pwd = 3;
	bool Active = 4;
	bool SetPwd = 5;
	string CreatedAt = 6;
	string UpdatedAt = 7;
}

message LoginResponse {
	Login Login = 1;
	Result result = 2;
}

message LoginRepeatResponse {
	repeated Login Login = 1;
	Result result = 2;
}

service LoginService {
	rpc GetLogin(LoginIDIn) returns (LoginResponse);
	rpc SearchLogin(Login) returns (LoginRepeatResponse);
	rpc CreateLogin(Login) returns (LoginResponse);
	rpc UpdateLogin(Login) returns (Result);
	rpc DeleteLogin(LoginIDIn) returns (Result);
}

message LoginIDIn {
	string Id = 1;
}

message Role {
	string Id = 1;
	string Name = 2;
	string Description = 3;
}

message RoleResponse {
	Role Role = 1;
	Result result = 2;
}

message RoleRepeatResponse {
	repeated Role Role = 1;
	Result result = 2;
}

service RoleService {
	rpc GetRole(RoleIDIn) returns (RoleResponse);
	rpc SearchRole(Role) returns (RoleRepeatResponse);
	rpc CreateRole(Role) returns (RoleResponse);
	rpc UpdateRole(Role) returns (Result);
	rpc DeleteRole(RoleIDIn) returns (Result);
}

message RoleIDIn {
	string Id = 1;
}

message LoginRole {
	string LoginId = 1;
	string RoleId = 2;
}

message LoginRoleResponse {
	LoginRole LoginRole = 1;
	Result result = 2;
}

message LoginRolePatch {
	string LoginId = 1;
	repeated string RoleIds = 2;
}

message LoginRoleRepeatResponse {
	repeated LoginRole LoginRole = 1;
	Result result = 2;
}

service LoginRoleService {
	rpc GetLoginRole(LoginRoleIDsIn) returns (LoginRoleResponse);
	rpc SearchLoginRole(LoginRole) returns (LoginRoleRepeatResponse);
	rpc CreateLoginRole(LoginRole) returns (LoginRoleResponse);
	rpc UpdateLoginRole(LoginRolePatch) returns (Result);
	rpc DeleteLoginRole(LoginRoleIDsIn) returns (Result);
}

message LoginRoleIDsIn {
	string LoginId = 1;
	string RoleId = 2;
}`
	file, err := os.OpenFile(protoFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening the %s file: %v\n", protoFile, err)
		return
	}
	defer file.Close()
	if _, err := file.WriteString(proto); err != nil {
		fmt.Println("Error writing to the proto file")
		return
	}
}

func TemplateFiles(p m.Project) {
	type Files struct {
		Src   string
		Dest  string
		Mkdir bool
	}
	files := []Files{
		{fmt.Sprintf("%s/modules/login/middleware/auth.tmpl", os.Getenv("FORGE_PATH")), fmt.Sprintf("%s/internal/middleware/auth.go", p.ProjectFile.FullPath), false},
		{fmt.Sprintf("%s/modules/login/tools_admin/main.tmpl", os.Getenv("FORGE_PATH")), fmt.Sprintf("%s/tools/admin/main.go", p.ProjectFile.FullPath), true},
		{fmt.Sprintf("%s/modules/login/util/password_test.tmpl", os.Getenv("FORGE_PATH")), fmt.Sprintf("%s/internal/util/password_test.go", p.ProjectFile.FullPath), false},
		{fmt.Sprintf("%s/modules/login/util/password.tmpl", os.Getenv("FORGE_PATH")), fmt.Sprintf("%s/internal/util/password.go", p.ProjectFile.FullPath), false},
		{fmt.Sprintf("%s/modules/login/util/email/email.tmpl", os.Getenv("FORGE_PATH")), fmt.Sprintf("%s/internal/util/email/email.go", p.ProjectFile.FullPath), true},
	}
	for _, f := range files {
		t, errParse := template.ParseFiles(f.Src)
		if errParse != nil {
			fmt.Printf("Template could not parse auth file: %s; %s\n", f.Src, errParse)
			return
		}
		if f.Mkdir {
			dirOnly := path.Dir(f.Dest)
			if err := os.MkdirAll(dirOnly, os.ModePerm); err != nil {
				fmt.Printf("Unable to make dest directory: %s; %v\n", f.Dest, err)
				return
			}
		}
		file, err := os.Create(f.Dest)
		if err != nil {
			fmt.Printf("File: %s was not able to be created; %v\n", f.Dest, err)
			return
		}
		err = t.Execute(file, p)
		if err != nil {
			fmt.Println("Error execution of template:", err)
		}
	}
}

func MigrationScripts(p m.Project) {
	// set postgres by default
	uuid := "UUID"
	ts := "TIMESTAMP"
	b := "BOOL"
	if p.ProjectFile.SqlStorage == "m" {
		// mysql
		uuid = "CHAR(36)"
		ts = "DATETIME"
		b = "BOOLEAN"
	}
	if p.ProjectFile.SqlStorage == "s" {
		// sqlite
		uuid = "TEXT"
		ts = "INTEGER"
		b = "INTEGER"
	}
	loginScript := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS login (
	id %s NOT NULL,
	email_addr VARCHAR(100) NOT NULL,
	pwd VARCHAR(250) NOT NULL,
	active %s DEFAULT true NOT NULL,
	set_pwd %s DEFAULT false NOT NULL,
	created_at %s NOT NULL,
	updated_at %s NULL,
	PRIMARY KEY(id)
);`, uuid, b, b, ts, ts)

	resetScript := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS login_reset (
	login_id %s NOT NULL,
	reset_token %s NOT NULL,
	created_at %s NOT NULL,
	updated_at %s NULL,
	PRIMARY KEY(login_id, reset_token)
);`, uuid, uuid, ts, ts)

	roleScript := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS role (
	id %s NOT NULL,
	name VARCHAR(50) NOT NULL,
	description VARCHAR(500) NULL,
	PRIMARY KEY(id)
);`, uuid)

	loginRoleScript := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS login_role (
	login_id %s NOT NULL,
	role_id %s NOT NULL,
	PRIMARY KEY(login_id, role_id)
);`, uuid, uuid)

	authAuthorizeScript := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS auth_authorize (
		id VARCHAR(32) NOT NULL,
		client_id VARCHAR(32) NOT NULL,
		verifier TEXT,
		verifier_encode_method VARCHAR(10),
		state VARCHAR(100),
		scope VARCHAR(256),
		authorized_at %s NOT NULL,
		auth_code_at %s,
		auth_code VARCHAR(256),
		PRIMARY KEY(id)
);`, ts, ts)

	authClientScript := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS auth_client (
		id VARCHAR(32) NOT NULL,
		name VARCHAR(100) NOT NULL,
		description VARCHAR(1000),
		homepage_url VARCHAR(500) NOT NULL,
		callback_url VARCHAR(500) NOT NULL,
		PRIMARY KEY(id)
);`)

	authClientSecretScript := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS auth_client_secret (
		client_id VARCHAR(32) NOT NULL,
		secret VARCHAR(256) NOT NULL,
		PRIMARY KEY(client_id, secret)
);`)

	authRefreshScript := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS auth_refresh (
		client_id VARCHAR(32) NOT NULL,
		token VARCHAR(256) NOT NULL,
		created_at %s NOT NULL,
		PRIMARY KEY(client_id, token)
);`, ts)

	adminId := util.GenerateUUID()
	userId := util.GenerateUUID()
	roleInsert := fmt.Sprintf(`INSERT INTO role (id, name, description) VALUES
	('%s', 'admin', 'System admin - this should be limited'),
	('%s', 'user', 'Default user');`, adminId, userId)

	scriptDir := fmt.Sprintf("%s/scripts/migrations", p.ProjectFile.FullPath)
	if err := os.MkdirAll(scriptDir, os.ModePerm); err != nil {
		fmt.Println("Creating scripts/migrations dir", err)
		return
	}
	scripts := []ScriptFile{
		{"%s/%s-create-table-login.sql", loginScript},
		{"%s/%s-create-table-reset-login.sql", resetScript},
		{"%s/%s-create-table-role.sql", roleScript},
		{"%s/%s-create-table-login-role.sql", loginRoleScript},
		{"%s/%s-insert-role.sql", roleInsert},
		{"%s/%s-create-table-auth-authorize.sql", authAuthorizeScript},
		{"%s/%s-create-table-auth-client.sql", authClientScript},
		{"%s/%s-create-table-auth-client-secret.sql", authClientSecretScript},
		{"%s/%s-create-table-auth-refresh.sql", authRefreshScript},
	}
	for _, script := range scripts {
		now := time.Now().Format("20060102150405")
		loginName := fmt.Sprintf(script.FileName, scriptDir, now)
		f, err := os.Create(loginName)
		if err != nil {
			fmt.Printf("Unable to creating %s file: %s\n", loginName, err)
			return
		}
		f.WriteString(script.Script)
		f.Close()
		time.Sleep(time.Second) // let's make sure a second has passed
	}
	// note: this has to run last
	// compile the admin tool, move binary to scripts/migrations
	now := time.Now().Format("20060102150405")
	execDest := fmt.Sprintf("-o=%s/%s-admin-tool.bin", scriptDir, now)
	execAdmin := exec.Command("go", "build", execDest, "tools/admin/main.go")
	output, errAdminCmd := execAdmin.CombinedOutput()
	fmt.Printf("%s\n", output)
	if errAdminCmd != nil {
		fmt.Println("Unable to create admin tool file:", errAdminCmd)
	}
}

func buildRoleEntities(project *m.Project) {
	uuid := "UUID"
	if project.ProjectFile.SqlStorage == "m" {
		// mysql
		uuid = "CHAR(36)"
	}
	if project.ProjectFile.SqlStorage == "s" {
		// sqlite
		uuid = "TEXT"
	}
	// role
	idName := m.Name{}
	idName.BuildName("id", []string{})
	nameName := m.Name{}
	nameName.BuildName("name", []string{})
	descName := m.Name{}
	descName.BuildName("description", []string{})
	roleName := m.Name{}
	rName := roleName.BuildName("role", project.KnownAliases)
	project.KnownAliases = append(project.KnownAliases, rName)
	role := m.Entity{
		Columns: []m.Column{
			{
				ColumnName:   idName,
				DBType:       uuid,
				GoType:       "string",
				GoTypeNonSql: "string",
				Null:         false,
				PrimaryKey:   true,
			},
			{
				ColumnName:   nameName,
				DBType:       "varchar",
				GoType:       "null.String",
				GoTypeNonSql: "string",
				Length:       50,
				Null:         false,
				PrimaryKey:   false,
			},
			{
				ColumnName:   descName,
				DBType:       "varchar",
				GoType:       "null.String",
				GoTypeNonSql: "string",
				Length:       500,
				Null:         true,
				PrimaryKey:   false,
			},
		},
		Name:            roleName,
		ColumnExistence: m.ColumnExistence{HaveNullColumns: true},
		ModuleName:      "role",
	}
	project.Entities = append(project.Entities, role)
	// login_role
	loginIdName := m.Name{}
	loginIdName.BuildName("login_id", []string{})
	roleIdName := m.Name{}
	roleIdName.BuildName("role_id", []string{})
	loginRoleName := m.Name{}
	lrName := loginRoleName.BuildName("login_role", project.KnownAliases)
	project.KnownAliases = append(project.KnownAliases, lrName)
	loginRole := m.Entity{
		Columns: []m.Column{
			{
				ColumnName:   loginIdName,
				DBType:       uuid,
				GoType:       "string",
				GoTypeNonSql: "string",
				Null:         false,
				PrimaryKey:   true,
			},
			{
				ColumnName:   roleIdName,
				DBType:       uuid,
				GoType:       "string",
				GoTypeNonSql: "string",
				Null:         false,
				PrimaryKey:   true,
			},
		},
		Name:            loginRoleName,
		ColumnExistence: m.ColumnExistence{HaveNullColumns: false},
		ModuleName:      "loginrole",
	}
	project.Entities = append(project.Entities, loginRole)
}

func Readme(p m.Project) {
	readmeLines := `
**Auth**: provides the ` + "`rest`" + ` service to authenticate via JWT token.

Additional features:
- ` + "`login`" + `, ` + "`role`" + ` and ` + "`login-role`" + ` endpoints for CRUD
- creation/validation of JWT token for authentication
- authorization of endpoints per set of ` + "`role`" + `s
- ` + "`login`" + `: sign-in and forget-password with send email logic
- initialize admin user tool

The following env vars are needed to be supplied for the **auth** feature
` + "`{{.ProjectNameEnv}}_LOGIN_PWD_COST`" + `: [int] the cost the encryption algorithm needs, 
` + "`{{.ProjectNameEnv}}_LOGIN_RESET_DURATION`" + `: [int] the value in days after the reset record will expire
` + "`{{.ProjectNameEnv}}_LOGIN_EXPIRES_AT_DURATION`" + `: [int] the value when the JWT will expire
` + "`{{.ProjectNameEnv}}_LOGIN_AUTH_ALG`" + `: [string] encryption algorithm name ` + "`HMAC | RSA | ECDSA | EdDSA`" + ` default is ` + "`HMAC`" + `, see ` + "`JWT Signing`" + ` for more info
` + "`{{.ProjectNameEnv}}_LOGIN_AUTH_PRIVATE`" + `: [string] base64 encoded private key for ` + "`RSA, ECDSA or EdDSA`" + ` or base64 encoded password for ` + "`HMAC`" + `, see ` + "`JWT Signing`" + ` for more info
` + "`{{.ProjectNameEnv}}_LOGIN_AUTH_PUBLIC`" + `: [string] base64 encoded public key for ` + "`RSA, ECDSA or EdDSA`" + ` should match the private file, see ` + "`JWT Signing`" + ` for more info
` + "`{{.ProjectNameEnv}}_LOGIN_BASIC_AUTH_USER`" + `: [string] optional user name for basic auth
` + "`{{.ProjectNameEnv}}_LOGIN_BASIC_AUTH_PASS`" + `: [string] optional user password for basic auth
` + "`{{.ProjectNameEnv}}_EMAIL_PORT`" + `: [int] port for the stmp service
` + "`{{.ProjectNameEnv}}_EMAIL_PWD`" + `: [string] password for the stmp service
` + "`{{.ProjectNameEnv}}_EMAIL_FROM`" + `: [string] user name for the stmp service
` + "`{{.ProjectNameEnv}}_EMAIL_RESET_URL`" + `: [string] from email address used for the reset-password email
` + "`{{.ProjectNameEnv}}_ADMIN_EMAIL`" + `: [string] the admin's email address for the initialization tool

` + "`NOTE`" + `: the nature of env vars, the value that are of type ` + "`int`" + ` or ` + "`bool`" + ` are treat as string but the code will try to cast them as the type they will need to be.  Putting quotes around these value is safe and encouraged.

#### JWT Signing
Depending on the type of algorithm you set, the following are examples to create ` + "`.pem`" + ` files for the authentication/login feature.  These can be removed, just here as helper commands.

The contents of these files will need to be ` + "`base64`" + ` encoded before you save them to the env vars.  This helps to minimize the extra lines with your profile/rc file.

##### RSA
` + "```" + `
# private key
openssl genpkey -algorithm rsa -out rsa-private.pem
# public key
openssl pkey -in rsa-private.pem -pubout -out rsa-public.pem
` + "```" + `
##### ECDSA
` + "```" + `
# private key
openssl ecparam -genkey -name secp521r1 -noout -out es521-private.pem
# public key
openssl ec -in es521-private.pem -pubout -out es521-public.pem 
` + "```" + `
##### EdDSA
` + "```" + `
# private key
openssl genpkey -algorithm ed25519 -out eddsa-private.pem
# public key
openssl pkey -in eddsa-private.pem -pubout -out eddsa-public.pem
` + "```" + `
`

	readmeFile := fmt.Sprintf("%s/README.md", p.ProjectFile.FullPath)
	if _, err := os.Stat(readmeFile); os.IsNotExist(err) {
		fmt.Printf("%s is missing unable to write in hooks\n", readmeFile)
	} else {
		var readmeReplace bytes.Buffer
		tReadme := template.Must(template.New("readme").Parse(readmeLines))
		errReadme := tReadme.Execute(&readmeReplace, p)
		if errReadme != nil {
			fmt.Printf("%s: template error [%s]\n", readmeFile, errReadme)
			return
		}
		file, err := os.OpenFile(readmeFile, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Error opening the %s file: %v\n", readmeFile, err)
			return
		}
		defer file.Close()
		if _, err := file.WriteString(readmeReplace.String()); err != nil {
			fmt.Println("Error writing to the readme file")
			return
		}
	}
}
