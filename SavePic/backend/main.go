package main

import (
	"log"
	"net/http"

	"savepic/database"
	"savepic/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := database.Init("data/savepic.db"); err != nil {
		log.Fatalf("database init failed: %v", err)
	}

	r := gin.Default()
	r.Use(corsMiddleware())
	routes.Setup(r)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "success",
			"data": gin.H{"status": "ok"},
		})
	})

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("server start failed: %v", err)
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
