package main

import (
	"log"

	"savepic/backend/app"
	"savepic/backend/database"
)

func main() {
	if err := database.InitFromEnv(); err != nil {
		log.Fatalf("database init failed: %v", err)
	}

	r := app.NewEngine()
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("server start failed: %v", err)
	}
}
