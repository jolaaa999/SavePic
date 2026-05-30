package sqlitedb

import (
	"os"
	"path/filepath"
	"strings"

	"savepic/backend/database"
	"savepic/backend/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

/**
 * Init 连接本地 SQLite 并自动迁移表结构（仅本地开发使用）。
 */
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

	database.SetDB(db)
	return nil
}

/**
 * InitFromEnv 按环境变量 SQLITE_PATH 初始化 SQLite。
 */
func InitFromEnv() error {
	dbPath := strings.TrimSpace(os.Getenv("SQLITE_PATH"))
	if dbPath == "" {
		dbPath = "data/savepic.db"
	}
	return Init(dbPath)
}
