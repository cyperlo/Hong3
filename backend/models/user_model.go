package models

import (
	"crypto/rand"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Username  string    `gorm:"uniqueIndex;type:varchar(50);not null" json:"username"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-"` // 不返回密码
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// BeforeCreate 创建前钩子
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		// 如果 ID 为空，生成新的 ID
		u.ID = generateID()
	}
	return nil
}

// Token token 模型
type Token struct {
	Token     string    `gorm:"primaryKey;type:varchar(36)" json:"token"`
	UserID    string    `gorm:"type:varchar(36);not null;index" json:"user_id"`
	Username  string    `gorm:"type:varchar(50);not null" json:"username"`
	ExpiresAt time.Time `gorm:"not null;index" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 指定表名
func (Token) TableName() string {
	return "tokens"
}

// BeforeCreate 创建前钩子
func (t *Token) BeforeCreate(tx *gorm.DB) error {
	if t.Token == "" {
		// 如果 Token 为空，生成新的 Token
		t.Token = generateID()
	}
	return nil
}

// generateID 生成唯一ID
func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

