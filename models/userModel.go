package models

import (
	"strings"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model

	ID             string `gorm:"primary_ key"`
	UserName       string
	Email          string
	PID            string `gorm:"<-:create"`
	ProfilePicture string
	AuthType       string `gorm:"<-:create"`
	AreTagsChoosen bool
}

func (u *UserModel) BeforeCreate(tx *gorm.DB) (err error) {
	// Generate a new UUID
	uuid, err := uuid.NewV4()

	if err != nil {
		return err
	}

	// Convert to string without hyphens
	u.ID = strings.ReplaceAll(uuid.String(), "-", "")
	return
}
