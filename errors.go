package mig

import "fmt"

// Returned when an error occurs running one of the migrations. If the
// version is 0, the initial db_version table creation failed.
type MigrationError struct {
	Cause   error
	Version int
}

func (e MigrationError) Error() string {
	return fmt.Sprintf("Error migrating #%d: %s", e.Version, e.Cause.Error())
}
