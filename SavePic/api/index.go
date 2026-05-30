package handler

import (
	"net/http"

	"savepic/backend/app"
	"savepic/backend/database"
)

var engine http.Handler

func init() {
	if err := database.InitFromEnv(); err != nil {
		panic(err)
	}
	engine = app.NewEngine()
}

/**
 * Handler 供 Vercel Go Serverless 调用，处理 /api/* 等后端路由。
 */
func Handler(w http.ResponseWriter, r *http.Request) {
	engine.ServeHTTP(w, r)
}
