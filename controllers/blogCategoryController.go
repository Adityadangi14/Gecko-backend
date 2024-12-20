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
	"github.com/go-redis/redis"
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

	val,err:=	initializers.RC.Get(rediskeys.BlogCategories).Result()


	if err == redis.Nil {
		fmt.Println("Key does not exist.")
	} else if err != nil {
		fmt.Println("Error fetching value:", err)
	} else {

	var result []models.BlogCategoryModel
		
err:=json.Unmarshal([]byte(val),&result)

	if err !=nil{
		fmt.Println(err)
		return
	}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    result,
		})
		return
	}
	blogCategories, err := getAllBlogCategories()
	
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "unable to fetch categories",
		})
		return
	}

	json,er:= json.Marshal(blogCategories)

	if(er !=nil){
		fmt.Println("failed to set cache")
	}

	initializers.RC.Set(rediskeys.BlogCategories,json,0);

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


func GetBlogsByCat(c *gin.Context){

	var body struct{
		CatName string
	}
	fmt.Println(c.Request.Body)
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Failed to read body",
		})
		return
	}

	pageNo,err := strconv.Atoi(c.Query("pageNo"))

	var total int64
	initializers.DB.Model(&models.BlogModel{}).Count(&total)

	val ,_ := services.GetCachedData(rediskeys.BlogCategories+"/"+c.Query("pageNo")+"/"+body.CatName)

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

	var cachedBlogs []CachedBlogs

	if val != ""{

	err :=	json.Unmarshal([]byte(val),&cachedBlogs)

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
	
	var catModel models.BlogCategoryModel

	rows :=initializers.DB.Where("category_name  = ?", body.CatName). Find(&catModel).RowsAffected

	if rows==0{
		if rows == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Failed to get blogs",
			})
			return
		}
	}
	fmt.Println(catModel.TagsId)
	var blogs []models.BlogModel

	err = initializers.DB .Model(&models.BlogModel{}).Where("tags_id && ?", pq.Int64Array(catModel.TagsId)).Limit(5).Offset((pageNo-1)*5).Find(&blogs).Count(&total). Error

	if err !=nil{

		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "unable to fetch blogs",
		})

		return
	}

	var blogArray []map[string]interface{}	



	for _,blog := range blogs{
		var company models.CompanyModel
		err = initializers.DB.Where("id = ?", blog.CompanyId).Find(&company).Error

		if(err != nil){
			return
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

	blogsStr,err:=	json.Marshal(blogArray)

	if err !=nil{
		fmt.Println(err)
	}

	initializers.RC.Set(rediskeys.BlogCategories+"/"+c.Query("pageNo")+"/"+body.CatName,blogsStr,time.Hour);

c.JSON(http.StatusOK, gin.H{
	"success": true,
	"data":    blogArray,
	"isNextPageExist": total> (5*int64(pageNo)),
})
}