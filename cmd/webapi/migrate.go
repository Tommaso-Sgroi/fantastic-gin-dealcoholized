package main

import (
	"database/sql"
	"os"

	"github.com/pressly/goose/v3"
)

func Migrate(db *sql.DB, dialect string) error {
	embedMigrations := os.DirFS("db" + string(os.PathSeparator) + "migrations")
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect(dialect); err != nil {
		return err
	}
	if err := goose.Up(db, "."); err != nil {
		return err
	}
	return nil
}
