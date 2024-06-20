package models

import (
	pq "github.com/lib/pq"
	"gorm.io/gorm"
)

type BlogModel struct {
	gorm.Model

	Title             string
	ThumbnailUrl      string
	Discription       string
	PubTime           string
	BlogUrl           string
	CompanyId         int
	TagsId            pq.Int64Array `gorm:"type:integer[]"`
	ThumbnailBlurhash string        `gorm:"size:255"`
}
