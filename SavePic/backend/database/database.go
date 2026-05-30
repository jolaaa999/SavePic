package database

import (
	"os"
	"strings"

	"savepic/backend/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

/**
 * InitFromEnv 连接 Postgres（Vercel 生产环境）。
 * 本地开发若未配置 DATABASE_URL，由 server 包回退到 SQLite。
 */
func InitFromEnv() error {
	dsn := postgresDSN()
	if dsn == "" {
		return ErrPostgresRequired
	}
	return initPostgres(dsn)
}

/** ErrPostgresRequired 未配置 Postgres 连接串 */
var ErrPostgresRequired = &initError{msg: "未配置 DATABASE_URL 或 POSTGRES_URL"}

type initError struct {
	msg string
}

func (e *initError) Error() string {
	return e.msg
}

func postgresDSN() string {
	for _, key := range []string{
		"POSTGRES_URL_NON_POOLING",
		"DATABASE_URL",
		"POSTGRES_URL",
	} {
		if v := strings.TrimSpace(os.Getenv(key)); v != "" {
			return v
		}
	}
	return ""
}

func initPostgres(dsn string) error {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxOpenConns(2)
	sqlDB.SetMaxIdleConns(1)

	return migrate(db)
}

/**
 * SetDB 供本地 SQLite 初始化后注入全局 DB 实例。
 */
func SetDB(db *gorm.DB) {
	DB = db
}

func migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.Category{}, &models.Meme{}, &models.Tag{}); err != nil {
		return err
	}
	DB = db
	return nil
}
