# forge-go
Clone, build and deploy an API written in Golang, quick and easy.  A base framework to add entities to create standard CRUD endpoints based on a SQL create table statement.  From there, customize as needed.

Why I stuck with using singular nouns: https://stackoverflow.com/questions/6845772/should-i-use-singular-or-plural-name-convention-for-rest-resources

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
`<path>` is full path from `$GOPATH/src` e.g.: `github.com/blackflagsoftware/new_project`, if the folder doesn't exists, the folder will be created and the project will be copied into this folder.  If the folder exists, bad things happen, kidding, and message will show indicating it exists.

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
`Paste as input` as suggested, this will prompt to paste in the SQL create table statement.  It prompts for multiple statements if needed.
`Prompt as input`, this will prompt for entity name and then add as many fields as you need.  It will prompt for column type, is it null, etc.  When finished it will save the SQL to `./prompt_schema` for your reference.
`Blank Structure`, this will create all the files as if it was creating the CRUD, used only if you want to keep the same file flow, see `Boilerplate code`

The table name will become the name of the struct and the endpoint's grouping name.  For each group/endpoint it will create:
- GET - read by primary key
- GET - search all records
- POST - create new record
- PATCH - update increment
- DELETE - delete by primary key

##### Boilerplate code
The general setup of boilerplate code follows the convention of [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html).

Each layer is abstracted from each other by a Go interface.  Each layer can easily be tested without knowing *who* called it and *what* the downstream layer is.  See `manager_test.go` for some testing examples since most of the *business logic* would live in this file (layer).

Both a RESTful and GRPC server code is created for you.  See `cmd/rest` and `cmd/grpc`.

### Admin Screen
This 

### Notes
`forge` has some hooks that this code base inserts for you, so if you see a comment like `// --- replace ***** text - do not remove ---`, please don't remove.  `forge` uses these lines to "hook into" and manipulate the base code.

If `Blank Structure` is chosen, just know most of the boilerplate code will not be include but each file (layer) is created and, at least, should compile.  It is up to you then to add the logic for each layer.

I've don't have all sql types represented and they may very from engine to engine on the effectiveness.  So this code generation is gives you an *as-is* end product, but it is *go* and you can change it how you need/want.

As of this writing, I have written but not tested multi-key functionality, so I can't say it works 100%.
