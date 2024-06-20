package models

import "gorm.io/gorm"

type TagModel struct {
	gorm.Model

	TagName string
}
