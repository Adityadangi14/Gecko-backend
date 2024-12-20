package controllers

import (
	"encoding/json"
	"fmt"
	rediskeys "gecko_backend/constants/redisKeys"
	"gecko_backend/initializers"
	"gecko_backend/models"
	"gecko_backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func GetTrendingBlogs(c *gin.Context) {

data,_:=	services.GetCachedData(rediskeys.TrendingBlogs)

if data!= ""{

	type  TrendingBlogModel struct{
		Title             string
		ThumbnailUrl      string
		Description       string
		PubTime           string
		BlogUrl           string
		CompanyId         int
		Company   models.CompanyModel
		TagsId            pq.Int64Array 
		ThumbnailBlurhas string     
	}

	var trendingBlogs []TrendingBlogModel

	err:=json.Unmarshal([]byte(data),&trendingBlogs)
	if err!=nil{
		fmt.Println(err)
	}
	fmt.Println("cache hit")
	c.JSON(http.StatusOK, gin.H{
		"success": true,

		"data": trendingBlogs,
	})

	return
}

	blogs, err := getTrendingBlogs()

	if err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"success": false,

			"message": "Unable to get blogs",
		})

		return
	}

	json,er:= json.Marshal(blogs)

	if(er !=nil){
		fmt.Println("failed to set cache")
	}

	initializers.RC.Set(rediskeys.TrendingBlogs,json,0);


	c.JSON(http.StatusOK, gin.H{
		"success": true,

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

	rows := initializers.DB.Where(models.TrendingBlogModel{BlogId: body.BlogId}).Delete(&trendingBlog).RowsAffected

	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,

			"message": "Unbale to delete blog",
		})

		return
	}

	blogs, err := getTrendingBlogs()

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,

			"message": "Unbale to get blog",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,

		"data": blogs,
	})

}

func getTrendingBlogs() ([]map[string]interface{}, error) {

	var trendingBlogs []models.TrendingBlogModel

	var blogArray []map[string]interface{}
	

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

		var company models.CompanyModel

		err = initializers.DB.Where("id = ?", blog.CompanyId).Find(&company).Error

		if(err!=nil){
			return nil, err
		}
		blogArray = append(blogArray, map[string]interface{}{
				"Title":blog.Title,
				"ThumbnailUrl":blog.ThumbnailUrl,
				"Description":blog.Discription,
				"PubTime":blog.PubTime,
				"BlogUrl":blog.BlogUrl,
				"CompanyId":blog.CompanyId,
				"Company":company,
				"TagsId":blog.TagsId,
				"ThumbnailBlurhas":blog.ThumbnailBlurhash,
		})
	}

	if blogArray == nil {
		return nil, nil
	}

	return blogArray, nil
}
