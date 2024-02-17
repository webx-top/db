//go:build !sqlitego

package sqlite

import (
	_ "github.com/mattn/go-sqlite3" // SQLite3 driver.
)
