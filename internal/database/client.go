package database

import (
	"context"
	"database/sql"

	"github.com/LalatinaHub/common/database"
)

// GetDB returns a singleton database connection pool.
// Forwarded to github.com/LalatinaHub/common/database.GetDB.
func GetDB() (*sql.DB, error) {
	return database.GetDB()
}

// Ping verifies that the database connection is alive.
// Forwarded to github.com/LalatinaHub/common/database.Ping.
func Ping(ctx context.Context) error {
	return database.Ping(ctx)
}

// InitIndexes ensures critical performance indexes exist on tables.
// Forwarded to github.com/LalatinaHub/common/database.InitIndexes.
func InitIndexes(ctx context.Context) error {
	return database.InitIndexes(ctx)
}

// Close closes the database connection pool.
// Forwarded to github.com/LalatinaHub/common/database.Close.
func Close() error {
	return database.Close()
}

// SetDB explicitly overrides the singleton DB instance.
// Forwarded to github.com/LalatinaHub/common/database.SetDB.
func SetDB(db *sql.DB) {
	database.SetDB(db)
}

// ResetDBInstance resets the singleton state.
// Forwarded to github.com/LalatinaHub/common/database.ResetDBInstance.
func ResetDBInstance() {
	database.ResetDBInstance()
}
