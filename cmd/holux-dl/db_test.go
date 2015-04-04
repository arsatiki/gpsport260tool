package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

func TestDBInit(t *testing.T) {
	var name string

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	err = initialize(db)
	if err != nil {
		t.Fatal(err)
	}

	for _, table := range TABLES {
		err := db.QueryRow(NEEDS_INIT, table).Scan(&name)
		if err != nil || name != table {
			t.Fatalf("expected table %s, got %s (err: %v)",
				table, name, err)
		}
	}
}
