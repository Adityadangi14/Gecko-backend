package main

import (
	"gecko_backend/controllers"
	"gecko_backend/initializers"
	"gecko_backend/middleware"

	"github.com/gin-gonic/gin"
)

func init() {

	gin.SetMode(gin.ReleaseMode)
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDB()
	initializers.ConnectToRedis()
	initializers.DB.Config.QueryFields = true
}

func main() {
	r := gin.Default()
	// blog
	r.POST("/createBlog", controllers.CreateBlog)
	r.DELETE("/deleteBlog", controllers.DeleteBlog)
	r.GET("/getBlogs", controllers.GetBlogs)
	r.POST("/getBlogsByTags", controllers.GetBlogsByTags)
	// trending blog
	r.GET("/getTrendingBlogs",middleware.RequiredAuth, controllers.GetTrendingBlogs)
	r.POST("/addTrendingBlog",  controllers.AddTrendingBlog)
	r.DELETE("/deleteTrendingBlog",  controllers.DeleteTrendingBlog)
	// category
	r.POST("/addBlogCategory",middleware.RequiredAuth,  controllers.CreateBlogCategory)
	r.GET("/getBlogCategory",middleware.RequiredAuth,  controllers.GetBlogCategories)
	r.POST("/getBlogByCategory",middleware.RequiredAuth,  controllers.GetBlogsByCat)
	// companies
	r.POST("/createCompany",middleware.RequiredAuth,  controllers.CreateCompany)
	r.GET("/getCompanies", middleware.RequiredAuth, controllers.GetCompanies)
	r.DELETE("/deleteCompany",middleware.RequiredAuth,  controllers.DeleteCompany)
	r.PUT("/editCompany",middleware.RequiredAuth,  controllers.EditCompany)
	// tag
	r.POST("/createTag",middleware.RequiredAuth,  controllers.CreateTag)
	r.DELETE("/deleteTag",middleware.RequiredAuth, controllers.DeleteTag)
	r.GET("/getTags",middleware.RequiredAuth, controllers.GetAllTags)
	r.PUT("/editTag",middleware.RequiredAuth, controllers.EditTag)
	r.POST("/uploadThumblanilFile", controllers.UploadThumbnailFiles)
	// user
	r.GET("/getUser",middleware.RequiredAuth,  controllers.GetUser)
	r.PUT("/editUserTagSelectionStatus",middleware.RequiredAuth, controllers.EditUserTagSelectionStatus)
	r.POST("/createUserTags",middleware.RequiredAuth, controllers.CreateUserTags)

	r.POST("/auth", controllers.AuthHandler)

	r.Run(":3000")
}
