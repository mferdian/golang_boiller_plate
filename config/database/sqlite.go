package database

import (
	"log"
	"os"
	"testing"

	"github.com/mferdian/golang_boiller_plate/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB membuat SQLite in-memory DB
func SetupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}

	// Auto migrate schema
	err = db.AutoMigrate(
		&model.User{},
	)
	if err != nil {
		t.Fatalf("failed to migrate db: %v", err)
	}

	return db
}

func TestMain(m *testing.M) {
	code := m.Run()
	log.SetOutput(os.Stdout)
	os.Exit(code)
}
