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
)

func AddLogin(p *m.Project) {
	Config(*p)
	Errors(*p)
	ProtoFile(*p)
	TemplateFiles(*p)
	loginEntity := m.Entity{ModuleName: "login"}
	loginEntity.Name.BuildName("login", p.ProjectFile.KnownAliases)
	p.Entities = append(p.Entities, loginEntity)
	temp.StartTemplating(p)
	// BuildStorage(*p)
	// loginEntity.BuildAPIHooks()
	// UpdateModFiles(p.ProjectFile.AppName)
	MigrationScripts(*p)
	p.ProjectFile.Modules = append(p.ProjectFile.Modules, "login")
	p.ProjectFile.SaveProjectFile()
	// p.ProjectFile.LoadProjectFile()
	// TODO: this is a mess, it all comes down with passing ProjectFile around everywhere!
}

func Config(p m.Project) {
	configLine := `LoginPwdCost           = GetEnvOrDefault("{{.ProjectNameUpper}}_LOGIN_PWD_COST", "10")             \/\/ algorithm cost
	LoginResetDuration     = GetEnvOrDefault("{{.ProjectNameUpper}}_LOGIN_RESET_DURATION", "7")        \/\/ in days
	LoginExpiresAtDuration = GetEnvOrDefault("{{.ProjectNameUpper}}_LOGIN_EXPIRES_AT_DURATION", "168") \/\/ in hours (7 days)
	LoginAuthSecret        = GetEnvOrDefault("{{.ProjectNameUpper}}_LOGIN_AUTH_SECRET", "")
	LoginEmailHost         = GetEnvOrDefault("{{.ProjectNameUpper}}_EMAIL_HOST", "")
	LoginEmailPort         = GetEnvOrDefault("{{.ProjectNameUpper}}_EMAIL_PORT", "587")
	LoginEmailPwd          = GetEnvOrDefault("{{.ProjectNameUpper}}_EMAIL_PWD", "")
	LoginEmailFrom         = GetEnvOrDefault("{{.ProjectNameUpper}}_EMAIL_FROM", "")
	LoginEmailResetUrl     = GetEnvOrDefault("{{.ProjectNameUpper}}_EMAIL_RESET_URL", "")
	LoginAdminEmail        = GetEnvOrDefault("{{.ProjectNameUpper}}_ADMIN_EMAIL", "")
	\/\/ --- replace config text - do not remove ---
	`
	configFile := fmt.Sprintf("%s/config/config.go", p.ProjectFile.FullPath)
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		fmt.Printf("%s is missing unable to write in hooks\n", configFile)
	} else {
		var configReplace bytes.Buffer
		tConfig := template.Must(template.New("config").Parse(configLine))
		errConfig := tConfig.Execute(&configReplace, p)
		if errConfig != nil {
			fmt.Printf("%s: template error [%s]\n", configFile, errConfig)
		} else {
			cmdConfig := fmt.Sprintf(`perl -pi -e 's/\/\/ --- replace config text - do not remove ---/%s/g' %s`, configReplace.String(), configFile)
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
}`

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
	string Uid = 1;
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
	string Uid = 1;
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
			fmt.Printf("Template could not parse login file: %s; %s\n", f.Src, errParse)
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
	uid %s NOT NULL,
	email_addr VARCHAR(100) NOT NULL,
	pwd VARCHAR(250) NOT NULL,
	active %s DEFAULT true NOT NULL,
	set_pwd %s DEFAULT false NOT NULL,
	created_at %s NOT NULL,
	updated_at %s NULL,
	PRIMARY KEY(uid)
);`, uuid, b, b, ts, ts)

	resetScript := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS reset_login (
	login_uid %s NOT NULL,
	reset_token %s NOT NULL,
	created_at %s NOT NULL,
	updated_at %s NULL,
	PRIMARY KEY(login_uid, reset_token)
);`, uuid, uuid, ts, ts)

	scriptDir := fmt.Sprintf("%s/scripts/migrations", p.ProjectFile.FullPath)
	if err := os.MkdirAll(scriptDir, os.ModePerm); err != nil {
		fmt.Println("Creating scripts/migrations dir", err)
		return
	}
	now := time.Now().Format("20060102150405")
	loginName := fmt.Sprintf("%s/%s-create-table-login.sql", scriptDir, now)
	f, err := os.Create(loginName)
	if err != nil {
		fmt.Println("Unable to creating scripts/migrations/login.sql file", err)
		return
	}
	f.WriteString(loginScript)
	f.Close()
	time.Sleep(time.Second) // let's make sure a second has passed
	now = time.Now().Format("20060102150405")
	resetName := fmt.Sprintf("%s/%s-create-table-reset-login.sql", scriptDir, now)
	f, err = os.Create(resetName)
	if err != nil {
		fmt.Println("Unable to create scripts/migrations/reset-login.sql file", err)
		return
	}
	f.WriteString(resetScript)
	f.Close()
	time.Sleep(time.Second) // let's make sure a second has passed
	// compile the admin tool, move binary to scripts/migrations
	now = time.Now().Format("20060102150405")
	execDest := fmt.Sprintf("-o=%s/%s-admin-tool.bin", scriptDir, now)
	execAdmin := exec.Command("go", "build", execDest, "tools/admin/main.go")
	output, errAdminCmd := execAdmin.CombinedOutput()
	fmt.Printf("%s\n", output)
	if errAdminCmd != nil {
		fmt.Println("Unable to create admin tool file:", errAdminCmd)
	}
}
