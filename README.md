# forge-go
Clone, build and deploy an API written in Golang, quick and easy.  A base framework to add entities to create standard CRUD endpoints based on a SQL create table statement.  From there, customize as needed.

### Go Versions
Built with Go 1.17+, verified 1.18

### Install
```
get clone github.com/blackflagsoftware/forge-go
cd forge-go
go mod tidy
bash ./build_forge.sh // helper script to put the binary in the root directory
```
Then you will need to add an `env var` to your system called `FORGE_PATH` that points to the project's root path.

### Usage
To clone a new project:
From anywhere in terminal or console (as long it can be found via the `PATH` list, see `Install`)
```
forge clone <path>
```
`<path>` is full path from `$GOPATH/src` e.g.: `github.com/blackflagsoftware/new_project`, if the folder doesn't exists, the folder will be created and the project will be copied into this folder.  If the folder exists, bad things happen, kidding, and message will show indicating it exists and will exit.

With a cloned project, go ahead and load it into your favorite text editor.

To add an entity and create your CRUD code:
From the root directory of the new project, run this in terminal or console
```
forge
```
Through a console CLI program `forge` will ask a few starter questions:
- correct path to start the project
- which storage type to create code?
	- SQL
	- File
	- MongoDB
- If SQL is chosen then which engine?
	- Postgres
	- Mysql
	- Sqlite
- Tag Format to use for the 'json' struct tags

Once this is completed then a `.forge` file is saved in the project's root directory.  With this in place then these questions will not be asked again.

The normal flow when starting this CLI application within your project's root directory will ask questions to add a new endpoint based on SQL syntax.

How to add your *object* to the project.  The next menu asks this:
- File as input
- Paste as input
- Prompt as input
- Blank Structure
- Admin

`File as input` as suggested, it will prompt for a file to load, this can be a file of multiple SQL create table statements, `forge` will create CRUD for each of the tables

`Paste as input` as suggested, this will prompt to paste in the SQL create table statement.  It prompts for multiple statements if needed. **NOTE** when pasting in a `create table` statement be aware NOT to have any blank lines between any of the other lines.  This will cause the prompt parser not to get the whole statement, and the process will not create the boilerplate code.

`Prompt as input`, this will prompt for entity name and then add as many fields as you need.  It will prompt for column type, is it null, etc.  When finished it will save the SQL to `./prompt_schema` for your reference.

`Blank Structure`, this will create all the files as if it was creating the CRUD, used only if you want to keep the same file flow, see `Boilerplate code`

The table name will become the name of the struct and the endpoint's grouping name.  For each group/endpoint it will create:
- GET - read by primary key
- GET - search all records
- POST - create new record
- PATCH - update increment
- DELETE - delete by primary key

##### Boilerplate code
The general setup of the boilerplate code follows the convention of [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html).

Each layer is abstracted from each other by a Go interface.  Each layer can easily be tested without knowing *who* called it and *what* the downstream layer is.  See `manager_test.go` for some testing examples since most of the *business logic* would live in this file or layer.

Both a RESTful and GRPC server code is created for you.  See `cmd/rest` and `cmd/grpc`.

##### Env Vars
Inspect `config/config.go` for the names of the env vars that will need to be set.  Some have defaults that you may not want to set.  When adding the first entity, this file will change based on the `Storage Type` you set for the project.

### Admin Screen
This screen gives the ability to change the project values when setting up the project at the start.  You can change the `Storage Type` and `TagFormat`.

The 3rd option is adding a module, see `Login Module` for more details.

The last option is if your `Storage Type` is SQL, then you can use an ORM, GORM in particular.  I'm not a fan of ORM's so this option has no guarantee it working out of the box but I do create the CRUD data layer with GORM vs SQL statements when this option is turned on.

### Login Module
This module will insert an entity called `Login`, it will follow the same format as any other entity you may add.  This entity provides some extra functionality around creating password, resetting of passwords and signing in and providing a JWT token for authentication.

The reset password functionality also includes the ability to email a "reset password" email, this will just point to a web page to take the email address and reset token, add provide the user the ability to enter in a new password (and confirm password) then send those pieces of data to the API endpoint to reset the user's password.

A few tables are provided called `login` and `login_reset` and migration files are provided for these tables, see `tools/migration/README.md`.

`Admin Tool` is added in the `tools` directory, this should also be set up as a migration script to allow the `admin` to set their password.  Through the migration process, it guaranteed to run once.  See `tools/migration/README.md` for more information.

```
EnvVars added to `config/config.go` when this is enabled, see:
FORGE_GO_BASE_LOGIN_PWD_COST: [int] see https://pkg.go.dev/golang.org/x/crypto/bcrypt for more explanation, min: 4; max: 31
FORGE_GO_BASE_LOGIN_RESET_DURATION: [int] number of days the reset token is valid for
FORGE_GO_BASE_LOGIN_EXPIRES_AT_DURATION: [int] number of hours the JWT token is valid for after successful login
FORGE_GO_BASE_LOGIN_AUTH_SECRET: [string] you secret phrase that is used to encrypt the passwords within the login table
FORGE_GO_BASE_EMAIL_HOST: [string] the smtp host the reset emails will be sent to
FORGE_GO_BASE_EMAIL_PORT: [int] the port of the smtp host
FORGE_GO_BASE_EMAIL_PWD: [string] the stmp password
FORGE_GO_BASE_EMAIL_FROM: [string] the from email address for the reset email
FORGE_GO_BASE_EMAIL_RESET_URL: [string] the base url to the this service, the code will tack on the correct endpoint and the query stirng params for email and the reset token
FORGE_GO_BASE_ADMIN_EMAIL: [string] the email for the admin user, used when setting up the admin password for the first time

Note: those env vars that are indicated as 'int', should be quoted, env var need to be string but the values will turn into an 'int'.

The prefix of the above env vars should change to the name of new project.
```
Walkthrough :
- Setup the above env var(s)
- If using SQL, enable `FORGE_GO_BASE_MIGRATION_ENABLED` and use the provided migration scripts (see `tools/migration/README.md`) or run the table create statements manually for `login` and `login_reset`.  Scripts are provided for you in the `scripts/migrations` folder.
- The `tools/admin` is saved as a `.bin` file for to run under `migration`, if not, run the `tools/admin` manually by doing `go run tools/admin/main.go`.
- Build and run the `rest` server. If migration is enabled this will create the tables and run the admin binary, or if ran manually, a user will be created with the email from `FORGE_BASE_GO_ADMIN_EMAIL` and a reset token will be created in the `login_reset` table.
- An email will be sent (if set up correctly), that will use the `FORGE_GO_BASE_EMAIL_*` env vars to do this.
- If you have a web site that points to `FORGE_GO_BASE_EMAIL_RESET_URL`, then use that to enter in new password (and confirm password), else you can just call the `/login/reset/pwd` with the payload needed, the token is found in the email or get it from the `login_reset` table.  This will ensure the admin's password is encrypted correctly into the storage (DB, mongo, etc).
- Once the admin user is set and a password, than those credentials can be used to add, through the POST endpoint, new users.

Note: `forge` is designed to only allow this module to be added once, unless you remove it from `.forge`, bad things may happen.  You have been warned.

### Tools
There are 2 other tools provided to make `forge` a better overall tool.

`Migration` and `Regression`, please see their respected `README` files for set up and explanation.

### Features
There are a few built in features and have env vars to configure if needed.

##### Audit
By enabling `FORGE_GO_BASE_ENABLE_AUDITING` (default `false`) this will turn on an overall auditing whenever an entity is created, updated or deleted.  Set: `FORGE_GO_BASE_AUDIT_STORAGE` [`file` or `sql`], if `file` is chosen then set `FORGE_GO_BASE_AUDIT_FILE_PATH` to where you want to save this json formatted file.  If `sql` is chosen than a table structure of:
```
CREATE TABLE IF NOT EXISTS audit (
	id INT AUTO_INCREMENT, -- or SERIAL
	user_name VARCHAR(50) NULL,
	entity VARCHAR(50) NOT NULL,
	entity_id VARCHAR(50) NOT NULL,
	changes JSON NOT NULL, -- may not work with sqlite
	created_at DATETIME NOT NULL, -- or TIMESTAMP
	PRIMARY KEY(id)
);
-- change according to your sql engine and it's correct syntax
```
in order for the audit functionality to work correctly, see `internal/audit/audit.go` for more info.

##### Metrics
By default the rest server will create Prometheus metrics for each endpoint.  Set `FORGE_GO_BASE_ENABLE_METRICS` to `false` to turn off because it is set to `true` by default.

### Notes
Why I stuck with using singular nouns: https://stackoverflow.com/questions/6845772/should-i-use-singular-or-plural-name-convention-for-rest-resources

`forge` has some hooks that this code base inserts for you, so if you see a comment like `// --- replace ***** text - do not remove ---`, please don't remove.  `forge` uses these lines to "hook into" and manipulate the base code.

If `Blank Structure` is chosen, just know most of the boilerplate code will not be include but each file or layer is created and, at least, should compile.  It is up to you then, to add the logic for each layer.

I've don't have all sql types represented and they may very from engine to engine on the effectiveness.  So this code generation is gives you is an *as-is* end product, but it is *go* and you can change it how you need/want.

As of this writing, I have written but not tested multi-key functionality, so I can't say it works 100%.

There are a lot of TODO comments, regarding feed-back loops or reporting.  My thinking was, in the case of `login`, that there are times where the user may only need a success but a feed-back loop should go to the developer, admin or owner of the service to know when something has gone wrong.  I will leave that up to you, you maybe be happy using STDOUT and looking at the logs or another type of service would be good here.

All the helper scripts and internal code uses `bash`, unless you install `bash` for Windows, sorry, at this time it will probably not work.  I will put it on my TODO list, or come up with a solution and provide a PR.

##### TODO
- `bash` scripts to os appropriate scripting mechanism
- Roles for the admin/jwt still need to be implemented