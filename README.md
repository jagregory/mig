# mig

[![GoDoc](https://godoc.org/github.com/jagregory/mig?status.png)](https://godoc.org/github.com/jagregory/mig)

Go SQL migration package

## Usage

```go
import "github.com/jagregory/mig"

mig.Define(`
	create table foo ( id integer );
`)

mig.DefineVersion(2, `
	create table bar ( id integer );
`)

mig.Migrate(db)
```

The `Define` and `DefineVersion` functions are used to define a migration
and the SQL that will be executed for it. The `Migrate` function actually
executes the migrations. An `error` will be returned from `Migrate` if
anything bad happens.

## How it works

Mig will create a `db_version` table in your database with a `version int`
column. Each migration will be executed in the order they're defined in
unless there's a corresponding row in the `db_version` table. After each
successful migration a row will be inserted in `db_version` for that
migration.
