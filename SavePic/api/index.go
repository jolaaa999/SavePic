package handler

import (
	"encoding/json"
	"net/http"
	"sync"

	"savepic/backend/app"
	"savepic/backend/database"
)

var (
	engine  http.Handler
	initErr error
	once    sync.Once
)

func ensureEngine() http.Handler {
	once.Do(func() {
		if err := database.InitFromEnv(); err != nil {
			initErr = err
			return
		}
		engine = app.NewEngine()
	})
	return engine
}

/**
 * Handler 供 Vercel Go Serverless 调用，处理 /api/* 等后端路由。
 */
func Handler(w http.ResponseWriter, r *http.Request) {
	if h := ensureEngine(); h != nil {
		h.ServeHTTP(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"code": 500,
		"msg":  "服务初始化失败: " + initErr.Error(),
		"data": nil,
	})
}
