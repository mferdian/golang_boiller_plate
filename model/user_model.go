package model

import (
	"github.com/google/uuid"
	"github.com/mferdian/golang_boiller_plate/helpers"
	"gorm.io/gorm"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey" json:"user_id"`
	Name     string    `json:"user_name"`
	Email    string    `gorm:"unique; not null" json:"user_email"`
	Password string    `json:"user_password"`
	NoTelp   string    `json:"user_no_telp"`
	Address  string    `json:"user_address"`
	Role     string    `json:"role"`

	TimeStamp
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	var err error
	u.Password, err = helpers.HashPassword(u.Password)
	if err != nil {
		return err
	}

	return nil
}
