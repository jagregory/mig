# mig

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
