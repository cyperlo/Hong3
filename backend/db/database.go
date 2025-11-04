package db

import (
	"fmt"
	"log"
	"time"

	"github.com/chenhailong/hong3/config"
	"github.com/chenhailong/hong3/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() (*gorm.DB, error) {
	// 从配置获取数据库连接信息
	cfg := config.AppConfig
	if cfg == nil {
		cfg = config.LoadConfig()
	}

	dsn := cfg.Database.GetDSN()

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 获取底层 sql.DB 连接池
	sqlDB, err := DB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connected successfully")
	return DB, nil
}

// AutoMigrate 自动迁移数据库表
func AutoMigrate() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	// 自动迁移所有模型（token 已迁移到 Redis，不再需要 tokens 表）
	err := DB.AutoMigrate(
		&models.User{},
		// &models.Token{}, // Token 已迁移到 Redis
	)

	if err != nil {
		return fmt.Errorf("failed to auto migrate: %w", err)
	}

	log.Println("Database migration completed")
	return nil
}
