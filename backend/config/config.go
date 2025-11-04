package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/yaml.v3"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Redis    RedisConfig    `json:"redis"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string `json:"port"`
	Host string `json:"host"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
	SSLMode  string `json:"ssl_mode"`
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
	Enabled  bool   `json:"enabled"`
}

var AppConfig *Config

// LoadConfig 加载配置
// 优先级：环境变量 > 配置文件 > 默认值
func LoadConfig() *Config {
	config := &Config{}

	// 首先尝试从配置文件加载
	loadFromFile(config)

	// 然后用环境变量覆盖（环境变量优先级最高）
	applyEnvOverrides(config)

	AppConfig = config
	logConfig(config)
	return config
}

// loadFromFile 从配置文件加载
func loadFromFile(config *Config) {
	// 查找配置文件，按优先级：
	// 1. 当前目录的 config.yaml
	// 2. backend 目录的 config.yaml
	configPaths := []string{
		"config.yaml",
		filepath.Join("backend", "config.yaml"),
		filepath.Join(".", "backend", "config.yaml"),
	}

	var configFile string
	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			configFile = path
			break
		}
	}

	if configFile == "" {
		log.Println("Config file not found, using defaults and environment variables")
		return
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		log.Printf("Warning: Failed to read config file %s: %v, using defaults", configFile, err)
		return
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		log.Printf("Warning: Failed to parse config file %s: %v, using defaults", configFile, err)
		return
	}

	log.Printf("Configuration loaded from file: %s", configFile)
}

// applyEnvOverrides 应用环境变量覆盖配置
func applyEnvOverrides(config *Config) {
	// 服务器配置
	if port := getEnv("PORT", ""); port != "" {
		config.Server.Port = port
	}
	if host := getEnv("HOST", ""); host != "" {
		config.Server.Host = host
	}

	// 数据库配置
	if host := getEnv("DB_HOST", ""); host != "" {
		config.Database.Host = host
	}
	if port := getEnv("DB_PORT", ""); port != "" {
		config.Database.Port = port
	}
	if user := getEnv("DB_USER", ""); user != "" {
		config.Database.User = user
	}
	if password := getEnv("DB_PASSWORD", ""); password != "" {
		config.Database.Password = password
	}
	if name := getEnv("DB_NAME", ""); name != "" {
		config.Database.Name = name
	}
	if sslMode := getEnv("DB_SSLMODE", ""); sslMode != "" {
		config.Database.SSLMode = sslMode
	}

	// Redis 配置
	if host := getEnv("REDIS_HOST", ""); host != "" {
		config.Redis.Host = host
	}
	if port := getEnv("REDIS_PORT", ""); port != "" {
		config.Redis.Port = port
	}
	if password := getEnv("REDIS_PASSWORD", ""); password != "" {
		config.Redis.Password = password
	}
	if db := getEnv("REDIS_DB", ""); db != "" {
		if dbInt := getEnvAsInt("REDIS_DB", -1); dbInt >= 0 {
			config.Redis.DB = dbInt
		}
	}
	if enabled := getEnv("REDIS_ENABLED", ""); enabled != "" {
		config.Redis.Enabled = getEnvAsBool("REDIS_ENABLED", false)
	}

	// 设置默认值（如果配置文件和环境变量都没有设置）
	setDefaults(config)
}

// setDefaults 设置默认值
func setDefaults(config *Config) {
	if config.Server.Port == "" {
		config.Server.Port = "8080"
	}
	if config.Server.Host == "" {
		config.Server.Host = "0.0.0.0"
	}

	if config.Database.Host == "" {
		config.Database.Host = "postgres"
	}
	if config.Database.Port == "" {
		config.Database.Port = "5432"
	}
	if config.Database.User == "" {
		config.Database.User = "postgres"
	}
	if config.Database.Password == "" {
		config.Database.Password = "postgres"
	}
	if config.Database.Name == "" {
		config.Database.Name = "hong3"
	}
	if config.Database.SSLMode == "" {
		config.Database.SSLMode = "disable"
	}

	if config.Redis.Host == "" {
		config.Redis.Host = "redis"
	}
	if config.Redis.Port == "" {
		config.Redis.Port = "6379"
	}
}

// GetDSN 获取数据库连接字符串
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode)
}

// GetRedisURL 获取 Redis 连接 URL
func (c *RedisConfig) GetRedisURL() string {
	if c.Password != "" {
		return "redis://:" + c.Password + "@" + c.Host + ":" + c.Port + "/" + strconv.Itoa(c.DB)
	}
	return "redis://" + c.Host + ":" + c.Port + "/" + strconv.Itoa(c.DB)
}

// GetServerAddr 获取服务器地址
func (c *ServerConfig) GetServerAddr() string {
	return c.Host + ":" + c.Port
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvAsInt 获取环境变量并转换为整数
func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("Warning: Invalid integer value for %s, using default: %d", key, defaultValue)
		return defaultValue
	}
	return intValue
}

// getEnvAsBool 获取环境变量并转换为布尔值
func getEnvAsBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		log.Printf("Warning: Invalid boolean value for %s, using default: %v", key, defaultValue)
		return defaultValue
	}
	return boolValue
}

// logConfig 打印配置信息（不打印敏感信息）
func logConfig(config *Config) {
	log.Println("=== Application Configuration ===")
	log.Printf("Server: %s:%s", config.Server.Host, config.Server.Port)
	log.Printf("Database: %s@%s:%s/%s", config.Database.User, config.Database.Host, config.Database.Port, config.Database.Name)
	if config.Redis.Enabled {
		log.Printf("Redis: %s:%s (DB: %d)", config.Redis.Host, config.Redis.Port, config.Redis.DB)
	} else {
		log.Println("Redis: Disabled")
	}
	log.Println("=================================")
}
