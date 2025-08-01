package migrations

import (
	"github.com/mferdian/golang_boiller_plate/model"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&model.User{},
	); err != nil {
		return err
	}

	return nil
}
