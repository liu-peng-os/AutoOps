package db

import "testing"

type tableNamer interface {
	TableName() string
}

func collectTableNames(t *testing.T, models []interface{}) map[string]struct{} {
	t.Helper()

	names := make(map[string]struct{}, len(models))
	for _, model := range models {
		namer, ok := model.(tableNamer)
		if !ok {
			t.Fatalf("model %T does not implement TableName", model)
		}
		names[namer.TableName()] = struct{}{}
	}
	return names
}

func TestSQLManagedTablesExcludedFromAutoMigrateFallback(t *testing.T) {
	sqlManaged := collectTableNames(t, sqlManagedModels)
	legacyManaged := collectTableNames(t, legacyAutoMigrateModels)

	for table := range sqlManaged {
		if _, exists := legacyManaged[table]; exists {
			t.Fatalf("table %s is managed by both SQL migrations and AutoMigrate fallback", table)
		}
	}
}

func TestSQLManagedTableNamesMatchManagedModels(t *testing.T) {
	expected := collectTableNames(t, sqlManagedModels)
	actual := SQLManagedTableNames()

	if len(expected) != len(actual) {
		t.Fatalf("SQLManagedTableNames length = %d, want %d", len(actual), len(expected))
	}

	for _, table := range actual {
		if _, exists := expected[table]; !exists {
			t.Fatalf("unexpected SQL managed table: %s", table)
		}
	}
}

func TestParseLegacyMigrationVersion(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int64
		wantErr bool
	}{
		{name: "sequential sql", input: "0004_delivery_flow.up.sql", want: 4},
		{name: "goose sql", input: "0004_delivery_flow.sql", want: 4},
		{name: "invalid", input: "delivery_flow.sql", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseLegacyMigrationVersion(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error for %s", tt.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseLegacyMigrationVersion(%q) error = %v", tt.input, err)
			}
			if got != tt.want {
				t.Fatalf("parseLegacyMigrationVersion(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}
