# FORGE_GO_BASE

Description of the app here

### Usage
This application provides 2 servers:

**rest** - RESTful service
```
cd cmd/rest
go build && ./rest
or
go run main.go
(optional) -restPort <port number> e.g.: ./rest -restPort 10001
```

**grpc** - GPRC service
```
cd cmd/grpc
go build && ./grpc
or
go run main.go
```

### Env Vars
See config/config.go for the complete list of `env vars`, most are self-explanatory. Some additional explanation will follow in this document.

### Features
**Prometheus**: collect metrics in memory, set to `true` by default.  Set this to `false` to disable it:
`FORGE_GO_BASE_ENABLE_METRICS`: [bool] true/false

**SQL Migration**: tool to manager sql migration files, the feature is disabled by default.  See `tools/migration/README.md` for more info.
`FORGE_GO_BASE_MIGRATION_ENABLED`: [bool] true/false
`FORGE_GO_BASE_MIGRATION_PATH`: [string] path to your migration scripts
`FORGE_GO_BASE_MIGRATION_SKIP_INIT`: [bool] true/false (optional) set to true if you don't want the feature to make a DB with the projects name and create the `migration` table

**Audit**: will save changes per row to storage device `file|sql`, set to `false` by default
`FORGE_GO_BASE_ENABLE_AUDITING`: [bool] true/false to enable/disable
`FORGE_GO_BASE_AUDIT_STORAGE`: [string] which storage type to save row audit data `file | sql`
`FORGE_GO_BASE_AUDIT_FILE_PATH`: [string] if `file` is the storage type, path to read/save the audit file

**Add your documentation here**