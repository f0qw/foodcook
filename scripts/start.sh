#!/bin/bash

# FoodCook 启动脚本

set -e

echo "🚀 启动 FoodCook 应用..."

# 检查Go是否安装
if ! command -v go &> /dev/null; then
    echo "❌ Go 未安装，请先安装 Go 1.24.1+"
    exit 1
fi

# 检查MySQL是否运行
if ! mysqladmin ping -h127.0.0.1 -uroot -proot --silent; then
    echo "❌ MySQL 未运行，请先启动 MySQL 服务"
    echo "   确保数据库 foodcook 已创建"
    exit 1
fi

# 检查Redis是否运行
if ! redis-cli ping &> /dev/null; then
    echo "❌ Redis 未运行，请先启动 Redis 服务"
    exit 1
fi

# 安装依赖
echo "📦 安装依赖..."
go mod tidy

# 构建应用
echo "🔨 构建应用..."
go build -o bin/foodcook cmd/foodcook/main.go

# 设置环境变量
export DB_HOST=127.0.0.1
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=root
export DB_NAME=foodcook
export REDIS_HOST=127.0.0.1
export REDIS_PORT=6379
export JWT_SECRET=your-secret-key-change-in-production

# 启动应用
echo "🌟 启动应用..."
echo "   应用将在 http://localhost:8080 启动"
echo "   API文档: http://localhost:8080/api"
echo "   按 Ctrl+C 停止应用"
echo ""

./bin/foodcook 