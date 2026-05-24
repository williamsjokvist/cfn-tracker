package sql

import (
	"testing"
)

func TestInMemorySchema(t *testing.T) {
	storage, err := NewStorage(true)
	if err != nil {
		t.Fatalf("NewStorage(true) failed: %v", err)
	}
	defer storage.DB().Close()

	err = storage.DB().Ping()
	if err != nil {
		t.Fatalf("database ping failed: %v", err)
	}

	var tables []string
	err = storage.DB().Select(&tables, "SELECT name FROM sqlite_master WHERE type='table' ORDER BY name")
	if err != nil {
		t.Fatalf("failed to query tables: %v", err)
	}

	expectedTables := map[string]bool{
		"matches":    false,
		"migrations": false,
		"sessions":   false,
		"users":      false,
	}

	for _, table := range tables {
		if _, exists := expectedTables[table]; exists {
			expectedTables[table] = true
		}
	}

	for table, found := range expectedTables {
		if !found {
			t.Errorf("expected table %q not found in schema", table)
		}
	}
}
