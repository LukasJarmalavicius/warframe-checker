package services

import (
	"path/filepath"
	"testing"
	"warframe-checker/internal/storage"
)

func newTestDB(t *testing.T) *storage.DB {
	t.Helper()
	db, err := storage.New(filepath.Join(t.TempDir(), "test.sqlite"))
	if err != nil {
		t.Fatal("failed to create test database:", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}
