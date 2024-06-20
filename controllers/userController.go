package controllers

import (
	"gecko_backend/initializers"
	"gecko_backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    user,
	})
}

func EditUserTagSelectionStatus(c *gin.Context) {

	var body struct {
		AreTagsChoosen bool
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Invalid request body format",
		})
		return
	}

	var usr models.UserModel

	user, ok := c.Get("user")

	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"success": true,
			"message": "user not found",
		})

		return
	}

	if user, ok := user.(models.UserModel); ok {
		usr = user
	} else {

		c.JSON(http.StatusNotFound, gin.H{
			"success": true,
			"message": "Error in managing user",
		})
		return
	}

	usr.AreTagsChoosen = body.AreTagsChoosen

	er := initializers.DB.Save(&usr).Error

	if er != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Unable to update user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    usr,
	})

}

func CreateUserTags(c *gin.Context) {
	var body struct {
		TagsId []int64
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Failed to read body",
		})
		return
	}

	var usr models.UserModel

	user, ok := c.Get("user")

	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"success": true,
			"message": "user not found",
		})

		return
	}

	if user, ok := user.(models.UserModel); ok {
		usr = user
	} else {

		c.JSON(http.StatusNotFound, gin.H{
			"success": true,
			"message": "Error in managing user",
		})
		return
	}

	usr.AreTagsChoosen = true

	er := initializers.DB.Save(&usr).Error

	if er != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": true,
			"message": "Error in managing user",
		})
		return
	}

	if user, ok := user.(models.UserModel); ok {
		usr = user
	} else {

		c.JSON(http.StatusNotFound, gin.H{
			"success": true,
			"message": "Error in managing user",
		})
		return
	}

	userTags := models.UserTagsModel{
		UserId: usr.ID,
		TagsId: body.TagsId,
	}

	result := initializers.DB.Create(&userTags)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Failed to create tags for user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,

		"message": " successfully creted tags",
	})
}

// func getFieldNames(obj interface{}) []string {
// 	typ := reflect.TypeOf(obj)
// 	if typ.Kind() != reflect.Struct {
// 		return nil
// 	}

// 	fieldNames := make([]string, 0, typ.NumField())
// 	for i := 0; i < typ.NumField(); i++ {
// 		field := typ.Field(i)
// 		fieldNames = append(fieldNames, field.Name)
// 	}
// 	return fieldNames
// }
