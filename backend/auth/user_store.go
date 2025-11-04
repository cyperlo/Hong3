package auth

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/chenhailong/hong3/db"
	"github.com/chenhailong/hong3/models"
	"github.com/chenhailong/hong3/redis"
)

// UserStore 用户存储（使用 GORM + Redis）
type UserStore struct {
}

var defaultStore *UserStore

// InitStore 初始化用户存储
func InitStore() {
	defaultStore = &UserStore{}
}

// GetStore 获取默认的用户存储
func GetStore() *UserStore {
	if defaultStore == nil {
		InitStore()
	}
	return defaultStore
}

// Register 注册新用户
func (s *UserStore) Register(username, password, name string) (*models.User, error) {
	// 检查用户名是否已存在
	var existingUser models.User
	result := db.DB.Where("username = ?", username).First(&existingUser)
	if result.Error == nil {
		// 用户已存在
		return nil, ErrUserExists
	}

	// 创建新用户
	user := &models.User{
		Username: username,
		Password: hashPassword(password),
		Name:     name,
	}

	if err := db.DB.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// Login 登录
func (s *UserStore) Login(username, password string) (*models.User, string, error) {
	// 查找用户
	var user models.User
	result := db.DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, "", ErrInvalidCredentials
	}

	// 验证密码
	if user.Password != hashPassword(password) {
		return nil, "", ErrInvalidCredentials
	}

	// 生成 token
	token := generateID()
	expiresAt := time.Now().Add(7 * 24 * time.Hour) // 7天过期
	expiration := 7 * 24 * time.Hour

	// 存储 token 到 Redis
	tokenData := &redis.TokenData{
		UserID:    user.ID,
		Username:  user.Username,
		ExpiresAt: expiresAt,
	}

	if err := redis.SetToken(token, tokenData, expiration); err != nil {
		return nil, "", fmt.Errorf("failed to save token: %w", err)
	}

	return &user, token, nil
}

// ValidateToken 验证token
func (s *UserStore) ValidateToken(token string) (*models.User, error) {
	// 从 Redis 获取 token 数据
	tokenData, err := redis.GetToken(token)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// 检查 token 是否过期
	if time.Now().After(tokenData.ExpiresAt) {
		// 删除过期的 token
		redis.DeleteToken(token)
		return nil, ErrTokenExpired
	}

	// 查找用户
	var user models.User
	result := db.DB.Where("id = ?", tokenData.UserID).First(&user)
	if result.Error != nil {
		return nil, ErrUserNotFound
	}

	return &user, nil
}

// GetUserByUsername 根据用户名获取用户
func (s *UserStore) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	result := db.DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, ErrUserNotFound
	}

	return &user, nil
}

// hashPassword 简单的密码哈希（生产环境应该使用bcrypt）
func hashPassword(password string) string {
	hash := md5.Sum([]byte(password))
	return hex.EncodeToString(hash[:])
}

// generateID 生成唯一ID
func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
