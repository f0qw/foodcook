#!/bin/bash

echo "🍽️  FoodCook 家庭菜单管理系统 - Docker版本"
echo "=============================================="

# 检查Docker是否运行
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker未运行，请先启动Docker"
    exit 1
fi

# 检查Docker Compose是否可用
if ! docker-compose --version > /dev/null 2>&1; then
    echo "❌ Docker Compose不可用，请安装Docker Compose"
    exit 1
fi

# 创建必要的目录
echo "📁 创建必要的目录..."
mkdir -p uploads

# 构建并启动服务
echo "🚀 构建并启动服务..."
docker-compose up -d --build

# 等待服务启动
echo "⏳ 等待服务启动..."
sleep 15

# 检查服务状态
echo "🔍 检查服务状态..."
docker-compose ps

# 检查健康状态
echo "🏥 检查服务健康状态..."
for service in app mysql redis; do
    if docker-compose ps $service | grep -q "Up"; then
        echo "✅ $service 服务运行正常"
    else
        echo "❌ $service 服务启动失败"
    fi
done

echo ""
echo "🎉 启动完成！"
echo "=============================================="
echo "🌐 前端地址: http://localhost"
echo "🔧 后端API: http://localhost/api"
echo "🗄️  数据库: localhost:3306"
echo "🔑 初始账户: root / password"
echo ""
echo "📊 可选工具:"
echo "   phpMyAdmin: http://localhost:8081 (用户名: root, 密码: root)"
echo ""
echo "🔧 常用命令:"
echo "   查看日志: docker-compose logs -f"
echo "   停止服务: docker-compose down"
echo "   重启服务: docker-compose restart"
echo "==============================================" 