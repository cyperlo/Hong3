#!/bin/bash
# 安装 Redis 客户端依赖

echo "Installing Redis dependencies..."

cd "$(dirname "$0")/.."

go get github.com/redis/go-redis/v9
go mod tidy

echo "Redis dependencies installed successfully!"

