package bolt_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nuln/dbase/dbasetest"
	"github.com/nuln/dbase/driver/bolt"
)

func TestBolt(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	db, err := bolt.New(dbPath)
	if err != nil {
		t.Fatalf("failed to open bolt: %v", err)
	}
	defer func() {
		db.Close()
		os.Remove(dbPath)
	}()

	dbasetest.Suite(t, db)
}
