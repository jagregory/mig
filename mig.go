package mig

import (
	"database/sql"
	"fmt"
)

// migration is a combination of a SQL snippet and a version number which
// represents it.
type migration struct {
	version int
	script  string
}

var migrations = make([]migration, 0, 100)

// Define will define and register a migration for running the next time
// Mig is run. A version is automatically chosen based on the sequence of
// previous calls to Define (aka version=len(migrations)+1).
//
//     mig.Define(`
//       create table foo ( id integer );
//     `)
func Define(script string) {
	DefineVersion(len(migrations)+1, script)
}

// DefineVersion will define and register a migration for running the next
// time Mig is run with a specific version number.
//
//     mig.DefineVersion(2, `
//       alter table foo add column bar varchar(10);
//     `)
func DefineVersion(version int, script string) {
	migrations = append(migrations, migration{version, script})
}

// Migrate will execute the defined migrations against the sql.DB. For
// each migration Mig will look in a db_version table to see if it has
// already been run. If there's a row in db_version for a migration it
// wont be run again, otherwise the migration will be run and a row stored
// in db_version.
//
// If this is the first time Mig has been run, Migrate will first create
// a db_version table and store a zero row.
func Migrate(db *sql.DB) error {
	fmt.Println("Running migrations")

	if err := createVersionTable(db); err != nil {
		return MigrationError{err, 0}
	}

	for _, m := range migrations {
		ok, err := shouldRun(m.version, db)
		if err != nil {
			return MigrationError{err, m.version}
		}

		if ok {
			fmt.Printf("Executing #%d...", m.version)

			if err := execute(m.script, db); err != nil {
				return MigrationError{err, m.version}
			}

			if err := insertVersion(m.version, db); err != nil {
				return MigrationError{err, m.version}
			}

			fmt.Println(" done.")
		} else {
			fmt.Printf("Skipping #%d, already run\n", m.version)
		}
	}

	fmt.Println("Migration complete")
	return nil
}

// createVersionTable creates the db_version table if it doesn't exist
func createVersionTable(db *sql.DB) error {
	return execute(`
		create table if not exists db_version (
			version     integer not null unique,
			updated_at  timestamp not null default(current_timestamp)
		);

		if not exists(select 1 from db_version where version = 0) then
			insert into db_version (version) values (0);
		end if;
	`, db)
}

// shouldRun tests whether a migration should be executed
func shouldRun(version int, db *sql.DB) (bool, error) {
	sql := fmt.Sprintf("select 1 from db_version where version = %d;", version)

	result, err := db.Exec(sql)
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rows == 0, nil
}

// execute runs a migration. This assumes the migration hasn't already
// been run.
func execute(script string, db *sql.DB) error {
	_, err := db.Exec(fmt.Sprintf(`DO $$
		BEGIN
		%s
		END;
		$$;`, script))

	return err
}

// insertVersion inserts a version into the db_version table. Assumes it
// doesn't already exist in there.
func insertVersion(version int, db *sql.DB) error {
	sql := fmt.Sprintf("insert into db_version (version) values (%d);", version)
	_, err := db.Exec(sql)
	return err
}
