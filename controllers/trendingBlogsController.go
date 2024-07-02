package controllers

import (
	"gecko_backend/initializers"
	"gecko_backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTrendingBlogs(c *gin.Context) {
	blogs, err := getTrendingBlogs()

	if err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"success": false,

			"message": "Unable to get blogs",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": false,

		"data": blogs,
	})

}

func AddTrendingBlog(c *gin.Context) {
	var body struct {
		BlogId uint
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Failed to read body",
		})
		return
	}

	var trendingBlog models.TrendingBlogModel

	trendingBlog.BlogId = body.BlogId

	err := initializers.DB.Save(&trendingBlog).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,

			"message": "Unable to save blogs",
		})
	}
	blogs, err := getTrendingBlogs()

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,

			"message": "Something went wrong",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,

		"data": blogs,
	})

}

func DeleteTrendingBlog(c *gin.Context) {
	var body struct {
		BlogId uint
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Failed to read body",
		})
		return
	}

	var trendingBlog models.TrendingBlogModel

	err := initializers.DB.Where("id = ?", body.BlogId).Delete(&trendingBlog).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,

			"message": "Unbale to delete blog",
		})
	}

	blogs, err := getTrendingBlogs()

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,

			"message": "Unbale to get blog",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,

		"data": blogs,
	})

}

func getTrendingBlogs() ([]models.BlogModel, error) {

	var trendingBlogs []models.TrendingBlogModel

	var blogArray []models.BlogModel

	err := initializers.DB.Find(&trendingBlogs).Error

	if err != nil {
		return nil, err
	}

	for _, blogId := range trendingBlogs {

		var blog models.BlogModel

		err := initializers.DB.Where("id = ?", blogId.BlogId).Find(&blog).Error

		if err != nil {
			return nil, err
		}

		blogArray = append(blogArray, blog)
	}

	return blogArray, nil
}
