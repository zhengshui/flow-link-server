# ============================================
# 构建阶段 (Builder Stage)
# ============================================
FROM golang:1.21-alpine AS builder

# 安装构建所需的工具
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /build

# 先复制依赖文件，利用 Docker 层缓存
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# 复制源代码
COPY . .

# 构建二进制文件
# CGO_ENABLED=0: 禁用 CGO，生成静态链接的二进制文件
# -ldflags="-w -s": 去除调试信息，减小二进制体积
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.Version=${VERSION:-dev}" \
    -o main cmd/main.go

# ============================================
# 运行阶段 (Runtime Stage)
# ============================================
FROM alpine:3.19

# 安装运行时必要的包
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /build/main .

# 创建非 root 用户运行应用
RUN addgroup -g 1000 appgroup && \
    adduser -D -u 1000 -G appgroup appuser && \
    chown -R appuser:appgroup /app

USER appuser

# 声明端口（实际端口由环境变量控制）
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:${PORT:-8080}/health || exit 1

# 启动应用
CMD ["./main"]
