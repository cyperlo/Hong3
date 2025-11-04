# 部署指南

## 使用 Docker Compose 部署

### 快速开始

1. **克隆项目**
```bash
git clone <repository-url>
cd Hong3
```

2. **构建并启动所有服务**
```bash
docker-compose up -d
```

3. **查看服务状态**
```bash
docker-compose ps
```

4. **查看日志**
```bash
# 查看所有服务日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f backend
docker-compose logs -f frontend
```

5. **停止服务**
```bash
docker-compose down
```

### 生产环境部署

使用生产环境配置：

```bash
docker-compose -f docker-compose.prod.yml up -d
```

### 重新构建

如果代码有更新，需要重新构建：

```bash
# 重新构建并启动
docker-compose up -d --build

# 只重新构建某个服务
docker-compose build backend
docker-compose up -d backend
```

### 访问服务

- 前端: http://localhost
- 后端 API: http://localhost:8080
- WebSocket: ws://localhost:8080/ws

### 健康检查

检查后端服务健康状态：

```bash
curl http://localhost:8080/api/rooms
```

### 清理

```bash
# 停止并删除容器
docker-compose down

# 停止并删除容器、网络、卷
docker-compose down -v

# 删除所有相关镜像
docker-compose down --rmi all
```

## 单独构建 Docker 镜像

### 构建后端镜像

```bash
cd backend
docker build -t hong3-backend:latest .
```

### 构建前端镜像

```bash
cd frontend
docker build -t hong3-frontend:latest .
```

## 环境变量配置

如果需要使用 Redis 或 PostgreSQL，取消 `docker-compose.yml` 中对应服务的注释，并配置环境变量。

## 故障排查

1. **查看容器日志**
```bash
docker-compose logs -f [service-name]
```

2. **进入容器调试**
```bash
docker-compose exec backend sh
docker-compose exec frontend sh
```

3. **检查网络连接**
```bash
docker network inspect hong3_hong3-network
```

4. **检查端口占用**
```bash
# Linux/Mac
lsof -i :80
lsof -i :8080

# Windows
netstat -ano | findstr :80
netstat -ano | findstr :8080
```

