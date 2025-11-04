# 数据库设置指南

## 概述

项目已从内存存储迁移到 PostgreSQL 数据库，使用 GORM 作为 ORM 框架。

## 安装依赖

### 方法 1：使用脚本（推荐）

```bash
cd backend/db
chmod +x install_deps.sh
./install_deps.sh
```

### 方法 2：手动安装

```bash
cd backend
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
go mod tidy
```

## SQL 文件

### 迁移文件

- `backend/db/migrations/001_create_users_table.sql` - 创建用户表和 token 表的 SQL 脚本

### 初始化脚本

- `backend/db/init.sql` - 完整的数据库初始化脚本

## 数据库配置

### Docker Compose（推荐）

使用 docker-compose 启动服务时，PostgreSQL 会自动启动并配置：

```bash
make build
make up
```

数据库配置：
- 主机: `postgres` (容器内) 或 `localhost` (宿主机)
- 端口: `5432`
- 用户名: `postgres`
- 密码: `postgres`
- 数据库名: `hong3`

### 环境变量

可以通过环境变量自定义数据库配置：

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=your_password
export DB_NAME=hong3
export DB_SSLMODE=disable
```

## 自动迁移

应用启动时会自动执行数据库迁移，创建所需的表结构：

- `users` 表 - 存储用户信息
- `tokens` 表 - 存储认证 token

## 手动执行 SQL

如果需要手动初始化数据库：

```bash
# 连接到 PostgreSQL
docker-compose exec postgres psql -U postgres -d hong3

# 或者使用本地 PostgreSQL
psql -h localhost -U postgres -d hong3
```

然后执行 SQL 脚本：

```sql
\i backend/db/migrations/001_create_users_table.sql
```

或者直接复制 `backend/db/init.sql` 的内容执行。

## 数据库表结构

### users 表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | VARCHAR(36) | 主键，用户ID |
| username | VARCHAR(50) | 用户名（唯一） |
| password | VARCHAR(255) | 密码（MD5哈希） |
| name | VARCHAR(100) | 显示名称 |
| created_at | TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | 更新时间 |

### tokens 表

| 字段 | 类型 | 说明 |
|------|------|------|
| token | VARCHAR(36) | 主键，token值 |
| user_id | VARCHAR(36) | 用户ID（外键） |
| username | VARCHAR(50) | 用户名 |
| expires_at | TIMESTAMP | 过期时间 |
| created_at | TIMESTAMP | 创建时间 |

## 注意事项

1. **数据持久化**：数据存储在 Docker volume `postgres-data` 中，删除容器不会丢失数据
2. **自动清理**：系统会每小时自动清理过期的 token
3. **密码安全**：当前使用 MD5 哈希，生产环境建议使用 bcrypt

## 故障排除

### 数据库连接失败

1. 检查 PostgreSQL 容器是否运行：`docker-compose ps`
2. 检查环境变量是否正确设置
3. 检查网络连接：确保后端容器可以访问 postgres 容器

### 表不存在错误

1. 确保自动迁移已执行（查看启动日志）
2. 手动执行 SQL 脚本创建表

