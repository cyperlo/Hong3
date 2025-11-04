package auth

import "errors"

var (
	ErrUserExists       = errors.New("用户名已存在")
	ErrInvalidCredentials = errors.New("用户名或密码错误")
	ErrInvalidToken     = errors.New("无效的token")
	ErrTokenExpired     = errors.New("token已过期")
	ErrUserNotFound     = errors.New("用户不存在")
)

