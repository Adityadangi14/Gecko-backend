package controllers

import (
	"gecko_backend/initializers"
	"gecko_backend/models"
	"strconv"

	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateCompany(c *gin.Context) {
	var body struct {
		CompanyName    string
		CompanyLogoURL string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	blog := models.CompanyModel{CompanyName: body.CompanyName, CompanyLogoURL: body.CompanyLogoURL}

	result := initializers.DB.Create(&blog)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create company",
		})
		return
	}

	var companies []models.CompanyModel

	initializers.DB.Find(&companies)

	companyjson := make([]map[string]interface{}, 0)
	for _, company := range companies {
		companyData := map[string]any{"id": company.ID, "company": company.CompanyName, "companyLogo": company.CompanyLogoURL}
		companyjson = append(companyjson, companyData)
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"message":   "Company has been added with id " + strconv.FormatUint(uint64(blog.ID), 10),
		"compaines": companyjson,
	})
}

func GetCompanies(c *gin.Context) {

	var companies []models.CompanyModel

	initializers.DB.Find(&companies)

	if len(companies) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"success":   true,
			"compaines": make([]models.CompanyModel, 0),
		})

		return
	}

	companyjson := make([]map[string]interface{}, 0)
	for _, company := range companies {
		companyData := map[string]any{"id": company.ID, "company": company.CompanyName, "companyLogo": company.CompanyLogoURL}
		companyjson = append(companyjson, companyData)
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"compaines": companyjson,
	})

}

func DeleteCompany(c *gin.Context) {
	var body struct {
		CompanyId uint
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Invalid request body format",
		})
		return
	}
	var company models.CompanyModel
	err := initializers.DB.Where("id = ?", body.CompanyId).Delete(&company).Error

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"success": false,
			"message": "unable to delete the company :" + err.Error(),
		})
		return
	}

	var companies []models.CompanyModel

	initializers.DB.Find(&companies)
	companyjson := make([]map[string]interface{}, 0)
	for _, company := range companies {
		companyData := map[string]any{"id": company.ID, "company": company.CompanyName, "companyLogo": company.CompanyLogoURL}
		companyjson = append(companyjson, companyData)
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"message":   "Compnay with id " + strconv.FormatUint(uint64(body.CompanyId), 10) + " has been deleted successfully",
		"compaines": companyjson,
	})
}

func EditCompany(c *gin.Context) {
	var body struct {
		companyId      uint
		companyName    string
		companyLogoURL string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	var company models.CompanyModel

	result := initializers.DB.Where("id = ?", body.companyId).First(&company)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	company.CompanyLogoURL = body.companyLogoURL
	company.CompanyName = body.companyName
	company.ID = body.companyId

	err := initializers.DB.Save(&company)

	if err != nil {
		var companies []models.CompanyModel

		initializers.DB.Find(&companies)
		companyjson := make([]map[string]interface{}, 0)
		for _, company := range companies {
			companyData := map[string]any{"id": company.ID, "company": company.CompanyName, "companyLogo": company.CompanyLogoURL}
			companyjson = append(companyjson, companyData)
		}

		c.JSON(http.StatusOK, gin.H{
			"success":   true,
			"message":   "Compnay with id " + strconv.FormatUint(uint64(body.companyId), 10) + " has been updated successfully",
			"compaines": companyjson,
		})
	}
}
