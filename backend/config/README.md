# 配置说明

## 概述

配置系统**仅从环境变量读取配置**，不再支持配置文件。所有配置都通过 `docker-compose.yml` 中的环境变量设置。

## 环境变量

### 服务器配置

- `PORT`: 服务器端口（默认: 8080）
- `HOST`: 服务器主机（默认: 0.0.0.0）

### 数据库配置（PostgreSQL）

- `DB_HOST`: 数据库主机（默认: postgres）
- `DB_PORT`: 数据库端口（默认: 5432）
- `DB_USER`: 数据库用户名（默认: postgres）
- `DB_PASSWORD`: 数据库密码（默认: postgres）
- `DB_NAME`: 数据库名称（默认: hong3）
- `DB_SSLMODE`: SSL 模式（默认: disable）

**PostgreSQL 容器环境变量：**
- `POSTGRES_USER`: PostgreSQL 用户名（默认: postgres）
- `POSTGRES_PASSWORD`: PostgreSQL 密码（默认: postgres）
- `POSTGRES_DB`: PostgreSQL 数据库名（默认: hong3）

**注意**：`DB_USER` 和 `DB_PASSWORD` 必须与 `POSTGRES_USER` 和 `POSTGRES_PASSWORD` 保持一致。

### Redis 配置

- `REDIS_HOST`: Redis 主机（默认: redis）
- `REDIS_PORT`: Redis 端口（默认: 6379）
- `REDIS_PASSWORD`: Redis 密码（默认: 空，无密码）
- `REDIS_DB`: Redis 数据库编号（默认: 0）
- `REDIS_ENABLED`: 是否启用 Redis（默认: false）

## 使用方法

### 在 docker-compose.yml 中配置

```yaml
services:
  postgres:
    environment:
      - POSTGRES_USER=myuser          # 自定义数据库用户名
      - POSTGRES_PASSWORD=mypassword  # 自定义数据库密码
      - POSTGRES_DB=hong3

  backend:
    environment:
      - DB_USER=myuser                # 必须与 POSTGRES_USER 一致
      - DB_PASSWORD=mypassword        # 必须与 POSTGRES_PASSWORD 一致
      - REDIS_ENABLED=true
```

### 使用环境变量文件（推荐）

创建 `.env` 文件（不应提交到版本控制）：

```bash
# PostgreSQL 配置
POSTGRES_USER=myuser
POSTGRES_PASSWORD=mypassword
POSTGRES_DB=hong3

# 后端数据库配置（需与 POSTGRES_* 一致）
DB_USER=myuser
DB_PASSWORD=mypassword
DB_NAME=hong3

# Redis 配置
REDIS_ENABLED=true
```

然后在 `docker-compose.yml` 中使用：

```yaml
services:
  postgres:
    env_file:
      - .env
    environment:
      - POSTGRES_USER=${POSTGRES_USER}
      - POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
```

### 在代码中使用配置

```go
import "github.com/chenhailong/hong3/config"

// 加载配置（从环境变量）
cfg := config.LoadConfig()

// 使用配置
dbDSN := cfg.Database.GetDSN()
redisURL := cfg.Redis.GetRedisURL()
serverAddr := cfg.Server.GetServerAddr()
```

## 设置自定义数据库账号和密码

### 方法 1：在 docker-compose.yml 中直接设置

```yaml
services:
  postgres:
    environment:
      - POSTGRES_USER=admin
      - POSTGRES_PASSWORD=your_strong_password
      - POSTGRES_DB=hong3

  backend:
    environment:
      - DB_USER=admin
      - DB_PASSWORD=your_strong_password
      - DB_NAME=hong3
```

### 方法 2：使用环境变量（推荐）

创建 `.env` 文件：

```bash
POSTGRES_USER=admin
POSTGRES_PASSWORD=your_strong_password
POSTGRES_DB=hong3
```

然后更新 `docker-compose.yml`：

```yaml
services:
  postgres:
    env_file:
      - .env
    environment:
      - POSTGRES_USER=${POSTGRES_USER}
      - POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
      - POSTGRES_DB=${POSTGRES_DB}

  backend:
    env_file:
      - .env
    environment:
      - DB_USER=${POSTGRES_USER}
      - DB_PASSWORD=${POSTGRES_PASSWORD}
      - DB_NAME=${POSTGRES_DB}
```

## Redis 配置

Redis **默认无密码**，如需设置密码：

```yaml
redis:
  command: redis-server --appendonly yes --bind 0.0.0.0 --requirepass your_password

backend:
  environment:
    - REDIS_PASSWORD=your_password
```

## 注意事项

1. **数据库账号密码同步**：`DB_USER`/`DB_PASSWORD` 必须与 `POSTGRES_USER`/`POSTGRES_PASSWORD` 一致
2. **首次启动**：首次设置数据库密码后，需要重新创建数据库容器（删除 volume）
3. **安全性**：生产环境请使用强密码，不要使用默认值
4. **环境变量优先级**：docker-compose.yml 中的环境变量会覆盖 `.env` 文件

