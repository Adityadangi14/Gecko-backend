package models

import (
	pq "github.com/lib/pq"
	"gorm.io/gorm"
)

type UserTagsModel struct {
	gorm.Model

	UserId string

	TagsId pq.Int64Array `gorm:"type:integer[]"`
}
