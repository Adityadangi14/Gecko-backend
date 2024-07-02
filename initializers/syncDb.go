package initializers

import "gecko_backend/models"

func SyncDB() {
	DB.AutoMigrate(&models.BlogModel{})
	DB.AutoMigrate(&models.CompanyModel{})
	DB.AutoMigrate(&models.TagModel{})
	DB.AutoMigrate(&models.UserModel{})
	DB.AutoMigrate(&models.UserTagsModel{})
	DB.AutoMigrate(&models.TrendingBlogModel{})

}
