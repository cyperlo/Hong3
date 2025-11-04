# Redis 配置说明

## 安装依赖

在项目根目录运行以下命令安装 Redis 客户端：

```bash
cd backend
go get github.com/redis/go-redis/v9
go mod tidy
```

## 配置

在 `backend/config.yaml` 中配置 Redis：

```yaml
redis:
  host: "redis"
  port: "6379"
  password: ""
  db: 0
  enabled: true
```

或通过环境变量：

```bash
REDIS_ENABLED=true
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

## Token 存储

Token 现在存储在 Redis 中，而不是 PostgreSQL：

- **Key 格式**: `token:{token_value}`
- **过期时间**: 7 天（自动过期）
- **数据结构**: JSON 格式，包含 `user_id`, `username`, `expires_at`

## 优势

1. **性能更好**: Redis 内存存储，读写速度快
2. **自动过期**: Redis 自动删除过期 token，无需手动清理
3. **可扩展**: 支持分布式部署，可以共享 token 状态
4. **减少数据库负载**: Token 查询不再占用数据库资源

