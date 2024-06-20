package controllers

import (
	"fmt"
	"gecko_backend/initializers"
	"gecko_backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateTag(c *gin.Context) {
	var body struct {
		Tag string
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Invalid request body format",
		})
		return
	}

	tag := models.TagModel{TagName: body.Tag}

	result := initializers.DB.Create(&tag)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create company",
		})
		return
	}

	tags, err := getTags()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "unable to fetch Tags",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Tag has been added with id " + strconv.FormatUint(uint64(tag.ID), 10),
		"tags":    tags,
	})
}

func DeleteTag(c *gin.Context) {
	var body struct {
		TagId string
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Invalid request body format",
		})
		return
	}

	tag := models.TagModel{}

	rows := initializers.DB.Delete(&tag, body.TagId).RowsAffected

	if rows == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Tag id " + body.TagId + " not found",
		})
		return
	}

	result, err := getTags()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "unable to fetch Tags",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Tag has been deleted with id " + body.TagId,
		"tags":    result,
	})
}

func GetAllTags(c *gin.Context) {
	result, err := getTags()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "unable to fetch Tags",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"tags":    result,
	})

}

func EditTag(c *gin.Context) {
	var body struct {
		TagId   string
		TagName string
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Invalid request body format",
		})
		return
	}
	var tag models.TagModel
	result := initializers.DB.Where("id = ?", body.TagId).First(&tag)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Unable to fetch tag",
		})
		return
	}
	fmt.Print(result.Statement.Vars...)
	tag.TagName = body.TagName

	er := initializers.DB.Save(&tag).Error
	if er != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Unable to edit tag",
		})
		return
	}

	tags, err := getTags()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "unable to fetch Tags",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Tag has been updated with id " + body.TagId,
		"tags":    tags,
	})
}

func getTags() ([]map[string]interface{}, error) {

	var tags []models.TagModel

	tagJson := make([]map[string]interface{}, len(tags))

	err := initializers.DB.Find(&tags).Error

	if err != nil {
		return nil, fmt.Errorf("unable to find tags: %w", err)
	}
	for _, tag := range tags {
		tagData := map[string]any{"tagName": tag.TagName, "id": tag.ID}
		tagJson = append(tagJson, tagData)
	}

	return tagJson, nil
}
