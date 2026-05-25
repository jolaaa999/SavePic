package handler

import (
	"log"
	"net/http"
	"sync"

	"savepic/backend/app"
	"savepic/backend/database"
)

var (
	initOnce sync.Once
	engine   http.Handler
	initErr  error
)

func initApp() {
	initOnce.Do(func() {
		if err := database.InitFromEnv(); err != nil {
			initErr = err
			log.Printf("database init failed: %v", err)
			return
		}
		engine = app.NewEngine()
	})
}

// Handler Vercel Serverless 入口（Gin 引擎转为 http.HandlerFunc）
func Handler(w http.ResponseWriter, r *http.Request) {
	initApp()
	if initErr != nil {
		http.Error(w, "database unavailable", http.StatusInternalServerError)
		return
	}
	engine.ServeHTTP(w, r)
}
