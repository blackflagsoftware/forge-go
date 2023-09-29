# Migration

This tool is used to help this particular api/microservice to migrate database updates and changes if so desired.

### Setup

##### EnvVars
`FORGE_GO_BASE_MIGRATION_ENABLED`: [true/false] set to `true` to enable it when your api/service starts up, see `Integration` for more details.

`FORGE_GO_BASE_MIGRATION_PATH`: [string] path to the migration scripts, by default it will use the `<project_folder>/scripts/migrations`. Leaving this to the default value may be useful during debugging or running on a local development maching. A different directory path may be provided as needed, this is particularly true when doing something like containerization and deploying to the cloud.

`FORGE_GO_BASE_MIGRATION_SKIP_INIT`: [true/false] set to `true` if there is another outside process that will create the DB and migration table initialization process.

`FORGE_GO_BASE_ADMIN_DB_USER`: [string] (optional) the root admin user name in case the normal user (`FORGE_GO_BASE_DB_USER`) doesn't not have access to create a DB or tables

`FORGE_GO_BASE_ADMIN_DB_PASS`: [string] (optional) the root admin password in case the normal user (`FORGE_GO_BASE_DB_PASS`) doesn't not have access to create a DB or tables

The following env vars are copied from the project's main `config.go` file and should be set for the project's connection to the DB engine, mentioning this because this tool is dependent on these env vars:

`FORGE_GO_BASE_DB_HOST`
`FORGE_GO_BASE_DB_DB`
`FORGE_GO_BASE_DB_USER`
`FORGE_GO_BASE_DB_PASS`
`FORGE_GO_BASE_MIGRATION_DB_ENGINE`

Note: The prefix of the above env vars will change in the `config.go` file, the prefix will change to the new project name.

##### Integration
If `FORGE_GO_BASE_MIGRATION_ENABLED` is set then you service will automatically run the code within this tool to process all the migration scripts to date.

##### Stand-alone
This tool can be compiled as a normal Golang binary and executed in stand-alone mode.  There are a few reasons you will want to do this:

- Normalize your script names, (see `NormalizeNames` and `Usage`)
- Interactive script creation process, (see `InteractiveMode`)
- Process migration scripts (for debugging, etc)

#### Scripts
The script shoule be composed of the SQL language statement.  Due to limitations on some of the sql libraries in use, having only one single SQL statement per script, i.e. `create table` statement and populating that table should be separate files, files should end in *.sql, see `NormalizeNames` for more details.

#### Executables
The tool will also allow you to run a more complex set of instructions via a separate binary.  Warning, you will have to compile and make available that file in your process of deployment, files should end in *.bin and should be in normalized name format, see `NormalizeNames` for more details.

#### NormalizeNames
In order for this tool to run each script just once, the script/binary files will need to have a normalized name, in this case the file needs to start with a date format of YYYYMMDDhhmmss, i.e. 20230213010559-create-my-first-table.sql.  This will put the files in order they need to run.  The rest of the name is only for human readablity and is up to the user on what that is.  Though the normalization process will replace any `"_" or " "` to `-`, just to keep the file names consistent.

### Usage

During development the following arguments can be used:
- `f`: file name, found in `<project_name>/scripts/migrations`, this will normalize the file name.  This is `rename` process of that file.
- `m`: this will run the migration process
- `h`: host name or the sql server name/port
- `d`: database name
- `u`: migration user name
- `p`: migration password
- `au`: admin user name
- `ap`: admin password
- `t`: script directory

The arguments `m, h, d, u, p, au, ap and t` are optional and mostly used for development/debugging overrides.  For deployment/production uses, the env vars in the `<project_path>/config/config.go` file should point you to which env vars need to be provided for your deployment process.

#### InteractiveMode
If `f or m` is not provided, then the tool will go into interactive mode and will ask you to:
- enter your sql code
- enter a descriptive name

This will save the file in `NormalizeName` format in the `<project_name>/scripts/migrations` folder.

##### TODO
Rollback