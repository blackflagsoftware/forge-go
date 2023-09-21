package project

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"text/template"

	e "github.com/blackflagsoftware/forge-go/internal/entity"
)

func (p *Project) AddLogin() {
	p.Config()
	p.Errors()
	p.ProtoFile()
	p.TemplateFiles()
	e.PopulateConfig(&p.ProjectFile)
	BuildStorage(*p)
	UpdateModFiles(p.ProjectFile.AppName)
}

func (p Project) Config() {
	configLine := `LoginPwdCost           = GetEnvOrDefault("{{.Upper}}_LOGIN_PWD_COST", "10")             \/\/ algorithm cost
	LoginResetDuration     = GetEnvOrDefault("{{.Upper}}_LOGIN_RESET_DURATION", "7")        \/\/ in days
	LoginExpiresAtDuration = GetEnvOrDefault("{{.Upper}}_LOGIN_EXPIRES_AT_DURATION", "168") \/\/ in hours (7 days)
	LoginAuthSecret        = GetEnvOrDefault("{{.Upper}}_LOGIN_AUTH_SECRET", "")
	LoginEmailHost         = GetEnvOrDefault("{{.Upper}}_EMAIL_HOST", "")
	LoginEmailPort         = GetEnvOrDefault("{{.Upper}}_EMAIL_PORT", "587")
	LoginEmailPwd          = GetEnvOrDefault("{{.Upper}}_EMAIL_PWD", "")
	LoginEmailFrom         = GetEnvOrDefault("{{.Upper}}_EMAIL_FROM", "")
	LoginEmailResetUrl     = GetEnvOrDefault("{{.Upper}}_EMAIL_RESET_URL", "")
	LoginAdminEmail        = GetEnvOrDefault("{{.Upper}}_ADMIN_EMAIL", "")
	\/\/ --- replace config text - do not remove ---
	`
	configFile := fmt.Sprintf("%s/config/config.go", p.ProjectFile.FullPath)
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		fmt.Printf("%s is missing unable to write in hooks\n", configFile)
	} else {
		var configReplace bytes.Buffer
		tConfig := template.Must(template.New("config").Parse(configLine))
		name := p.ProjectFile.Name
		errConfig := tConfig.Execute(&configReplace, name)
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

func (p Project) Errors() {
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

func (p Project) ProtoFile() {
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

func (p Project) TemplateFiles() {
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
		{fmt.Sprintf("%s/modules/login/v1_login/grpc.tmpl", os.Getenv("FORGE_PATH")), fmt.Sprintf("%s/internal/v1/login/grpc.go", p.ProjectFile.FullPath), true},
		{fmt.Sprintf("%s/modules/login/v1_login/manager_test.tmpl", os.Getenv("FORGE_PATH")), fmt.Sprintf("%s/internal/v1/login/manager_test.go", p.ProjectFile.FullPath), false},
		{fmt.Sprintf("%s/modules/login/v1_login/manager.tmpl", os.Getenv("FORGE_PATH")), fmt.Sprintf("%s/internal/v1/login/manager.go", p.ProjectFile.FullPath), false},
		{fmt.Sprintf("%s/modules/login/v1_login/model.tmpl", os.Getenv("FORGE_PATH")), fmt.Sprintf("%s/internal/v1/login/model.go", p.ProjectFile.FullPath), false},
		{fmt.Sprintf("%s/modules/login/v1_login/rest.tmpl", os.Getenv("FORGE_PATH")), fmt.Sprintf("%s/internal/v1/login/rest.go", p.ProjectFile.FullPath), false},
		{fmt.Sprintf("%s/modules/login/v1_login/sql.tmpl", os.Getenv("FORGE_PATH")), fmt.Sprintf("%s/internal/v1/login/sql.go", p.ProjectFile.FullPath), false},
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
		err = t.Execute(file, p.ProjectFile)
		if err != nil {
			fmt.Println("Execution of template:", err)
		}
	}
}
