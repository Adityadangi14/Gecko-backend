package controllers

import (
	"gecko_backend/initializers"
	"gecko_backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func CreateBlogCategory(c *gin.Context) {

	var body struct {
		CategoryName string
		TagIds       pq.Int64Array `gorm:"type:integer[]"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Failed to read body",
		})
		return
	}

	blogCategory := models.BlogCategoryModel{
		CategoryName: body.CategoryName,
		TagsId:       body.TagIds,
	}

	rows := initializers.DB.Create(&blogCategory).RowsAffected

	if rows == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Failed to create blog",
		})
		return
	}

	blogCategories, err := getAllBlogCategories()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    blogCategories,
	})

}

func GetBlogCategories(c *gin.Context) {
	blogCategories, err := getAllBlogCategories()

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "unable to fetch categories",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    blogCategories,
	})

}

func getAllBlogCategories() ([]models.BlogCategoryModel, error) {

	var blogCategory []models.BlogCategoryModel

	err := initializers.DB.Find(&blogCategory).Error

	if err != nil {
		return nil, err
	}

	if blogCategory == nil {
		return make([]models.BlogCategoryModel, 0), nil
	}

	return blogCategory, nil
}
