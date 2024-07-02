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
	initializers.DB.Config.QueryFields = true
}

func main() {
	r := gin.Default()

	r.POST("/createBlog", middleware.RequiredAuth, controllers.CreateBlog)
	r.DELETE("/deleteBlog", middleware.RequiredAuth, controllers.DeleteBlog)
	r.GET("/getBlogs", middleware.RequiredAuth, controllers.GetBlogs)
	r.POST("/getBlogsByTags", middleware.RequiredAuth, controllers.GetBlogsByTags)
	r.GET("/getTrendingBlogs", middleware.RequiredAuth, controllers.GetTrendingBlogs)
	r.POST("/addTrendingBlog", middleware.RequiredAuth, controllers.AddTrendingBlog)
	r.DELETE("/deleteTrendingBlog", middleware.RequiredAuth, controllers.DeleteTrendingBlog)
	r.POST("/createCompany", middleware.RequiredAuth, controllers.CreateCompany)
	r.GET("/getCompanies", middleware.RequiredAuth, controllers.GetCompanies)
	r.DELETE("/deleteCompany", middleware.RequiredAuth, controllers.DeleteCompany)
	r.PUT("/editCompany", middleware.RequiredAuth, controllers.EditCompany)
	r.POST("/createTag", middleware.RequiredAuth, controllers.CreateTag)
	r.DELETE("/deleteTag", middleware.RequiredAuth, controllers.DeleteTag)
	r.GET("/getTags", middleware.RequiredAuth, controllers.GetAllTags)
	r.PUT("/editTag", middleware.RequiredAuth, controllers.EditTag)
	r.POST("/uploadThumblanilFile", controllers.UploadThumbnailFiles)
	r.GET("/getUser", middleware.RequiredAuth, controllers.GetUser)
	r.PUT("/editUserTagSelectionStatus", middleware.RequiredAuth, controllers.EditUserTagSelectionStatus)
	r.POST("/createUserTags", middleware.RequiredAuth, controllers.CreateUserTags)

	r.POST("/auth", controllers.AuthHandler)

	r.Run(":3000")
}
