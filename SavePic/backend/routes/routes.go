package routes

import (
	"os"

	"savepic/handlers"

	"github.com/gin-gonic/gin"
)

// Setup 注册所有路由
func Setup(r *gin.Engine) {
	_ = os.MkdirAll("uploads", 0755)

	r.Static("/uploads", "./uploads")

	api := r.Group("/api")
	{
		api.GET("/categories", handlers.ListCategories)
		api.POST("/categories", handlers.CreateCategory)

		api.GET("/tags", handlers.ListTags)

		api.GET("/memes", handlers.ListMemes)
		api.GET("/categories/:id/memes", handlers.ListMemesByCategory)
		api.POST("/memes/upload", handlers.UploadMeme)
		api.PUT("/memes/:id", handlers.UpdateMeme)
		api.DELETE("/memes/:id", handlers.DeleteMeme)
	}
}
