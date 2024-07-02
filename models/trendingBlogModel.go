package models

import "gorm.io/gorm"

type TrendingBlogModel struct {
	gorm.Model

	BlogId uint
}
