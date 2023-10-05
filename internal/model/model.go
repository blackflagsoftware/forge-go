package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path"
	"strings"
	"unicode"

	"github.com/Masterminds/sprig"
)

type (
	Project struct {
		UseBlank              bool
		Entities              []Entity
		ProjectNameAllLower   string
		ProjectNameLower      string
		ProjectNameCamel      string
		ProjectNameAbbr       string
		ProjectNameLowerCamel string
		ProjectNameUpper      string
		ProjectNameEnv        string
		// SQLProvider           string // optional if using SQL as a storage, either Psql, MySql or Sqlite; this interfaces with sqlx
		// SQLProviderLower      string // optional is using SQL, lowercase of above
		CurrentEntity Entity // rotating entity or templating, will be emptied out after processing
		ProjectFile
		FileTemplate
		GrpcTemplate
		ManagerTemplate
		ModelTemplate
		RestTemplate
		StorageTemplate
		StorageVars
		MigrationVars
	}

	FileTemplate struct {
		FileKeys       string
		FileGetColumns string
		FilePostIncr   string
	}

	GrpcTemplate struct {
		GrpcImport       string // this needs to blank out
		GrpcArgsInit     string
		GrpcTranslateIn  string
		GrpcTranslateOut string
	}

	StorageVars struct {
		SQLProvider           string // optional if using SQL as a storage, either Psql, MySql or Sqlite; this interfaces with sqlx
		SQLProviderLower      string // optional if using SQL as a storage, either psql, mysql or sqlite; this interfaces with gorm
		SQLProviderConnection string // holds the connection string for gorm of the other sql types
	}

	// rename to Template
	MigrationVars struct {
		MigrationHeader     string
		MigrationConnection string
		MigrationVerify     string
	}

	Entity struct {
		Columns       []Column
		SqlLines      []string
		DefaultColumn string
		SortColumns   string
		ModuleName    string
		Name
		ColumnExistence
	}

	ModelTemplate struct {
		ModelImport      string
		ModelRows        string
		ModelInitStorage string
	}

	ManagerTemplate struct {
		ManagerTestImport    string
		ManagerTestGetRow    string
		ManagerTestPostRow   string
		ManagerTestPatchInit string
		ManagerTestDeleteRow string
		ManagerImport        string
		ManagerGetRows       string
		ManagerPostRows      string
		ManagerPatchInitArgs string
		ManagerPatchRows     string
		ManagerAuditKey      string
	}

	RestTemplate struct {
		RestStrConv         string
		RestGetDeleteUrl    string
		RestGetDeleteAssign string
		RestArgSet          string
	}

	StorageTemplate struct {
		StorageTable          string
		StorageTablePrefix    string
		StorageTablePostfix   string
		StorageGetColumns     string
		StorageTableKeyKeys   string
		StorageTableKeyValues string
		// StorageTableKeyListOrder string
		StoragePostColumns      string
		StoragePostColumnsNamed string
		StoragePostReturning    string
		StoragePostQuery        string
		StoragePostLastId       string
		StoragePatchColumns     string
		StoragePatchWhere       string
		StoragePatchWhereValues string
	}

	PostPutTest struct {
		Name         string
		ForColumn    string
		ColumnLength int
		Failure      bool
	}

	ColumnTest struct {
		GoType string
		DBType string
	}

	ProjectFile struct {
		Message              string   `json:"message"`
		AppName              string   `json:"app_name"`
		FullPath             string   `json:"full_path"`    // full path to the project
		SubDir               string   `json:"sub_dir"`      // directory, only, you will save the files to
		SubPackage           string   `json:"sub_package"`  // if sub directory is multipath i.e. internal/v1, will be used for package name
		ProjectPath          string   `json:"project_path"` // fullpath minus gopath/src; used for import statements
		Storage              string   `json:"storage"`
		SqlStorage           string   `json:"sql_storage"`
		ProjectPathEncoded   string   `json:"project_path_encoded"`   // encode to use in some of the templating
		SubDirEncoded        string   `json:"sub_dir_encoded"`        // encode to use in some of the templating
		DynamicSchema        bool     `json:"dynamic_schema"`         // TODO: this is for postgres, need to add an ability to add it
		Schema               string   `json:"schema"`                 // TODO: this is for postgres, need to add an ability to add it
		DynamicSchemaPostfix string   `json:"dynamic_schema_postfix"` // TODO: might not need anymore
		UseORM               bool     `json:"use_orm"`
		TagFormat            string   `json:"tag_format"`
		KnownAliases         []string `json:"known_aliases"`
		VersionPath          string   `json:"version_path"` // TODO: make an option to change this
		Modules              []string `json:"modules"`
		LoadedConfig         bool     `json:"loaded_config"` // tells the code to only run once
		Name
	}

	Column struct {
		ColumnName   Name
		DBType       string
		GoType       string
		GoTypeNonSql string
		Null         bool
		DefaultValue string
		Length       int64
		PrimaryKey   bool
	}

	ColumnExistence struct {
		HaveNullColumns bool
		TimeColumn      bool
	}

	Name struct {
		RawName    string `json:"-"` // name given by the user/script
		Lower      string `json:"-"`
		Camel      string `json:"-"`
		LowerCamel string `json:"-"`
		Abbr       string `json:"-"`
		AllLower   string `json:"-"`
		Upper      string `json:"-"`
		EnvVar     string `json:"-"`
	}

	Format struct {
		Name string
	}

	SqlEntity struct {
		Name         string
		Columns      []Column
		ColExistence ColumnExistence
	}
)

func (n *Name) BuildName(name string, knownAliases []string) string {
	n.RawName = name
	runeName := []rune(n.RawName)
	idxFirst := 0
	idxLast := len(runeName)
	for i := range runeName {
		if !unicode.IsLetter(runeName[i]) {
			idxFirst = i + 1
			continue
		}
		break
	}
	for i := len(runeName) - 1; i > 0; i-- {
		if !unicode.IsLetter(runeName[i]) {
			idxLast = i
			continue
		}
		break
	}

	rawName := string(runeName[idxFirst:idxLast])

	n.Camel = BuildAltName(rawName, "pascalCase")
	n.LowerCamel = BuildAltName(rawName, "camelCase")
	n.Lower = BuildAltName(rawName, "lowerCase")
	n.AllLower = strings.ToLower(n.Camel)
	n.Upper = strings.ToUpper(n.Camel)
	n.EnvVar = strings.ReplaceAll(strings.ToUpper(rawName), "-", "_")
	return n.DetermineAbbr(knownAliases)
}

func BuildAltName(name, mode string) string {
	f := Format{Name: name}
	var t *template.Template
	var err error
	switch mode {
	case "snakeCase":
		t, err = template.New("format").Funcs(sprig.GenericFuncMap()).Parse("{{.Name | snakecase}}")
	case "kebabCase":
		t, err = template.New("format").Funcs(sprig.GenericFuncMap()).Parse("{{.Name | kebabcase}}")
	case "camelCase":
		t, err = template.New("format").Funcs(sprig.GenericFuncMap()).Parse("{{.Name | camelcase}}")
	case "pascalCase":
		t, err = template.New("format").Funcs(sprig.GenericFuncMap()).Parse("{{.Name | camelcase}}")
	case "upperCase":
		t, err = template.New("format").Funcs(sprig.GenericFuncMap()).Parse("{{.Name | upper}}")
	default:
		// lowerCase
		t, err = template.New("format").Funcs(sprig.GenericFuncMap()).Parse("{{.Name | lower}}")
	}
	if err != nil {
		fmt.Println(err)
		return name
	}
	b := bytes.NewBufferString("")
	errE := t.Execute(b, f)
	if errE != nil {
		fmt.Println(errE)
		return name
	}
	if mode == "camelCase" {
		// finish off the camel case functionality
		n := []rune(b.String())
		return string(append([]rune{unicode.ToLower(n[0])}, n[1:]...))
	}
	return b.String()
}

func (n *Name) DetermineAbbr(knownAliases []string) string {
	// check all known aliases and create the best Name.Abbr
	if len(n.Lower) < 3 {
		n.Abbr = n.Lower
		return n.Abbr
	}
	try := 0
Loop:
	for {
		try++
		var retry bool
		n.Abbr, retry = nameVersion(try, n.Lower)
		if !retry {
			break
		}
		for _, ka := range knownAliases {
			if ka == n.Abbr {
				continue Loop
			}
		}
		break
	}
	return n.Abbr
}

func nameVersion(tried int, name string) (abbr string, retry bool) {
	switch tried {
	case 1:
		// first 3 character
		abbr = string(name[:3])
		retry = true
	case 2:
		// first 2 characters
		abbr = string(name[:2])
		retry = true
	case 3, 4:
		// first letters of each snake case
		charNumber := tried - 2
		split := strings.Split(name, "_")
		if len(split) == 1 {
			abbr = name
			return
		}
		for _, c := range split {
			abbr += string(c[:charNumber])
		}
		retry = true
	default:
		abbr = name
		retry = false
	}
	return
}

func (p *ProjectFile) CreateProjectFile(pwd string) bool {
	p.Message = "This is used for forge program for convenience"
	p.FullPath = pwd
	p.VersionPath = "v1"
	p.AppName = path.Base(pwd)
	p.Name.BuildName(p.AppName, []string{})
	p.SubDir = "internal/" + p.VersionPath

	// create projectpath and subpackage
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		fmt.Println("***** GOPATH is not set, some of your paths will not be correct! *****")
	} else {
		idx := len(fmt.Sprintf("%s/src/", goPath))
		p.ProjectPath = p.FullPath[idx:]
	}
	subPath := strings.Split(p.SubDir, "/")
	p.SubPackage = subPath[len(subPath)-1]
	// encode paths
	p.ProjectPathEncoded = strings.Replace(p.ProjectPath, `/`, `\/`, -1)
	p.SubDirEncoded = strings.Replace(p.SubDir, `/`, `\/`, -1)
	return true
}

func (p *ProjectFile) LoadProjectFile() bool {
	if _, errStat := os.Stat("./.forge"); os.IsNotExist(errStat) {
		return false
	}
	bContent, errRead := os.ReadFile("./.forge")
	if errRead != nil {
		fmt.Printf("Error reading from .forge: %s", errRead)
		return true
	}
	errUnmarshal := json.Unmarshal(bContent, &p)
	if errUnmarshal != nil {
		fmt.Printf("Error extracting data from .forge: %s", errUnmarshal)
		return true
	}
	return true
}

func (p *ProjectFile) SaveProjectFile() {
	bContent, errMarshal := json.MarshalIndent(p, "", "    ")
	if errMarshal != nil {
		fmt.Println("Saving project file: unable to save -", errMarshal)
		return
	}
	errSave := os.WriteFile("./.forge", bContent, 0644)
	if errSave != nil {
		fmt.Println("Saving project file: unable to save -", errSave)
	}
}

func StorageTypeToProper(abbr string) string {
	switch abbr {
	case "s":
		return "SQL"
	case "f":
		return "File"
	case "m":
		return "MongoDB"
	}
	return abbr
}

func SqlTypeToProper(abbr string) string {
	switch abbr {
	case "p":
		return "Postgres"
	case "m":
		return "MySQL"
	case "s":
		return "Sqlite"
	}
	return abbr
}

func TagFormatToProper(tag string) string {
	switch tag {
	case "snakeCase":
		return "Snake Case"
	case "kebabCase":
		return "Kebab Case"
	case "camelCase":
		return "Camel Case"
	case "pascalCase":
		return "Pascal Case"
	case "upperCase":
		return "Upper Case"
	case "lowerCase":
		return "Lower Case"
	}
	return tag
}

func (e *Entity) HasNullColumn() bool {
	for _, c := range e.Columns {
		if strings.Contains(c.GoType, "null") {
			return true
		}
	}
	return false
}

func (e *Entity) HasTimeColumn() bool {
	for _, c := range e.Columns {
		if c.GoType == "time.Time" || c.ColumnName.Lower == "created_at" || c.ColumnName.Lower == "updated_at" {
			return true
		}
	}
	return false
}

func (e *Entity) HasNullTimeColumn() bool {
	for _, c := range e.Columns {
		if c.GoType == "null.Time" {
			return true
		}
	}
	return false
}

func (e *Entity) HasJsonColumn() bool {
	for _, c := range e.Columns {
		if c.GoType == "*json.RawMessage" {
			return true
		}
	}
	return false
}

func (e *Entity) HasPrimaryUUIDColumn() bool {
	for _, c := range e.Columns {
		if c.DBType == "uuid" && c.PrimaryKey {
			return true
		}
	}
	return false
}
