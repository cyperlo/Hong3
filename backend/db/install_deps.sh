#!/bin/bash
# 安装 GORM 和 PostgreSQL 驱动依赖

echo "Installing GORM dependencies..."

cd "$(dirname "$0")/.."

go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
go mod tidy

echo "Dependencies installed successfully!"

