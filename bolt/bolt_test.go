package bolt_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nuln/dbase/bolt"
	"github.com/nuln/dbase/dbasetest"
)

func TestBolt(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	db, err := bolt.New(dbPath)
	if err != nil {
		t.Fatalf("failed to open bolt: %v", err)
	}
	defer func() {
		_ = db.Close()
		_ = os.Remove(dbPath)
	}()

	dbasetest.Suite(t, db)
}
