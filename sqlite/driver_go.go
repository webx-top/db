//go:build sqlitego

package sqlite

import (
	_ "github.com/glebarez/go-sqlite"        // SQLite3 driver. sqlite://
	_ "github.com/glebarez/go-sqlite/compat" // SQLite3 driver. sqlite3://
)
