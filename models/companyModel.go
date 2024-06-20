package models

import (
	"gorm.io/gorm"
)

type CompanyModel struct {
	gorm.Model
	CompanyName    string
	CompanyLogoURL string
}
