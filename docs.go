// Package mig is a simple SQL migration library. You use it to embed
// snippets of SQL into your application which will be run on startup (or
// whenever you tell them to). Mig is primarily used for migrating schemas
// between versions of applications.
//
// Mig is small and deliberately simple. There are no rollbacks in Mig,
// and it's expected that the application running Mig will have schema
// altering permissions on the target database.
//
// If you want to restrict your application's rights to a database, it's
// best you run Mig from a different process or under a different user.
package mig
