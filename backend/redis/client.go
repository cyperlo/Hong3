package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/chenhailong/hong3/config"
	redispkg "github.com/redis/go-redis/v9"
)

var Client *redispkg.Client
var ctx = context.Background()

// TokenData token 数据
type TokenData struct {
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	ExpiresAt time.Time `json:"expires_at"`
}

// InitRedis 初始化 Redis 连接
func InitRedis() (*redispkg.Client, error) {
	cfg := config.AppConfig
	if cfg == nil {
		cfg = config.LoadConfig()
	}

	if !cfg.Redis.Enabled {
		log.Println("Redis is disabled in configuration")
		return nil, fmt.Errorf("redis is disabled")
	}

	port, err := strconv.Atoi(cfg.Redis.Port)
	if err != nil {
		return nil, fmt.Errorf("invalid redis port: %w", err)
	}

	Client = redispkg.NewClient(&redispkg.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// 测试连接
	_, err = Client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	log.Println("Redis connected successfully")
	return Client, nil
}

// SetToken 设置 token
func SetToken(token string, data *TokenData, expiration time.Duration) error {
	if Client == nil {
		return fmt.Errorf("redis client not initialized")
	}

	key := fmt.Sprintf("token:%s", token)
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal token data: %w", err)
	}

	err = Client.Set(ctx, key, dataJSON, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set token: %w", err)
	}

	return nil
}

// GetToken 获取 token 数据
func GetToken(token string) (*TokenData, error) {
	if Client == nil {
		return nil, fmt.Errorf("redis client not initialized")
	}

	key := fmt.Sprintf("token:%s", token)
	val, err := Client.Get(ctx, key).Result()
	if err != nil {
		if err == redispkg.Nil {
			return nil, fmt.Errorf("token not found")
		}
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	var data TokenData
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token data: %w", err)
	}

	return &data, nil
}

// DeleteToken 删除 token
func DeleteToken(token string) error {
	if Client == nil {
		return fmt.Errorf("redis client not initialized")
	}

	key := fmt.Sprintf("token:%s", token)
	err := Client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}

	return nil
}

// TokenExists 检查 token 是否存在
func TokenExists(token string) (bool, error) {
	if Client == nil {
		return false, fmt.Errorf("redis client not initialized")
	}

	key := fmt.Sprintf("token:%s", token)
	count, err := Client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check token: %w", err)
	}

	return count > 0, nil
}

// CleanExpiredTokens 清理过期 token（Redis 会自动过期，此方法用于清理模式匹配的 key）
func CleanExpiredTokens() error {
	if Client == nil {
		return fmt.Errorf("redis client not initialized")
	}

	// Redis 会自动删除过期的 key，这里可以扫描并清理（如果需要）
	// 通常不需要手动清理，因为 Redis 的 TTL 会自动处理
	return nil
}
