package database

import (
	"os"
	"path/filepath"
	"strings"

	"savepic/backend/models"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitFromEnv 优先使用 Postgres（Vercel），否则回退 SQLite（本地开发）
func InitFromEnv() error {
	if dsn := postgresDSN(); dsn != "" {
		return initPostgres(dsn)
	}
	dbPath := strings.TrimSpace(os.Getenv("SQLITE_PATH"))
	if dbPath == "" {
		dbPath = "data/savepic.db"
	}
	return InitSQLite(dbPath)
}

func postgresDSN() string {
	for _, key := range []string{
		"DATABASE_URL",
		"POSTGRES_URL",
		"POSTGRES_URL_NON_POOLING",
	} {
		if v := strings.TrimSpace(os.Getenv(key)); v != "" {
			return v
		}
	}
	return ""
}

func initPostgres(dsn string) error {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	return migrate(db)
}

// InitSQLite 连接 SQLite 并自动迁移表结构
func InitSQLite(dbPath string) error {
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
	return migrate(db)
}

// Init 兼容旧调用，等价于 InitSQLite
func Init(dbPath string) error {
	return InitSQLite(dbPath)
}

func migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.Category{}, &models.Meme{}, &models.Tag{}); err != nil {
		return err
	}
	DB = db
	return nil
}
