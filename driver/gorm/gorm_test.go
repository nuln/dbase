package gorm_test

import (
	"testing"

	"gorm.io/driver/sqlite"

	"github.com/nuln/dbase/dbasetest"
	"github.com/nuln/dbase/driver/gorm"
)

func TestGormSQLite(t *testing.T) {
	db, err := gorm.New("sqlite", sqlite.Open(":memory:"))
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}
	defer func() { _ = db.Close() }()

	dbasetest.Suite(t, db)
}
