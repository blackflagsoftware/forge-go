package projectfile

import (
	n "github.com/blackflagsoftware/forge-go/internal/name"
)

type (
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
		DynamicSchema        bool     `json:"dynamic_schema"`         // TODO: might not need anymore
		Schema               string   `json:"schema"`                 // TODO: might not need anymore
		DynamicSchemaPostfix string   `json:"dynamic_schema_postfix"` // TODO: might not need anymore
		UseORM               bool     `json:"use_orm"`                // TODO: make an option to change this
		TagFormat            string   `json:"tag_format"`             // TODO: make an option to change this
		KnownAliases         []string `json:"known_aliases"`
		VersionPath          string   `json:"version_path"` // TODO: make an option to change this
		Modules              []string `json:"modules"`
		LoadedConfig         bool     `json:"loaded_config"` // tells the code to only run once
		n.Name
	}
)
