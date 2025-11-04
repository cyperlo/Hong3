# 外部连接数据库指南

## PostgreSQL 连接

### 连接信息

- **主机**: `localhost` 或服务器 IP 地址
- **端口**: `5432`
- **数据库**: `hong3`
- **用户名**: `postgres`
- **密码**: `postgres`

### 连接字符串示例

**psql 命令行：**
```bash
psql -h localhost -p 5432 -U postgres -d hong3
```

**连接 URL：**
```
postgresql://postgres:postgres@localhost:5432/hong3?sslmode=disable
```

### 常见客户端工具

1. **DBeaver / DataGrip / pgAdmin**
   - Host: `localhost`
   - Port: `5432`
   - Database: `hong3`
   - User: `postgres`
   - Password: `postgres`

2. **命令行工具**
   ```bash
   psql -h localhost -U postgres -d hong3
   ```

3. **编程语言连接**
   ```go
   // Go (GORM)
   dsn := "host=localhost user=postgres password=postgres dbname=hong3 port=5432 sslmode=disable"
   ```

## Redis 连接

### 连接信息

- **主机**: `localhost` 或服务器 IP 地址
- **端口**: `6379`
- **密码**: 无（如果配置了密码，请使用配置的密码）
- **数据库**: `0`

### 连接字符串示例

**redis-cli 命令行：**
```bash
redis-cli -h localhost -p 6379
```

**连接 URL：**
```
redis://localhost:6379/0
```

### 常见客户端工具

1. **RedisInsight / Redis Desktop Manager**
   - Host: `localhost`
   - Port: `6379`
   - Password: (留空)
   - Database: `0`

2. **命令行工具**
   ```bash
   redis-cli -h localhost -p 6379
   ```

3. **编程语言连接**
   ```go
   // Go (go-redis)
   rdb := redis.NewClient(&redis.Options{
       Addr:     "localhost:6379",
       Password: "",
       DB:       0,
   })
   ```

## 故障排除

### 无法连接 PostgreSQL

1. **检查容器是否运行**
   ```bash
   docker-compose ps
   ```

2. **检查端口是否被占用**
   ```bash
   netstat -tuln | grep 5432
   # 或
   lsof -i :5432
   ```

3. **检查防火墙/安全组**
   - 确保 5432 端口已开放
   - 如果使用云服务器，检查安全组规则

4. **查看 PostgreSQL 日志**
   ```bash
   docker-compose logs postgres
   ```

5. **测试连接**
   ```bash
   docker-compose exec postgres psql -U postgres -d hong3 -c "SELECT version();"
   ```

### 无法连接 Redis

1. **检查容器是否运行**
   ```bash
   docker-compose ps
   ```

2. **检查端口是否被占用**
   ```bash
   netstat -tuln | grep 6379
   # 或
   lsof -i :6379
   ```

3. **检查防火墙/安全组**
   - 确保 6379 端口已开放
   - 如果使用云服务器，检查安全组规则

4. **查看 Redis 日志**
   ```bash
   docker-compose logs redis
   ```

5. **测试连接**
   ```bash
   docker-compose exec redis redis-cli ping
   ```

## 安全建议

### 生产环境

1. **修改默认密码**
   - 在 `docker-compose.yml` 中设置强密码
   - 或使用环境变量文件

2. **限制访问**
   - 使用防火墙限制只允许特定 IP 访问
   - 配置 PostgreSQL 的 `pg_hba.conf` 限制来源 IP

3. **启用 SSL**
   - PostgreSQL: 设置 `sslmode=require`
   - Redis: 使用 Redis 6+ 的 TLS 支持

4. **不要使用 `--protected-mode no`**
   - 生产环境应该设置 Redis 密码并启用保护模式

## 从远程服务器连接

如果服务器在远程，使用服务器的公网 IP：

```bash
# PostgreSQL
psql -h <服务器IP> -p 5432 -U postgres -d hong3

# Redis
redis-cli -h <服务器IP> -p 6379
```

**注意**：确保服务器的防火墙和安全组已开放相应端口。

