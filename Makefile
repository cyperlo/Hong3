.PHONY: build up down restart logs clean help

# 默认目标
.DEFAULT_GOAL := help

# 构建所有镜像
build:
	docker-compose build

# 启动所有服务
up:
	docker-compose up -d

# 启动并查看日志
up-logs:
	docker-compose up

# 停止所有服务
down:
	docker-compose down

# 重启服务
restart:
	docker-compose restart

# 查看日志
logs:
	docker-compose logs -f

# 查看后端日志
logs-backend:
	docker-compose logs -f backend

# 查看前端日志
logs-frontend:
	docker-compose logs -f frontend

# 清理（停止并删除容器、网络）
clean:
	docker-compose down -v

# 完全清理（包括镜像）
clean-all:
	docker-compose down -v --rmi all

# 重新构建并启动
rebuild:
	docker-compose up -d --build

# 生产环境部署
prod-up:
	docker-compose -f docker-compose.prod.yml up -d

# 生产环境构建
prod-build:
	docker-compose -f docker-compose.prod.yml build

# 生产环境重新构建并启动
prod-rebuild:
	docker-compose -f docker-compose.prod.yml up -d --build

# 查看服务状态
ps:
	docker-compose ps

# 进入后端容器
shell-backend:
	docker-compose exec backend sh

# 进入前端容器
shell-frontend:
	docker-compose exec frontend sh

# 帮助信息
help:
	@echo "可用的命令:"
	@echo "  make build          - 构建所有镜像"
	@echo "  make up             - 启动所有服务（后台）"
	@echo "  make up-logs        - 启动所有服务（前台，显示日志）"
	@echo "  make down           - 停止所有服务"
	@echo "  make restart        - 重启所有服务"
	@echo "  make logs           - 查看所有服务日志"
	@echo "  make logs-backend   - 查看后端日志"
	@echo "  make logs-frontend  - 查看前端日志"
	@echo "  make clean          - 清理（停止并删除容器、网络）"
	@echo "  make clean-all      - 完全清理（包括镜像）"
	@echo "  make rebuild        - 重新构建并启动"
	@echo "  make ps             - 查看服务状态"
	@echo "  make shell-backend  - 进入后端容器"
	@echo "  make shell-frontend - 进入前端容器"
	@echo ""
	@echo "生产环境命令:"
	@echo "  make prod-up        - 生产环境启动"
	@echo "  make prod-build     - 生产环境构建"
	@echo "  make prod-rebuild   - 生产环境重新构建并启动"

