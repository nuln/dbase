package gorm_test

import (
	"testing"

	"github.com/nuln/dbase/dbasetest"
	"github.com/nuln/dbase/driver/gorm"
	"gorm.io/driver/sqlite"
)

func TestGormSQLite(t *testing.T) {
	db, err := gorm.New("sqlite", sqlite.Open(":memory:"))
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}
	defer db.Close()

	dbasetest.Suite(t, db)
}
