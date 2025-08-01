package migrations

import (
	"github.com/mferdian/golang_boiller_plate/model"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {
	err := SeedFromJSON[model.User](db, "./migrations/json/users.json", model.User{}, "Email")
	if err != nil {
		return err
	}

	return nil
}