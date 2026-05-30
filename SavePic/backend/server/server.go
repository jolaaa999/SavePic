package server

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"savepic/backend/app"
	"savepic/backend/database"
	"savepic/backend/sqlitedb"

	"github.com/gin-gonic/gin"
)

// Run 启动 HTTP 服务（本地开发：Postgres 或 SQLite）
func Run() {
	if err := database.InitFromEnv(); err != nil {
		if err := sqlitedb.InitFromEnv(); err != nil {
			panic(err)
		}
	}

	r := app.NewEngine()
	registerStatic(r, resolveDistDir())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := r.Run(":" + port); err != nil {
		panic(err)
	}
}

func resolveDistDir() string {
	if d := strings.TrimSpace(os.Getenv("STATIC_DIR")); d != "" {
		return d
	}
	for _, candidate := range []string{"frontend/dist", "../frontend/dist"} {
		if _, err := os.Stat(filepath.Join(candidate, "index.html")); err == nil {
			return candidate
		}
	}
	return "frontend/dist"
}

func registerStatic(r *gin.Engine, dist string) {
	assetsDir := filepath.Join(dist, "assets")
	if info, err := os.Stat(assetsDir); err == nil && info.IsDir() {
		r.Static("/assets", assetsDir)
	}

	for _, name := range []string{"favicon.svg", "icons.svg"} {
		p := filepath.Join(dist, name)
		if _, err := os.Stat(p); err == nil {
			r.StaticFile("/"+name, p)
		}
	}

	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.JSON(http.StatusNotFound, gin.H{
				"code": 404,
				"msg":  "接口不存在",
				"data": nil,
			})
			return
		}

		rel := strings.TrimPrefix(c.Request.URL.Path, "/")
		if rel != "" {
			candidate := filepath.Join(dist, filepath.FromSlash(rel))
			if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
				c.File(candidate)
				return
			}
		}

		c.File(filepath.Join(dist, "index.html"))
	})
}
