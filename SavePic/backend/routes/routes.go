package routes

import (
	"os"

	"savepic/backend/handlers"
	"savepic/backend/storage"

	"github.com/gin-gonic/gin"
)

// Setup 注册所有路由
func Setup(r *gin.Engine) {
	if storage.IsLocal() {
		_ = os.MkdirAll("uploads", 0755)
		r.Static("/uploads", "./uploads")
	}

	api := r.Group("/api")
	{
		api.GET("/categories", handlers.ListCategories)
		api.POST("/categories", handlers.CreateCategory)
		api.PUT("/categories/:id", handlers.UpdateCategory)
		api.DELETE("/categories/:id", handlers.DeleteCategory)

		api.GET("/tags", handlers.ListTags)
		api.PUT("/tags/:id", handlers.UpdateTag)

		api.GET("/memes", handlers.ListMemes)
		api.GET("/categories/:id/memes", handlers.ListMemesByCategory)
		api.POST("/memes/upload", handlers.UploadMeme)
		api.PUT("/memes/:id", handlers.UpdateMeme)
		api.DELETE("/memes/:id", handlers.DeleteMeme)
	}
}
