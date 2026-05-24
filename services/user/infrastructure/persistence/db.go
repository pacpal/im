// Package persistence 提供数据库连接与 GORM 封装（此文件为 Postgres 的实现）。
package persistence

import (
	"IM/pkg/config"
	zlog "IM/pkg/logger"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// PostgresDB 封装了 *gorm.DB，便于统一管理连接与生命周期。
type PostgresDB struct {
	db *gorm.DB
}

// NewPostgresDB 使用配置创建并返回 Postgres 数据库连接封装。
func NewPostgresDB(cfg config.DatabaseConfig) (res *PostgresDB, err error) {
	done := zlog.StartStep("persistence.NewPostgresDB", "host", cfg.Host, "db", cfg.DBName)
	defer func() { done(err) }()

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		err = fmt.Errorf("failed to connect database: %w", err)
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		err = fmt.Errorf("failed to get sql.DB: %w", err)
		return
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpen)
	sqlDB.SetMaxIdleConns(cfg.MaxIdle)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if e := sqlDB.Ping(); e != nil {
		err = fmt.Errorf("failed to ping database: %w", e)
		return
	}

	res = &PostgresDB{db: db}
	return
}

func (p *PostgresDB) GetDB() *gorm.DB {
	return p.db
}

func (p *PostgresDB) Close() error {
	sqlDB, err := p.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
