package main

import (
	"log"

	"github.com/chenhailong/hong3/api"
	"github.com/chenhailong/hong3/auth"
	"github.com/chenhailong/hong3/config"
	"github.com/chenhailong/hong3/db"
	"github.com/chenhailong/hong3/redis"
)

func main() {
	// 加载配置
	log.Println("Loading configuration...")
	cfg := config.LoadConfig()

	// 初始化数据库
	log.Println("Initializing database...")
	if _, err := db.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 自动迁移数据库表
	if err := db.AutoMigrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 初始化 Redis（如果启用）
	if cfg.Redis.Enabled {
		log.Println("Initializing Redis...")
		if _, err := redis.InitRedis(); err != nil {
			log.Fatalf("Failed to initialize Redis: %v", err)
		}
	} else {
		log.Println("Redis is disabled, token storage will fail. Please enable Redis in configuration.")
	}

	// 初始化用户存储
	auth.InitStore()

	// 启动服务器
	server := api.NewServer()
	serverAddr := cfg.Server.GetServerAddr()
	log.Printf("Starting Hong3 server on %s", serverAddr)
	if err := server.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}