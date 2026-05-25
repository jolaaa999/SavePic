package database

import (
	"os"
	"path/filepath"

	"savepic/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Init 连接 SQLite 并自动迁移表结构
func Init(dbPath string) error {
	if dbPath == "" {
		dbPath = "data/savepic.db"
	}

	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return err
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return err
	}

	if err := db.AutoMigrate(&models.Category{}, &models.Meme{}, &models.Tag{}); err != nil {
		return err
	}

	DB = db
	return nil
}
