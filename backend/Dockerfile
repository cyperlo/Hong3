# 构建阶段
FROM golang:1.24-alpine AS builder

# 设置 Alpine 镜像源（使用阿里云镜像加速）
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 设置工作目录
WORKDIR /app

# 安装构建依赖
RUN apk add --no-cache git

# 设置 Go 代理（使用中国镜像加速下载）
ENV GOPROXY=https://goproxy.cn,direct

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 运行阶段
FROM alpine:latest

# 设置 Alpine 镜像源（使用阿里云镜像加速）
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 安装 ca 证书和 wget（用于健康检查）
RUN apk --no-cache add ca-certificates wget

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./main"]

