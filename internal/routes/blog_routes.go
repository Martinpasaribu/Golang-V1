package routes

import (
	blogController "github.com/Martinpasaribu/Golang-V1/internal/controllers/blog"
	"github.com/gin-gonic/gin"
)

func RegisterBlogRoutes(router *gin.RouterGroup, blogCtrl *blogController.BlogController) {
	blogGroup := router.Group("/blogs")
	{
		blogGroup.POST("", blogCtrl.CreateBlog)
		blogGroup.POST("/image", blogCtrl.UploadBlogImages)
		
		// blogGroup.GET("/id/:id", blogCtrl.GetBlogByID)
		// blogGroup.PUT("/:id", blogCtrl.UpdateBlog)
		// blogGroup.DELETE("/:id", blogCtrl.DeleteBlog)

		blogGroup.GET("/data", blogCtrl.GetAllBlogs)
		blogGroup.GET("/detail-blog/:slug", blogCtrl.FindBlogBySlug)
		blogGroup.GET("/detail-category", blogCtrl.GetCategory)
		blogGroup.GET("/detail-category-navbar", blogCtrl.GetCategoryNavbar)
		blogGroup.GET("/articles-one", blogCtrl.GetArticles01)
		blogGroup.GET("/articles-list", blogCtrl.GetArticlesList)
	}
}