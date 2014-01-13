package mig

import (
	"database/sql"
	"fmt"
)

type migration struct {
	version int
	script  string
}

var migrations = make([]migration, 0, 100)

// Define a migration
func Define(script string) {
	DefineVersion(len(migrations)+1, script)
}

// Define a migration with a specific version
func DefineVersion(version int, script string) {
	migrations = append(migrations, migration{version, script})
}

// Execute the migrations
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

// Create the db_version table if it doesn't exist
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

// Test whether a migration should be executed
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

// Execute a migration. This assumes the migration hasn't already been run.
func execute(script string, db *sql.DB) error {
	_, err := db.Exec(fmt.Sprintf(`DO $$
		BEGIN
		%s
		END;
		$$;`, script))

	return err
}

// Inserts a version into the db_version table. Assumes it doesn't already
// exist in there.
func insertVersion(version int, db *sql.DB) error {
	sql := fmt.Sprintf("insert into db_version (version) values (%d);", version)
	_, err := db.Exec(sql)
	return err
}
