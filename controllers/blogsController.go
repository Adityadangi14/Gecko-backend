package controllers

import (
	"encoding/json"
	"fmt"
	rediskeys "gecko_backend/constants/redisKeys"
	"gecko_backend/initializers"
	"gecko_backend/models"
	"gecko_backend/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func CreateBlog(c *gin.Context) {
	var body struct {
		Title             string
		ThumbnailUrl      string
		Discription       string
		PubTime           string
		BlogUrl           string
		CompanyId         int
		TagsId            []int64
		ThumbnailBlurhash string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Failed to read body",
		})
		return
	}

	blog := models.BlogModel{

		Title: body.Title,

		ThumbnailUrl: body.ThumbnailUrl,

		Discription: body.Discription,

		PubTime: body.PubTime,

		BlogUrl: body.BlogUrl,

		TagsId: body.TagsId,

		CompanyId: body.CompanyId,

		ThumbnailBlurhash: body.ThumbnailBlurhash,
	}

	result := initializers.DB.Create(&blog)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Failed to create blog",
		})
		return
	}

	availableblogs, err := getBlogsPvt(1)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get blogs",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    availableblogs,
		"message": "Blog created with blog id " + strconv.FormatUint(uint64(blog.ID), 10),
	})
}

func GetBlogs(c *gin.Context) {

	type TagsInApi struct {
		TagName string
		ID      uint
	}

	type CachedBlogs struct{

	Title             string
	ThumbnailUrl      string
	Description       string
	PubTime           string
	BlogUrl           string
	CompanyId         int
	Company 		  models.CompanyModel
	Tags            []TagsInApi
	TThumbnailBlurhash string        `gorm:"size:255"`

	}

	pageNo,err := strconv.Atoi(c.Query("pageNo"))

	var total int64
	initializers.DB.Model(&models.BlogModel{}).Count(&total)

	if(err!=nil){
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to parse page",
		})
		return
	}

	val ,_ := services.GetCachedData(rediskeys.AllBlogs+"/"+c.Query("pageNo"))

	var cachedBlogs []CachedBlogs

	if val != ""{
		err= json.Unmarshal([]byte(val),&cachedBlogs)

		if err!=nil{
			fmt.Println(err)
		}else{
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"data":    cachedBlogs,
				"isNextPageExist": total> (5*int64(pageNo)),
				"isCached":true,
			})
			return
		}
	}

	availableblogs, err := getBlogsPvt(pageNo)

	blogsStr,err:=	json.Marshal(availableblogs)

	if err !=nil{
		fmt.Println(err)
	}

	initializers.RC.Set(rediskeys.AllBlogs+"/"+c.Query("pageNo"),blogsStr,time.Hour);

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get blogs",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    availableblogs,
		"isNextPageExist": total> (5*int64(pageNo)),
	})
}

func DeleteBlog(c *gin.Context) {
	var body struct {
		BlogId uint
	}

	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	blog := models.BlogModel{}

	rows := initializers.DB.Delete(&blog, body.BlogId).RowsAffected

	blogs, err := getBlogsPvt(1)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err,
		})
		return
	}

	if rows == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Blog id " + strconv.FormatUint(uint64(body.BlogId), 10) + " not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    blogs,
		"message": "Blog deleted with blog id " + strconv.FormatUint(uint64(body.BlogId), 10),
	})
}

func GetBlogsByTags(c *gin.Context) {
	var body struct {
		TagsId pq.Int64Array `gorm:"type:integer[]"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Failed to read body",
		})
		return
	}

	var blogTagIds pq.Int64Array = body.TagsId

	var blogs []models.BlogModel
	initializers.DB.Model(blogs).Where("tags_id && ?", pq.Int64Array(blogTagIds)).Find(&blogs)

	fmt.Println(blogTagIds)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"blogs":   blogs,
	})
}

func getBlogsPvt(pageNo int) ([]map[string]interface{}, error) {

	var blogs []models.BlogModel



	result := initializers.DB.Limit(5).Offset((pageNo-1)*5).Find(&blogs)

	if result.Error != nil {

		return nil, fmt.Errorf("unable to find blogs: %w", result.Error)
	}

	blogJson := make([]map[string]interface{}, 0, len(blogs))

	for _, i := range blogs {

		// getting tags
		type TagsInApi struct {
			TagName string
			ID      uint
		}

		var tagsApiItem TagsInApi

		var tags []TagsInApi

		for _, j := range i.TagsId {

			result := initializers.DB.Model(&models.TagModel{}).Where("id = ?", j).First(&tagsApiItem)

			if result.Error != nil {
				return nil, fmt.Errorf("unable to find tags in comapnaies")
			}

			tags = append(tags, tagsApiItem)
		}

		// getting company data
		type CompanyDetail struct {
			CompanyName    string
			CompanyLogoURL string
			ID             uint
		}

		var companyDetail CompanyDetail

		result := initializers.DB.Model(&models.CompanyModel{}).Where("id = ?", i.CompanyId).First(&companyDetail)

		if result.Error != nil {
			return nil, fmt.Errorf("unable to find comapnaies")
		}
		// putting all the things to return
		blog := map[string]interface{}{
			"ID":                 i.ID,
			"Title":              i.Title,
			"ThumbnailUrl":       i.ThumbnailUrl,
			"Description":        i.Discription,
			"BlogUrl":            i.BlogUrl,
			"Tags":               tags,
			"PubTime":            i.PubTime,
			"Company":            companyDetail,
			"TThumbnailBlurhash": i.ThumbnailBlurhash,
		}

		blogJson = append(blogJson, blog)

	}

	return blogJson, nil
}
