package models

import (
	pq "github.com/lib/pq"
	"gorm.io/gorm"
)

type BlogCategoryModel struct {
	gorm.Model
	CategoryName string
	TagsId       pq.Int64Array `gorm:"type:integer[]"`
}
