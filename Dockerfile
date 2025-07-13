# 多阶段构建 - 前端构建阶段
FROM node:18-alpine AS frontend-builder

# 设置工作目录
WORKDIR /app

# 复制前端依赖文件
COPY web/package*.json ./

# 安装前端依赖
RUN npm ci --only=production

# 复制前端源代码
COPY web/ ./

# 构建前端
RUN npm run build

# 后端构建阶段
FROM golang:1.24.1-alpine AS backend-builder

# 设置工作目录
WORKDIR /app

# 安装必要的包
RUN apk add --no-cache git

# 复制go mod文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制后端源代码
COPY . .

# 构建后端应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o foodcook cmd/foodcook/main.go

# 最终运行阶段
FROM alpine:latest

# 安装ca-certificates用于HTTPS请求和nginx
RUN apk --no-cache add ca-certificates nginx

# 创建非root用户
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# 设置工作目录
WORKDIR /app

# 从前端构建阶段复制构建结果
COPY --from=frontend-builder /app/dist /var/www/html

# 从后端构建阶段复制二进制文件
COPY --from=backend-builder /app/foodcook .

# 复制配置文件
COPY --from=backend-builder /app/configs ./configs

# 复制nginx配置
COPY docker/nginx.conf /etc/nginx/nginx.conf

# 创建上传目录
RUN mkdir -p uploads && \
    chown -R appuser:appgroup /app && \
    chown -R appuser:appgroup /var/www/html && \
    chown -R appuser:appgroup /var/log/nginx && \
    chown -R appuser:appgroup /var/lib/nginx

# 创建启动脚本
COPY docker/start.sh /start.sh
RUN chmod +x /start.sh

# 切换到非root用户
USER appuser

# 暴露端口
EXPOSE 80 8080

# 运行启动脚本
CMD ["/start.sh"] 