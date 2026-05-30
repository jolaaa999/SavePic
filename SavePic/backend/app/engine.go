package app

import (
	"net/http"

	"savepic/backend/routes"
	"savepic/backend/storage"

	"github.com/gin-gonic/gin"
)

// NewEngine 创建并配置 Gin 引擎（本地与 Vercel Serverless 共用）
func NewEngine() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery(), oidcMiddleware(), corsMiddleware())
	routes.Setup(r)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "success",
			"data": gin.H{"status": "ok"},
		})
	})

	return r
}

func oidcMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if token := c.GetHeader("x-vercel-oidc-token"); token != "" {
			ctx := storage.WithOIDCToken(c.Request.Context(), token)
			c.Request = c.Request.WithContext(ctx)
		}
		c.Next()
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
